package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"math"

	arch "github.com/stellviaproject/pipfil-arch"
)

func NewMultiFractal(cfg cfg.Configuration, image64InPipe, mfsImage64OutPipe arch.Pipe) arch.Filter {
	return arch.NewFilterWithPipes(
		"multifractal",
		func(measure pipes.Image64) pipes.Image64 {
			//Obtener el tamaño de ventana para el espectro MFS
			ratio := cfg.WindowRatio
			//Calcular logaritmo del tamaño de ventana
			wLog := math.Log(float64(2 * ratio))
			//Obtener ancho y alto de la imagen de medida
			width, height := measure.Width(), measure.Height()
			//Crear imagen de salida del espectro multifractal
			mfs := pipes.NewImage64(width, height)
			//Recorrer la imagen de medida
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					//Calcular exponente de hölder
					mfs.Set(x, y, math.Log(measure.At(x, y))/wLog)
				}
			}
			//Retornar espectro
			return mfs
		},
		arch.WithPipes(image64InPipe),
		arch.WithPipes(mfsImage64OutPipe),
		arch.WithLens(),
	)
}
