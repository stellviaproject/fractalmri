package lib

import (
	"image"
	"image/png"
	"os"
)

func ReadImage(fileName string) (*image.RGBA, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img.(*image.RGBA), nil
}

func WriteImage(fileName string, img *image.RGBA) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm|os.ModeDevice)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}
