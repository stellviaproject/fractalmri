package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"

	psync "github.com/stellviaproject/image-psync"
	arch "github.com/stellviaproject/pipfil-arch"
)

// Este filtro convierte una imagen Image64 (de valores float64) a una ImageBin (imagen binaria o de valores bool)
func NewBinarizeFilter(cfg cfg.Configuration, image64InPipe, imageBinOutPipe arch.Pipe) arch.Filter {
	//Crear el filtro
	return arch.NewFilterWithPipes(
		"binarizer", //El nombre del filtro
		func(img pipes.Image64) pipes.ImageBin { //La funcion que procesa los datos en el filtro
			parallel := cfg.Parallel                            //Obtener el nivel de paralelismo
			bin := pipes.NewImageBin(img.Width(), img.Height()) //Crear la imagen binaria
			//Recorrer la imagen en paralelo
			psync.ParallelRegionHorizontal(img.Width(), img.Height(), parallel, func(minX, minY, maxX, maxY int) {
				//Copiar una region de la imagen
				for x := minX; x < maxX; x++ { //Recorrer la coordenada x
					for y := minY; y < maxY; y++ { //Recorrer la coordeana y
						bin.Set(x, y, img.At(x, y) != 0.0) //Asignar True si es distinto de 0.0 (si el pixel no es vacio)
					}
				}
			})
			return bin //Retornar la imagen binaria
		},
		arch.WithPipes(image64InPipe),   // La tuberia de entrada Image64
		arch.WithPipes(imageBinOutPipe), // La tuberia de salida ImageBin
		arch.WithLens(),
	)
}


