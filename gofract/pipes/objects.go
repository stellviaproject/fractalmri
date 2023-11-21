package pipes

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type FractalDim struct {
	LogSizes   []float64 `json:"log-sizes"`
	LogMeasure []float64 `json:"log-measure"`
	FD         float64   `json:"fd"`
}

func NewFD(logSizes, logMeasures []float64, fd float64) *FractalDim {
	return &FractalDim{
		LogSizes:   logSizes,
		LogMeasure: logMeasures,
		FD:         fd,
	}
}

type Image64 [][]float64

func NewImage64(width, height int) Image64 {
	img := make([][]float64, height)
	for i := 0; i < height; i++ {
		img[i] = make([]float64, width)
	}
	return img
}

func ToImage64(img image.Image) Image64 {
	bd := img.Bounds()
	width, height := bd.Dx(), bd.Dy()
	i64 := NewImage64(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			i64.Set(x, y, float64(color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y))
		}
	}
	return i64
}

func (img Image64) Width() int {
	return len(img[0])
}

func (img Image64) Height() int {
	return len(img)
}

func (img Image64) Shape() (int, int) {
	return len(img[0]), len(img)
}

func (img Image64) At(x, y int) float64 {
	return img[y][x]
}

func (img Image64) Set(x, y int, value float64) {
	img[y][x] = value
}

func (img Image64) Normalized() Image64 {
	width, height := img.Width(), img.Height()
	norm := NewImage64(width, height)
	max := 0.0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			max = math.Max(max, img.At(x, y))
		}
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			norm.Set(x, y, img.At(x, y)/max)
		}
	}
	return norm
}

func (img Image64) Scaled(scalar float64) Image64 {
	width, height := img.Width(), img.Height()
	scaled := NewImage64(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			scaled.Set(x, y, img.At(x, y)*scalar)
		}
	}
	return scaled
}

func (img Image64) ToImage() image.Image {
	width, height := img.Width(), img.Height()
	gray := image.NewGray(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gray.SetGray(x, y, color.Gray{Y: uint8(img.At(x, y))})
		}
	}
	return gray
}

type ImageBin [][]bool

func NewImageBin(width int, height int) ImageBin {
	img := make([][]bool, height)
	for i := 0; i < height; i++ {
		img[i] = make([]bool, width)
	}
	return img
}

func (img ImageBin) Width() int {
	return len(img[0])
}

func (img ImageBin) Height() int {
	return len(img)
}

func (img ImageBin) At(x, y int) bool {
	return img[y][x]
}

func (img ImageBin) Set(x, y int, value bool) {
	img[y][x] = value
}

func (img ImageBin) Shape() (int, int) {
	return len(img[0]), len(img)
}

func (img ImageBin) ToImage() image.Image {
	width, height := img.Width(), img.Height()
	bin := image.NewGray(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if img.At(x, y) {
				bin.Set(x, y, color.Gray{Y: 255})
			}
		}
	}
	return bin
}

type ImageObject struct {
	Image ImageBin
	Area  int
}

type Umbral struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func NewUmbral(min, max float64) *Umbral {
	return &Umbral{
		Min: min,
		Max: max,
	}
}

func (um *Umbral) String() string {
	return fmt.Sprintf("(%f, %f)", um.Min, um.Max)
}

type FractalInfo struct {
	LogLogs []image.Image
	FDs     []*FractalDim
	Objects []*ImageObject
}

type ImageInfo struct {
	Umbrals       []*Umbral
	FractalInfoLs []*FractalInfo
}

type FractalAnalysisInfo struct {
	ImageList   []image.Image
	MFSs        []Image64
	ImageInfoLs []*ImageInfo
}
