package model

import (
	"fmt"
	"fractalmri/gofract/pipes"
)

var M = &AppModel{
	Profile:       Profile{},
	WindowRatio:   5,
	MinUmbral:     0.1,
	MaxUmbral:     1.4,
	Step:          0.1,
	IsUploading:   false,
	Images:        map[int]*ImageModel{},
	Selection:     map[int]*ImageModel{},
	LastSelection: -1,
	Results:       map[int]*ResultModel{},
	CurrentResult: -1,
	SegToolMin:    0.9,
	SegToolMax:    1.1,
}

type AppModel struct {
	Profile         Profile
	WindowRatio     int
	MinUmbral       float64
	MaxUmbral       float64
	Step            float64
	IsUploading     bool
	AllowRetry      bool
	UploadErr       string
	Images          map[int]*ImageModel
	Selection       map[int]*ImageModel
	LastSelection   int
	Results         map[int]*ResultModel
	CurrentResult   int
	ShowResult      bool
	HasSgTool       bool
	RunErr          string
	IsWaitingResult bool
	SegToolMin      float64
	SegToolMax      float64
}

func (m *AppModel) GetCurrentResult() *ResultModel {
	return m.Results[m.CurrentResult]
}

type ImageModel struct {
	URL     string
	Caption string
	ID      int
}

func ImageUrlForID(ID int) string {
	return fmt.Sprintf("/image?id=%d", ID)
}

func NewImageModel(ID int) *ImageModel {
	return &ImageModel{
		ID:      ID,
		URL:     ImageUrlForID(ID),
		Caption: fmt.Sprintf("Imagen %d", ID+1),
	}
}

type ResultModel struct {
	ID          int                 `json:"id"`
	FDs         []*pipes.FractalDim `json:"fds"`
	Umbrals     []*pipes.Umbral     `json:"umbrals"`
	LogLogCount int                 `json:"log-log-count"`
	IsCompleted bool                `json:"is-completed"`
}

func ResultURLForID(ID int) string {
	return fmt.Sprintf("/result?id=%d", ID)
}

func NewResultModel(ID int) *ResultModel {
	return &ResultModel{
		ID:      ID,
		FDs:     make([]*pipes.FractalDim, 0),
		Umbrals: make([]*pipes.Umbral, 0),
	}
}

type RunModel struct {
	ImageID    int     `json:"image-id"`
	Window     int     `json:"window"`
	MinUmbral  float64 `json:"min-umbral"`
	MaxUmbral  float64 `json:"max-umbral"`
	StepUmbral float64 `json:"step-umbral"`
}

type SegModel struct {
	ID        int     `json:"id"`
	MinUmbral float64 `json:"min-umbral"`
	MaxUmbral float64 `json:"max-umbral"`
}

type Profile struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}
