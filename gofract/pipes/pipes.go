package pipes

import (
	"fractalmri/gofract/cfg"
	"image"

	arch "github.com/stellviaproject/pipfil-arch"
)

func NewImagePipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, (*image.Image)(nil), cfg.Buffer)
}

func NewImage64Pipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, Image64{}, cfg.Buffer)
}

func NewUmbralListPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, []*Umbral{}, cfg.Buffer)
}

func NewImageBinPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, ImageBin{}, cfg.Buffer)
}

func NewFractalPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, &FractalDim{}, cfg.Buffer)
}

func NewBufferPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, []byte{}, cfg.Buffer)
}

func NewObjectPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, &ImageObject{}, cfg.Buffer)
}

func NewFractalInfoPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, &FractalInfo{}, cfg.Buffer)
}

func NewImageInfoPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, &ImageInfo{}, cfg.Buffer)
}

func NewFractalAnalysisInfoPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, &FractalAnalysisInfo{}, cfg.Buffer)
}

func NewLacPipe(name string, cfg cfg.Configuration) arch.Pipe {
	return arch.NewPipe(name, []float64{}, cfg.Buffer)
}
