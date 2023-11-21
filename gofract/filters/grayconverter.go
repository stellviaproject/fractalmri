package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"image"

	psync "github.com/stellviaproject/image-psync"

	arch "github.com/stellviaproject/pipfil-arch"
)

// Un filtro que convierte una imagen a escala de grises
func NewGrayConverterFilter(cfg cfg.Configuration, imageInPipe, grayImg64OutPipe arch.Pipe) arch.Filter {
	//Crear el filtro
	return arch.NewFilterWithPipes(
		"grayconverter", //Crear el convertidor
		func(img image.Image) pipes.Image64 {
			parallel := cfg.Parallel
			// Crear una nueva imagen en escala de grises
			bd := img.Bounds()
			gray := pipes.NewImage64(bd.Dx(), bd.Dy())
			// Copiar los valores de intensidad de cada píxel de la imagen original a la nueva imagen en escala de grises
			psync.ParallelRegionHorizontal(gray.Width(), gray.Height(), parallel, func(minX, minY, maxX, maxY int) {
				for x := minX; x < maxX; x++ {
					for y := minY; y < maxY; y++ {
						// Obtener el color del píxel en la imagen original
						r, g, b, _ := img.At(x, y).RGBA()

						// Calcular la intensity del píxel en la nueva imagen en escala de grises
						intensity := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0
						//Establecer el valor del gris en la imagen de salida
						gray.Set(x, y, intensity)
					}
				}
			})
			return gray
		},
		arch.WithPipes(imageInPipe),      //Tuberia de entrada
		arch.WithPipes(grayImg64OutPipe), //Tuberia de salida
		arch.WithLens(),
	)
}
