package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"sync"

	psync "github.com/stellviaproject/image-psync"
	arch "github.com/stellviaproject/pipfil-arch"
)

func NewUmbralizerFilter(cfg cfg.Configuration, image64InPipe, umbralImg64OutPipe, umbralListOutPipe arch.Pipe) arch.Filter {
	return arch.NewFilterWithPipes(
		"umbralizer",
		func(img pipes.Image64) ([]pipes.Image64, []*pipes.Umbral) {
			//Obtener los elementos de la configuracion
			umbral := cfg.Umbral                              //Obtener la secuencia de umbrales
			umbralLs := make([]*pipes.Umbral, 0, len(umbral)) //Generar la lista de umbrales maximos y mínimos
			for i := 1; i < len(umbral); i++ {
				umbralLs = append(umbralLs, pipes.NewUmbral(umbral[i-1], umbral[i]))
			}
			umbralLs = append(umbralLs, pipes.NewUmbral(cfg.MinUmbral, cfg.MaxUmbral)) //Agregar el umbral especifico

			outputs := make([]pipes.Image64, len(umbralLs)) //Perparar la lista de imagenes de salida
			//Inicializar variables de control de las gorutines
			wg := sync.WaitGroup{}               //grupo de espera
			prll := make(chan int, cfg.Parallel) //cantidad de gorutines

			//Segmentar por cada umbral
			for i := 0; i < len(umbralLs); i++ {
				//Indicar que se inicia una nueva gorutine
				wg.Add(1)
				prll <- 0
				//Iniciar gorutine
				go func(i int, minUmbral, maxUmbral float64) {
					defer wg.Done() //Notificar cuando se termine la gorutine
					//Crear la imagen de salida
					output := pipes.NewImage64(img.Width(), img.Height())
					//Recorrer la imagen por regiones horizontales en paralelo
					psync.ParallelRegionHorizontal(img.Width(), img.Height(), cfg.Parallel, func(minX, minY, maxX, maxY int) {
						//Recorrer la region de la imagen
						for x := minX; x < maxX; x++ {
							for y := minY; y < maxY; y++ {
								v := img.At(x, y) //Obtener el pixel en la posicion (x,y)
								//Comprobar si esta en el rango de umbral
								if v >= minUmbral && v <= maxUmbral {
									output.Set(x, y, v) //Establecer el pixel
								}
							}
						}
					})
					outputs[i] = output //Establecer la salida en la lista
					<-prll              //Indicar que se puede continuar con la siguiente
				}(i, umbralLs[i].Min, umbralLs[i].Max)
			}
			wg.Wait()                //Esperar porque terminen las gorutines
			return outputs, umbralLs //Retornar imagenes y lista de umbrales
		},
		arch.WithPipes(image64InPipe),
		arch.WithPipes(umbralImg64OutPipe, umbralListOutPipe),
		arch.WithLens())
}
