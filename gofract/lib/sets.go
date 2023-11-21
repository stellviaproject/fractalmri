package lib

import "image"

type ItemOfSet interface {
	GetImage() (*image.RGBA, error)
}

type ImageSetItem struct {
	Image    *image.RGBA
	MRILabel string
}

//export NewImageSetItem
func NewImageSetItem(img *image.RGBA, mriLabel string) *ImageSetItem {
	return &ImageSetItem{
		Image:    img,
		MRILabel: mriLabel,
	}
}

func (imgIt *ImageSetItem) GetImage() (*image.RGBA, error) {
	return imgIt.Image, nil
}

type FileSetItem struct {
	FileName string
	MRILabel string
}

//export NewFileSetItem
func NewFileSetItem(fileName string, mriLabel string) *FileSetItem {
	return &FileSetItem{
		FileName: fileName,
		MRILabel: mriLabel,
	}
}

func (fileIt *FileSetItem) GetImage() (*image.RGBA, error) {
	return ReadImage(fileIt.FileName)
}
