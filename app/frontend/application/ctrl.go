package application

import (
	"bytes"
	"fmt"
	"fractalmri/app/backend/domain/decoder"
	"fractalmri/app/frontend/application/client"
	"fractalmri/app/frontend/domain/model"
	"fractalmri/app/frontend/domain/msg"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"image"
	"image/color"
	"image/png"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var Ctrl = &MainController{}

type MainController struct {
	chart app.Value
	mtx   sync.Mutex
}

func (ctrl *MainController) Login(ctx app.Context) {
	var profile model.Profile
	ctx.LocalStorage().Get("profile", &profile)
	c := client.New()
	log.Printf("profile: {id:%d token:\"%s\"}", profile.ID, profile.Token)
	err := c.Post("/profile", profile, &model.M.Profile)
	_, mes := msg.NewUnauthorizedMsg()
	if err != nil {
		if mes.Message == err.Error() {
			if err := c.Post("/login", &profile, &model.M.Profile); err != nil {
				app.Window().Call("alert", err.Error())
				app.Window().Get("location").Call("reload")
			}
			ctx.LocalStorage().Set("profile", &model.M.Profile)
			log.Printf("login-profile: {id:%d token:\"%s\"}", model.M.Profile.ID, model.M.Profile.Token)
		} else {
			app.Window().Call("alert", err.Error())
			app.Window().Get("location").Call("reload")
		}
		ctx.LocalStorage().Set("profile", &model.M.Profile)
	}
}

func (ctrl *MainController) LastSelectionURL() string {
	if model.M.LastSelection == -1 {
		return ""
	}
	return model.M.Selection[model.M.LastSelection].URL
}

func (ctrl *MainController) GetList() {
	c := client.New()
	list := []*model.ImageModel{}
	count := 0
	var err error
	for err = c.Get("/list", &list); err != nil && count < 10; err = c.Get("/list", &list) {
		time.Sleep(time.Second * time.Duration(count))
		count++
	}
	if count >= 10 {
		log.Fatalln(err)
	}
	for i := 0; i < len(list); i++ {
		img := list[i]
		model.M.Images[img.ID] = img
	}
}

func (ctrl *MainController) GetResults() []*model.ResultModel {
	c := client.New()
	results := []*model.ResultModel{}
	count := 0
	var err error
	for err = c.Get("/results", &results); err != nil && count < 10; err = c.Get("/results", &results) {
		time.Sleep(time.Second * time.Duration(count))
		count++
	}
	if count >= 10 {
		log.Fatalln(err)
	}
	return results
}

func (ctrl *MainController) HandleOnDeleteImage(ctx app.Context, event app.Event) {
	ctx.Async(func() {
		id := model.M.LastSelection
		if id != -1 {
			c := client.New()
			list := []*model.ImageModel{}
			var err error
			if err = c.Get(fmt.Sprintf("/delete?id=%d", id), &list); err != nil {
				log.Println(err)
				return
			}
			model.M.Images = make(map[int]*model.ImageModel)
			for i := 0; i < len(list); i++ {
				img := list[i]
				model.M.Images[img.ID] = img
			}
			ctx.SetState("/image-list", model.M.Images)
		}
	})
}

func (ctrl *MainController) HandleOnChangeStep(ctx app.Context, event app.Event) {
	stepStr := event.Get("target").Get("value").String()
	step, err := strconv.ParseFloat(stepStr, 64)
	if err != nil {
		return
	}
	ctx.SetState("/step", step)
}

func (ctrl *MainController) LoadList(ctx app.Context) {
	ctx.Async(func() {
		ctrl.GetList()
		ctx.Dispatch(func(ctx app.Context) {
			ctx.SetState("/image-list", model.M.Images)
		})
	})
}

func (ctrl *MainController) HandleUpload(ctx app.Context, e app.Event) {
	fileInput := app.Window().GetElementByID("upload-file-input")
	files := fileInput.Get("files")
	if files.Length() > 0 {
		file := files.Index(0)
		ctx.SetState("/is-uploading", true)
		elem := app.Window().GetElementByID("modal-upload")
		instance := app.Window().Get("M").Get("Modal").Call("getInstance", elem)
		instance.Call("open")
		ctx.Async(func() {
			var err error
			var list model.ImageModel
			c := client.New()
			err = c.Upload("/upload", file, &list)
			fileInput.Set("value", "")
			ctx.Dispatch(func(ctx app.Context) {
				if err != nil {
					ctx.SetState("/upload-err", err.Error())
					if mStr := err.Error(); mStr != msg.ImageFormatErrStr && mStr != msg.ExceededFileSizeStr {
						ctx.SetState("/allow-retry", true)
					} else {
						ctx.SetState("/allow-retry", false)
					}
				} else {
					instance.Call("close")
					ctx.SetState("/allow-retry", false)
					ctx.SetState("/upload-err", "")
				}
				ctx.SetState("/is-uploading", false)
				ctrl.LoadList(ctx)
			})
		})
	}
}

func (ctrl *MainController) HandleSelect(ctx app.Context, e app.Event) {
	elem := e.Get("target")
	ID := ImageIDFromCheckID(elem.Get("id").String())
	if elem.Get("checked").Bool() {
		ctx.SetState("/last-selection", ID)
		model.M.Selection[ID] = model.NewImageModel(ID)
	} else {
		delete(model.M.Selection, ID)
		if model.M.LastSelection == ID {
			ctx.SetState("/last-selection", -1)
		}
		if len(model.M.Selection) == 0 {
			ctx.SetState("/last-selection", -1)
		}
	}
}

func ImageIDFromCheckID(checkID string) int {
	imageIDStr := strings.TrimLeft(checkID, "image-")
	ID, _ := strconv.Atoi(imageIDStr)
	return ID
}

func (ctrl *MainController) HandleRun(ctx app.Context, e app.Event) {
	ctx.SetState("/is-waiting-result", true)
	mfsWindowStr := app.Window().GetElementByID("input-mfs-window").Get("value").String()
	minUmbralStr := app.Window().GetElementByID("min-umbral-input").Get("value").String()
	maxUmbralStr := app.Window().GetElementByID("max-umbral-input").Get("value").String()
	stepUmbralStr := app.Window().GetElementByID("step-umbral-input").Get("value").String()
	mfsWindow, _ := strconv.Atoi(mfsWindowStr)
	minUmbral, err := strconv.ParseFloat(minUmbralStr, 64)
	if err != nil {
		app.Window().GetElementByID("min-umbral-input").Set("value", model.M.MinUmbral)
	}
	maxUmbral, err := strconv.ParseFloat(maxUmbralStr, 64)
	if err != nil {
		app.Window().GetElementByID("max-umbral-input").Set("value", model.M.MaxUmbral)
	}
	stepUmbral, err := strconv.ParseFloat(stepUmbralStr, 64)
	if err != nil {
		app.Window().GetElementByID("step-umbral-input").Set("value", model.M.Step)
	}
	if maxUmbral < minUmbral {
		maxUmbral, minUmbral = minUmbral, maxUmbral
	} else if maxUmbral == minUmbral {
		maxUmbral = minUmbral + stepUmbral
	}
	app.Window().GetElementByID("min-umbral-input").Set("value", minUmbral)
	app.Window().GetElementByID("max-umbral-input").Set("value", maxUmbral)
	ID := model.M.LastSelection
	ctx.Async(func() {
		c := client.New()
		result := model.ResultModel{}
		err := c.Post("/run", &model.RunModel{
			ImageID:    ID,
			Window:     mfsWindow,
			MinUmbral:  minUmbral,
			MaxUmbral:  maxUmbral,
			StepUmbral: stepUmbral,
		},
			&result)
		if err != nil {
			ctx.SetState("/run-err", err.Error())
		} else {
			model.M.Results[ID] = &result
			ctx.SetState("/current-result", result.ID)
		}
		ctx.SetState("/is-waiting-result", false)
	})
}

func (ctrl *MainController) DrawMFS(id int, ctx app.Context) (err error) {
	ctx.Async(func() {
		ctrl.mtx.Lock()
		defer ctrl.mtx.Unlock()
		c := client.New()
		var mfs []byte
		mfs, err = c.Download(fmt.Sprintf("/mfs?id=%d", id))
		if err != nil {
			return
		}
		var imgMfs pipes.Image64
		imgMfs, err = decoder.DecodeMFS(mfs)
		if err != nil {
			return
		}
		if ctrl.chart != nil {
			ctrl.chart.Call("destroy")
		}
		img := imgMfs.Normalized().Scaled(255)
		DrawImage(img.Width(), img.Height(), func(x, y int) (r uint8, g uint8, b uint8, a uint8) {
			c := uint8(img.At(x, y))
			return c, c, c, 255
		})
	})
	return
}

func (ctrl *MainController) HandleShowResult(ctx app.Context, e app.Event) {
	target := e.Get("currentTarget")
	idStr := target.Get("id").String()[len("button-"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(id)
	ctx.SetState("/current-result", id)
	ctx.SetState("/show-result", true)
	ctrl.DrawMFS(id, ctx)
}

func (ctrl *MainController) HandleShowFDChart(ctx app.Context, e app.Event) {
	ctrl.CloseModalResultList()
	ctx.SetState("/has-sg-tool", false)
	ctx.Async(func() {
		ctrl.mtx.Lock()
		defer ctrl.mtx.Unlock()
		ctx.SetState("/is-mfs", false)
		target := e.Get("currentTarget")
		loglogIndexStr := target.Get("id").String()[len("button-loglog-"):]
		loglogIndex, err := strconv.Atoi(loglogIndexStr)
		if err != nil {
			log.Panicln(err)
		}
		result := model.M.GetCurrentResult()
		canvas := app.Window().GetElementByID("canvas-vr")
		gctx := canvas.Call("getContext", "2d")
		counts, labels := []float64{}, []string{}
		sizes := cfg.GetCFG().BoxSizes
		logSizes := result.FDs[loglogIndex].LogSizes
		for i := len(sizes) - 1; i >= 0; i-- {
			labels = append(labels, fmt.Sprintf("%d (%.1f)", sizes[i], logSizes[i]))
		}
		counts = append(counts, result.FDs[loglogIndex].LogMeasure...)
		umbral := result.Umbrals[loglogIndex]
		fd := result.FDs[loglogIndex].FD
		chartParams := ValueOf(map[string]any{
			"type": "line",
			"data": map[string]any{
				"labels": labels,
				"datasets": []any{
					map[string]any{
						"label":                  fmt.Sprintf("Dimensión Fractal: %.5f Umbral[%.2f, %.2f]", fd, umbral.Min, umbral.Max),
						"data":                   counts,
						"fill":                   false,
						"borderColor":            "rgb(75, 192, 192)",
						"tension":                0.1,
						"cubicInterpolationMode": "monotone",
					},
				},
			},
			"options": map[string]any{
				"responsive": true,
				"scales": map[string]any{
					"x": map[string]any{
						"title": map[string]any{
							"display": true,
							"text":    "log(count)",
						},
					},
					"y": map[string]any{
						"display": true,
						"title": map[string]any{
							"display": true,
							"text":    "log(1/size)",
						},
					},
				},
			},
		})
		if ctrl.chart != nil {
			ctrl.chart.Call("destroy")
		}
		ctrl.chart = app.Window().Get("Chart").New(gctx, chartParams)
	})
}

func (ctrl *MainController) HandleShowMFS(ctx app.Context, e app.Event) {
	ctrl.CloseModalResultList()
	ctx.SetState("/has-sg-tool", false)
	ctrl.DrawMFS(model.M.CurrentResult, ctx)
}

func (ctrl *MainController) HandleShowResultList(ctx app.Context, e app.Event) {
	ctx.SetState("/show-result", false)
}

func (ctrl *MainController) CloseModalResultList() {
	elem := app.Window().GetElementByID("results-modal")
	instance := app.Window().Get("M").Get("Modal").Call("getInstance", elem)
	instance.Call("close")
}

func (ctrl *MainController) HandleShowSegTool(ctx app.Context, e app.Event) {
	ctrl.CloseModalResultList()
	ctx.SetState("/has-sg-tool", true)
	ctrl.DrawSegmentation(model.M.CurrentResult, ctx)
}

func (ctrl *MainController) HandleShowSegResult(ctx app.Context, e app.Event) {
	ctrl.DrawSegmentation(model.M.CurrentResult, ctx)
}

func (ctrl *MainController) DrawSegmentation(id int, ctx app.Context) (err error) {
	ctx.Async(func() {
		//Bloquear el controlador para que no puedan ser usado los recursos compartidos en otras tareas
		defer ctrl.mtx.Unlock()
		ctrl.mtx.Lock()
		//Cargar las imagenes desde el servidor
		resID := model.M.CurrentResult
		//Tarea 2: cargar imagen de resonancia
		c := client.New()
		var buffer []byte
		buffer, err = c.Download(fmt.Sprintf("/segment?id=%d&min=%f&max=%f", resID, model.M.SegToolMin, model.M.SegToolMax))
		if err != nil {
			return
		}
		reader := bytes.NewReader(buffer)
		var img image.Image
		img, err = png.Decode(reader)
		if err != nil {
			return
		}
		if ctrl.chart != nil { //Eliminar grafico chart si lo hay
			ctrl.chart.Call("destroy")
		}
		bd := img.Bounds()
		DrawImage(bd.Dx(), bd.Dy(), func(x, y int) (r uint8, g uint8, b uint8, a uint8) {
			rgba := img.At(x, y).(color.RGBA)
			return rgba.R, rgba.G, rgba.B, rgba.A
		})
		ctx.SetState("/is-waiting-mfs", false)
	})
	return
}

func DrawImage(imgWidth, imgHeight int, pixel func(x, y int) (r, g, b, a uint8)) {
	canvas := app.Window().GetElementByID("canvas-vr")
	gctx := canvas.Call("getContext", "2d") //Obtener contexto grafico del canvas
	width := canvas.Get("width").Int()      //Obtener ancho del canvas
	height := canvas.Get("height").Int()    //Obtener alto del canvas
	relX, relY := float64(imgWidth)/float64(width), float64(imgHeight)/float64(height)
	transform := func(x, y int) (int, int) {
		return int(float64(x) * relX), int(float64(y) * relY)
	}
	gctx.Call("clearRect", 0, 0, width, height) //Limpiar canvas
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := pixel(transform(x, y))
			gctx.Set("fillStyle", fmt.Sprintf("rgba(%d, %d, %d, %d)", r, g, b, a)) //Establecer color del pixel
			gctx.Call("fillRect", x, y, 1, 1)                                      //Dibujar el pixel
		}
	}
}

func (ctrl *MainController) HandleSgOnChange(ctx app.Context, event app.Event) {
	minInput := app.Window().GetElementByID("min-seg-input")
	maxInput := app.Window().GetElementByID("max-seg-input")
	minStr := minInput.Get("value").String()
	maxStr := maxInput.Get("value").String()
	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		min = model.M.SegToolMin
	}
	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		max = model.M.SegToolMax
	}
	if min >= max {
		min = model.M.SegToolMin
		max = model.M.SegToolMax
	}
	minInput.Set("value", min)
	maxInput.Set("value", max)
	model.M.SegToolMin = min
	model.M.SegToolMax = max
}

func (ctrl *MainController) HandleDownloadResult(ctx app.Context, event app.Event) {
	if model.M.CurrentResult != -1 {
		log.Println("handleDownloadResult")
		document := app.Window().Get("document")
		a := document.Call("createElement", "a")
		a.Set("href", fmt.Sprintf("/download-result?id=%d", model.M.CurrentResult))
		a.Set("download", fmt.Sprintf("result-%d.zip", model.M.CurrentResult))
		document.Get("body").Call("appendChild", a)
		a.Call("click")
		document.Get("body").Call("removeChild", a)
	}
}

func ValueOf(value any) app.Value {
	var jsv app.Value
	switch v := value.(type) {
	case map[string]any:
		jsv = app.ValueOf(map[string]any{})
		for k, v := range v {
			jsv.Set(k, ValueOf(v))
		}
	case []float64:
		jsv = app.ValueOf([]any{})
		for i := 0; i < len(v); i++ {
			jsv.Call("push", app.ValueOf(v[i]))
		}
	case []string:
		jsv = app.ValueOf([]any{})
		for i := 0; i < len(v); i++ {
			jsv.Call("push", app.ValueOf(v[i]))
		}
	case []any:
		jsv = app.ValueOf([]any{})
		for i := 0; i < len(v); i++ {
			jsv.Call("push", ValueOf(v[i]))
		}
	default:
		jsv = app.ValueOf(v)
	}
	return jsv
}
