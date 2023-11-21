package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"

	psync "github.com/stellviaproject/image-psync"
	arch "github.com/stellviaproject/pipfil-arch"
)

func NewMeasureFilter(cfg cfg.Configuration, image64InPipe, image64OutPipe arch.Pipe) arch.Filter {
	return arch.NewFilterWithPipes(
		"measure",
		func(img pipes.Image64) pipes.Image64 {
			//Obtener elementos de configuracion
			ratio := cfg.WindowRatio //Tamaño de ventana
			parallel := cfg.Parallel //Cantidad de gorutines
			//Obtener el ancho y alto de la imagen
			w, h := img.Width(), img.Height()
			//Crear la imagen de medida
			measure := pipes.NewImage64(w, h)
			//Recorrer cada pixel de la imagen en pararelo segun el numero de gorutines permitidas
			psync.ParallelForEach(img.Width(), img.Height(), parallel, func(minX, minY, maxX, maxY int) {
				sum := 0.0 //Variable para guardar la suma
				//Recorrer la imagen en las coordenadas (x,y)
				//para sumar los pixeles de los alrededores
				for x := -ratio; x <= ratio; x++ {
					for y := -ratio; y <= ratio; y++ {
						//Obtener la posicion en la imagen
						xp, yp := x+minX, y+minY
						//Si el pixel esta fuera de la region
						if xp < 0 || yp < 0 || xp >= w || yp >= h {
							continue //ignorar la coordenada
						}
						//Obtener pixel y sumar
						sum += img.At(xp, yp) / 255.0
					}
				}
				//Establecer la suma para la region actual
				measure.Set(minX, minY, sum)
			})
			return measure
		},
		arch.WithPipes(image64InPipe),
		arch.WithPipes(image64OutPipe),
		arch.WithLens(),
	)
}
