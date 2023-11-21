package application

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"

	"fractalmri/app/backend/domain/config"
	"fractalmri/app/backend/domain/decoder"
	"fractalmri/app/backend/domain/model"
	md "fractalmri/app/frontend/domain/model"
	"fractalmri/app/frontend/domain/msg"
	webapp "fractalmri/app/frontend/presentation"
	"fractalmri/gofract/pipes"

	"github.com/gin-gonic/gin"
	"github.com/maxence-charriere/go-app/v9/pkg/app"

	"crypto/rand"
	"encoding/hex"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/stellviaproject/zipper"
)

type App struct {
	router *gin.Engine
	webapp *app.Handler
}

func NewApp(cfgs ...config.WithConfigFn) *App {
	a := new(App)
	if model.C == nil {
		model.C = config.Default()
	}
	for i := range cfgs {
		cfgs[i](model.C)
	}
	os.Mkdir(model.C.LogPath, os.ModePerm)
	os.Mkdir(model.C.StorePath, os.ModePerm)
	os.Mkdir(model.C.UserPath, os.ModePerm)
	model.Load()
	return a
}

func (a *App) SetConfig(config *config.Configuration) {
	model.C = config
}

func (a *App) Run(addr string) {
	a.router = gin.New()
	// Iniciar engine
	a.router = gin.Default()
	// Crear clave privada
	pkey := []byte(generateToken())
	// Configurar las opciones de sesion
	store := cookie.NewStore(pkey)
	a.router.Use(sessions.Sessions("sessions", store))
	// Configurar las rutas
	//GET
	//a.router.GET("/", a.Main)
	a.RootAll()

	a.webapp = webapp.NewApp(a.cacheResources()...)
	a.router.Use(a.Main)
	//TODO: Habilitar y probar antes de entregar
	//file := logger.ConfigureLog(a.config.LogPath)
	//defer file.Close()
	go a.router.Run(addr)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("the service has received an interruption signal.\nstopping the service...")
	a.OnExit()
}

func generateToken() []byte {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return []byte("ak29jc9sjd89aduhe83y4837eyudhsad8uasyd7a8sydh3")
	}
	return []byte(hex.EncodeToString(b))
}

func (a *App) RootAll() {
	a.router.POST("/login", a.Login)
	a.router.POST("/profile", a.Profile)
	a.router.GET("/list", a.List)
	a.router.POST("/upload", a.Upload)
	a.router.POST("/run", a.RunAnalysis)
	a.router.GET("/loglog", a.LogLog)
	a.router.GET("/mfs", a.MFS)
	a.router.GET("/image", a.Image)
	a.router.GET("/results", a.Results)
	a.router.GET("/download-result", a.DownloadResult)
	a.router.GET("/segment", a.Segment)
	a.router.GET("/delete", a.DeleteImage)
}

func (a *App) Main(c *gin.Context) {
	a.webapp.ServeHTTP(c.Writer, c.Request)
}

func (a *App) OnExit() {

}

func (a *App) cacheResources() []string {
	queue := []string{model.C.WebDir}
	cache := []string{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		dir, err := os.ReadDir(cur)
		queue, cache = a.cache(cur, queue, cache, dir, err)
	}
	return cache
}

func (a *App) cache(cur string, queue, cache []string, dir []fs.DirEntry, err error) (q, c []string) {
	if err != nil {
		log.Println(err)
	}
	for i := range dir {
		if dir[i].IsDir() {
			queue = append(queue, path.Join(cur, dir[i].Name()))
		} else {
			fpath := path.Join(cur, dir[i].Name())
			log.Printf("Cache: %v", fpath)
			cache = append(cache, fpath)
		}
	}
	return queue, cache
}

// Obtiene el usuario correspondiente a la sesion
func (a *App) GetUser(c *gin.Context) *model.UserModel {
	session := sessions.Default(c)   //Obtener la sesion actual
	userid := session.Get("user_id") //Obtener el identificador de la sesion actual
	if userid == nil {               //Si no hay identificador de la sesion
		return nil //retornar nil
	}
	return model.M.GetUser(userid.(int))
}

func (a *App) Login(c *gin.Context) {
	u := a.GetUser(c)
	session := sessions.Default(c)
	if u == nil {
		profile := new(md.Profile)
		if err := c.BindJSON(profile); err != nil {
			c.JSON(msg.NewUnauthorizedMsg())
			return
		}
		u = model.M.GetUser(profile.ID)
		if u != nil && profile.ID != 0 {
			if u == nil {
				c.JSON(msg.NewUnauthorizedMsg())
				return
			}
			if profile.Token != u.Token {
				c.JSON(msg.NewTokenAuthFailMsg())
				return
			}
			session.Set("user_id", profile.ID)
			session.Save()
			c.JSON(msg.NewStartMsg(profile))
			return
		}
	}
	u = model.M.NewUser(string(generateToken()))
	session.Set("user_id", u.ID)
	if err := session.Save(); err != nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	profile := &md.Profile{
		Token: u.Token,
		ID:    u.ID,
	}
	c.JSON(msg.NewStartMsg(profile))
}

// Obtiene la informacion del usuario
// La url es "/profile" y el metodo es GET.
func (a *App) Profile(c *gin.Context) {
	//Obtiene el usuario de la sesion actual
	user := a.GetUser(c)
	//Si el usuario no ha iniciado sesion
	if user == nil {
		//Indicar que no hay usuario autenticado
		c.JSON(msg.NewUnauthorizedMsg())
	} else {
		//El usuario esta autenticado
		profile := new(md.Profile)
		if err := c.BindJSON(profile); err != nil {
			c.JSON(msg.NewProfileFailMsg())
			return
		}
		if user.ID != profile.ID {
			c.JSON(msg.NewUserIDAuthFailMsg())
			return
		}
		if user.Token != profile.Token {
			c.JSON(msg.NewTokenAuthFailMsg())
			return
		}
		profile.ID = user.ID
		profile.Token = user.Token
		c.JSON(msg.NewStartMsg(profile))
	}
}

func (a *App) Upload(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(msg.NewUploadFailedMsg())
		return
	}
	if file.Size > model.C.UploadMaxSize {
		c.JSON(msg.NewFileSizeExceededMsg())
		return
	}
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(msg.NewUploadFailedMsg())
		return
	}
	defer fileBytes.Close()
	data, err := io.ReadAll(fileBytes)
	if err != nil {
		c.JSON(msg.NewUploadFailedMsg())
		return
	}
	images, err := Decode(data)
	if err != nil {
		c.JSON(msg.NewImageFormatErrMsg())
		return
	}
	lastID := -1
	for i := 0; i < len(images); i++ {
		lastID = u.SaveImage(images[i])
	}
	c.JSON(msg.NewImageUploadOk(md.NewImageModel(lastID)))
}

func (a *App) List(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	list := make([]*md.ImageModel, 0, len(u.Images))
	for id := range u.Images {
		list = append(list, md.NewImageModel(id))
	}
	c.JSON(msg.NewImageListOk(list))
}

func (a *App) Results(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	list := make([]*md.ResultModel, 0, len(u.Results))
	for _, result := range u.Results {
		list = append(list, result)
	}
	c.JSON(msg.NewResultListOk(list))
}

func (a *App) RunAnalysis(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	var runInfo md.RunModel
	if err := c.BindJSON(&runInfo); err != nil {
		c.JSON(msg.NewRunFailedMsg())
		return
	}
	img := u.GetImage(runInfo.ImageID)
	if img == nil {
		c.JSON(msg.NewImageNoExistMsg())
		return
	}
	p := NewProcess(img, &runInfo)
	p.Run(u)
	result := u.GetResult(runInfo.ImageID)
	c.JSON(msg.NewResultMsg(result))
}

func (a *App) LogLog(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	resultIDStr := c.Query("id")
	imageIndexStr := c.Query("index")
	resultID, err := strconv.Atoi(resultIDStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	imageIndex, err := strconv.Atoi(imageIndexStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	data, err := u.GetResultImage(resultID, imageIndex)
	if err != nil {
		c.JSON(msg.NewImageNoExistMsg())
		return
	}
	c.Data(http.StatusOK, "image/png", data)
}

func (a *App) MFS(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	resultIDStr := c.Query("id")
	resultID, err := strconv.Atoi(resultIDStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	data, err := u.GetResultMFS(resultID)
	if err != nil {
		c.JSON(msg.NewImageNoExistMsg())
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (a *App) DeleteImage(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	imageIDStr := c.Query("id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	u.DeleteImage(imageID)
	list := make([]*md.ImageModel, 0, len(u.Images))
	for id := range u.Images {
		list = append(list, md.NewImageModel(id))
	}
	c.JSON(msg.NewImageListOk(list))
}

func (a *App) Image(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	imageIDStr := c.Query("id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	img := u.GetImage(imageID)
	if img == nil {
		c.JSON(msg.NewImageNotFoundMsg())
		return
	}
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, img)
	if err != nil {
		c.JSON(msg.NewImageFormatErrMsg())
		return
	}
	c.Data(http.StatusOK, "image/png", buffer.Bytes())
}

func (a *App) Segment(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	minStr := c.Query("min")
	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	maxStr := c.Query("max")
	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	img := u.GetImage(id)
	if img == nil {
		c.JSON(msg.NewImageNotFoundMsg())
		return
	}
	data, err := u.GetResultMFS(id)
	if err != nil {
		c.JSON(msg.NewImageNoExistMsg())
		return
	}
	mfs, err := decoder.DecodeMFS(data)
	if err != nil {
		c.JSON(msg.NewMFSDecodeErrMsg())
		return
	}
	width, height := mfs.Width(), mfs.Height()
	output := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			fd := mfs.At(x, y)
			if min <= fd && max >= fd {
				output.Set(x, y, img.At(x, y))
				// if IsBorder(x, y, width, height, min, max, mfs) {
				// 	for i := -2; i <= 2; i++ {
				// 		for j := -2; j <= 2; j++ {
				// 			if !IsOut(x, y, width, height) {
				// 				output.Set(x, y, color.RGBA{
				// 					R: 255,
				// 					G: 0,
				// 					B: 0,
				// 					A: 255,
				// 				})
				// 			}
				// 		}
				// 	}
				// } else {
				// 	output.Set(x, y, img.At(x, y))
				// }
			} else {
				color := img.At(x, y).(color.RGBA)
				color.A = uint8(float64(color.A) * 0.1)
				output.Set(x, y, img.At(x, y))
			}
		}
	}
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, output)
	if err != nil {
		c.JSON(msg.NewImageFormatErrMsg())
		return
	}
	c.Data(http.StatusOK, "image/png", buffer.Bytes())
}

func IsBorder(x, y, w, h int, min, max float64, mfs pipes.Image64) bool {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			xp, yp := x+i, y+j
			if xp != x && yp != y && !IsOut(xp, yp, w, h) {
				fd := mfs.At(xp, yp)
				return fd < min || fd > max
			}
		}
	}
	return false
}

func IsOut(x, y, w, h int) bool {
	return x-1 < 0 || y-1 < 0 || x+1 >= w || y+1 >= h
}

func (a *App) DownloadResult(c *gin.Context) {
	u := a.GetUser(c)
	if u == nil {
		c.JSON(msg.NewUnauthorizedMsg())
		return
	}
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(msg.NewInvalidURLParamMgs())
		return
	}
	srcFile := path.Join(model.C.UserPath, fmt.Sprintf("%d.json", u.ID))
	dstFile := path.Join(model.C.StorePath, fmt.Sprintf("%d/%d/data.json", u.ID, id))
	defer func() {
		if err := os.Remove(dstFile); err != nil {
			log.Println(err)
		}
	}()
	inFile, err := os.Open(srcFile)
	if err != nil {
		c.JSON(msg.NewFileNotFound())
		return
	}
	outFile, err := os.Create(dstFile)
	if err != nil {
		c.JSON(msg.NewFileNotFound())
		return
	}
	_, err = io.Copy(outFile, inFile)
	if err != nil {
		c.JSON(msg.NewFileNotFound())
		return
	}
	inFile.Close()
	outFile.Close()
	folderPath := path.Join(model.C.StorePath, fmt.Sprintf("%d", u.ID))
	zipData, err := zipper.ZipFolder(folderPath)
	if err != nil {
		c.JSON(msg.NewFileNotFound())
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", zipData)
}
