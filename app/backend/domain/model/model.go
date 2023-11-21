package model

import (
	"encoding/json"
	"fmt"
	"fractalmri/app/backend/domain/config"
	md "fractalmri/app/frontend/domain/model"
	"image"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path"
	"sync"
	"time"
)

var (
	M *MainModel
	C *config.Configuration
)

func Load() {
	M = &MainModel{}
	loadMainModel()
	saveMainModel()
}

func saveMainModel() {
	M.rw.Lock()
	defer M.rw.Unlock()
	file := Create(C.MainFile)
	defer file.Close()
	type MainModelJSON struct {
		Users   []int `json:"users"`
		Counter int   `json:"counter"`
	}
	main := &MainModelJSON{
		Users:   make([]int, 0, len(M.users)),
		Counter: M.userID,
	}
	for id, user := range M.users {
		main.Users = append(main.Users, id)
		userPath := path.Join(C.UserPath, fmt.Sprintf("%d.json", id))
		user.saveUser(userPath)
	}
	data, err := json.Marshal(main)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err = file.Write(data); err != nil {
		log.Fatalln(err)
	}
}

func loadMainModel() {
	M.rw.Lock()
	defer M.rw.Unlock()
	file, err := os.Open(C.MainFile)
	if err != nil {
		newMainModel()
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		newMainModel()
		return
	}
	type MainModelJSON struct {
		Users   []int `json:"users"`
		Counter int   `json:"counter"`
	}
	main := &MainModelJSON{}
	if err = json.Unmarshal(data, main); err != nil {
		log.Fatalln(err)
	}
	M = &MainModel{
		users:  make(map[int]*UserModel, len(main.Users)),
		userID: main.Counter,
	}
	for i := 0; i < len(main.Users); i++ {
		userPath := path.Join(C.UserPath, fmt.Sprintf("%d.json", main.Users[i]))
		user := &UserModel{}
		user.loadUser(userPath)
		M.users[main.Users[i]] = user
	}
}

func Create(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm|os.ModeDevice)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

type MainModel struct {
	users  map[int]*UserModel
	userID int
	rw     sync.RWMutex
}

func newMainModel() *MainModel {
	M = &MainModel{
		users:  make(map[int]*UserModel, 100),
		userID: 0,
	}
	return M
}

func (m *MainModel) NewUser(token string) *UserModel {
	m.rw.Lock()
	m.userID++
	if m.userID == math.MaxInt {
		m.userID = 1
	}
	user := NewUserModel(m.userID)
	user.Token = token
	m.users[user.ID] = user
	m.rw.Unlock()
	saveMainModel()
	return user
}

func (m *MainModel) UserIDs() []int {
	ids := make([]int, 0, len(m.users))
	for id, _ := range m.users {
		ids = append(ids, id)
	}
	return ids
}

func (m *MainModel) GetUser(id int) *UserModel {
	user, ok := m.users[id]
	if ok {
		return user
	}
	return nil
}

type UserModel struct {
	ID       int
	Images   map[int]*ImageModel     `json:"images"`
	Results  map[int]*md.ResultModel `json:"results"`
	Image_id int                     `json:"image-id"`
	Token    string                  `json:"token"`
	mtx      sync.Mutex
	lastTime time.Time
}

func (user *UserModel) UpdateTime() {
	user.lastTime = time.Now()
}

func (user *UserModel) DeleteStore() {
	if time.Since(user.lastTime) > C.MaxTime {
		for _, img := range user.Images {
			os.Remove(img.FileName)
		}
		user.Images = make(map[int]*ImageModel)
		os.RemoveAll(path.Join(C.StorePath, fmt.Sprintf("%d", user.ID)))
		user.Results = make(map[int]*md.ResultModel)
	}
}

func (user *UserModel) DeleteImage(id int) {
	user.mtx.Lock()
	defer user.mtx.Unlock()
	delete(user.Images, id)
	os.Remove(path.Join(C.StorePath, fmt.Sprintf("%d/%d.png", user.ID, id)))
	os.RemoveAll(path.Join(C.StorePath, fmt.Sprintf("%d/%d", user.ID, id)))
	delete(user.Results, id)
	user.saveUser(path.Join(C.UserPath, fmt.Sprintf("%d.json", user.ID)))
}

func (user *UserModel) loadUser(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, user); err != nil {
		return err
	}
	os.Mkdir(path.Join(C.StorePath, fmt.Sprintf("%d", user.ID)), os.ModePerm)
	return nil
}

func (user *UserModel) saveUser(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm|os.ModeDevice)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func NewUserModel(id int) *UserModel {
	os.Mkdir(path.Join(C.StorePath, fmt.Sprintf("%d", id)), os.ModePerm)
	return &UserModel{
		ID:       id,
		Images:   make(map[int]*ImageModel),
		Results:  make(map[int]*md.ResultModel),
		Image_id: 0,
		lastTime: time.Now(),
	}
}

func (m *UserModel) SaveImage(img image.Image) int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	currID := m.Image_id
	m.Image_id++
	imgMd := &ImageModel{
		ID:  currID,
		Img: img,
	}
	m.Images[currID] = imgMd
	fileName := path.Join(C.StorePath, fmt.Sprintf("%d/%d.png", m.ID, currID))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModeDevice|os.ModePerm)
	if err != nil {
		log.Println(err)
		return imgMd.ID
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Println(err)
		return imgMd.ID
	}
	m.saveUser(path.Join(C.UserPath, fmt.Sprintf("%d.json", m.ID)))
	return imgMd.ID
}

func (m *UserModel) GetImage(id int) image.Image {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	img := m.Images[id]
	if img.Img == nil {
		file, err := os.Open(path.Join(C.StorePath, fmt.Sprintf("%d/%d.png", m.ID, id)))
		if err != nil {
			return nil
		}
		defer file.Close()
		img.Img, err = png.Decode(file)
		if err != nil {
			return nil
		}
	}
	return img.Img
}

func (m *UserModel) GetResult(id int) *md.ResultModel {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	rs := m.Results[id]
	return rs
}

func (m *UserModel) AddResult(result *md.ResultModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.Results[result.ID] = result
	m.saveUser(path.Join(C.UserPath, fmt.Sprintf("%d.json", m.ID)))
}

func (m *UserModel) GetResultImage(id int, index int) ([]byte, error) {
	fileName := path.Join(C.StorePath, fmt.Sprintf("%d/%d/loglog/%d.png", m.ID, id, index))
	return os.ReadFile(fileName)
}

func (m *UserModel) GetResultMFS(id int) ([]byte, error) {
	fileName := path.Join(C.StorePath, fmt.Sprintf("%d/%d/mfs.bin", m.ID, id))
	return os.ReadFile(fileName)
}

type ImageModel struct {
	Img      image.Image `json:"-"`
	ID       int         `json:"id"`
	FileName string      `json:"fileName"`
	IsSaved  bool        `json:"isSaved"`
}
