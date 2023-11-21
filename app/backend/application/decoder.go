package application

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/tiff"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func Decode(data []byte) ([]image.Image, error) {
	reader := bytes.NewReader(data)
	var img image.Image
	var err error
	img, err = tiff.Decode(reader)
	if err == nil {
		return []image.Image{img}, nil
	}
	reader.Seek(0, io.SeekStart)
	img, err = png.Decode(reader)
	if err == nil {
		return []image.Image{img}, nil
	}
	reader.Seek(0, io.SeekStart)
	img, err = jpeg.Decode(reader)
	if err == nil {
		return []image.Image{img}, nil
	}
	reader.Seek(0, io.SeekStart)
	img, _, err = image.Decode(reader)
	if err == nil {
		return []image.Image{img}, nil
	}
	reader.Seek(0, io.SeekStart)
	dataset, err := dicom.Parse(reader, int64(len(data)), nil)
	if err != nil {
		return nil, err
	}
	pixelDataElement, _ := dataset.FindElementByTag(tag.PixelData)
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	imgls := []image.Image{}
	for _, fr := range pixelDataInfo.Frames {
		img, err := fr.GetImage()
		if err == nil {
			imgls = append(imgls, img)
		}
	}
	return imgls, nil
}
