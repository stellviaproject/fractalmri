package config

import (
	"bytes"
	"encoding/json"
	"image"
	"io"
	"log"
	"os"
	"time"
)

type WithConfigFn func(c *Configuration)

func WithQuality(quality int) WithConfigFn {
	if quality < 1 {
		panic("quality must be greater than 1")
	}
	return func(c *Configuration) {
		c.Quality = quality
	}
}

func WithMaxSize(maxSize int64) WithConfigFn {
	if maxSize < 1024*1024 {
		panic("upload size must be greater than 1 mega-bytes or 1048576 bytes")
	}
	return func(c *Configuration) {
		c.UploadMaxSize = maxSize
	}
}

func WithWebDir(webDir string) WithConfigFn {
	if _, err := os.Stat(webDir); err != nil {
		log.Panicln(err)
	}
	return func(c *Configuration) {
		c.WebDir = webDir
	}
}

type Configuration struct {
	UploadMaxSize int64         `json:"upload_max_size"`
	WebDir        string        `json:"-"`
	Quality       int           `json:"quality"`
	LogPath       string        `json:"log_path"`
	Port          int           `json:"port"`
	MainFile      string        `json:"model"`
	UserPath      string        `json:"users"`
	StorePath     string        `json:"store"`
	MaxTime       time.Duration `json:"max-time"`
}

func Default() *Configuration {
	return &Configuration{
		Port:          8080,
		UploadMaxSize: 20 * 1024 * 1024,
		WebDir:        "./web",
		Quality:       90,
		MainFile:      "./model.json",
		UserPath:      "./users",
		StorePath:     "./store",
		MaxTime:       time.Minute * 30,
	}
}

func Test(img image.Image) image.Image {
	return img
}

func (c *Configuration) test() {
	if _, err := os.Stat(c.WebDir); err != nil {
		log.Fatalln("web directory not found")
	}
}

func LoadConfig(filePath string) (*Configuration, error) {
	configuration := &Configuration{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, configuration); err != nil {
		return nil, err
	}
	configuration.WebDir = "./web"
	configuration.test()
	return configuration, nil
}

func (c *Configuration) SaveConfig() error {
	file, err := os.Create("./config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err = json.Indent(&buffer, data, "", "\t"); err != nil {
		return err
	}
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}
