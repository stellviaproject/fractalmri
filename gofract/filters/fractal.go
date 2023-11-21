package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"math"
	"reflect"

	"github.com/openacid/slimarray/polyfit"
	psync "github.com/stellviaproject/image-psync"
	arch "github.com/stellviaproject/pipfil-arch"
)

// Un filtro para calcular la dimension fractal de una imagen binaria
func NewFractalFilter(cfg cfg.Configuration, objInPipe, fdOutPipe arch.Pipe) arch.Filter {
	var fn any //La funcion del filtro
	//Si el tipo de dato de entrada de la tuberia es una imagen binaria
	if objInPipe.CheckType() == reflect.TypeOf(pipes.ImageBin{}) {
		fn = fractalFn(cfg) //Obtener la funcion del filtro
	} else if objInPipe.CheckType() == reflect.TypeOf(&pipes.ImageObject{}) { //Sino si es un objeto de una imagen
		//Crear una funcion que obtenga la imagen del objeto
		fn = func(obj *pipes.ImageObject) *pipes.FractalDim {
			//Obtener la imagen del objeto y llamar a la funcion de la dimension fractal
			return fractalFn(cfg)(obj.Image)
		}
	} else {
		//El tipo de dato de la tuberia de entrada no es el esperado
		panic("unknown type for fractal input")
	}
	//Crear el filtro
	return arch.NewFilterWithPipes(
		"fractal",                 //El nombre del filtro
		fn,                        //La funcion del filtro
		arch.WithPipes(objInPipe), //La tuberia de entrada del filtro
		arch.WithPipes(fdOutPipe), //La tuberia de salida del filtro
		arch.WithLens(),
	)
}

// Retorna una funcion que calcula la dimension fractal de una imagen binaria
func fractalFn(cfg cfg.Configuration) func(bin pipes.ImageBin) *pipes.FractalDim {
	//La funcion para calcular la dimension fractal
	return func(bin pipes.ImageBin) *pipes.FractalDim {
		//Obtener la configuracion del filtro
		boxSizes := cfg.BoxSizes                 //Tamanio de las cajas para cada conteo
		parallel := cfg.Parallel                 //Nivel de paralelismo del filtro
		counts := make([]float64, len(boxSizes)) //Numero de cajas por cada tamanio de caja

		width, height := bin.Width(), bin.Height() //Obtener las dimensiones de la imagen
		//Realizar los conteo de cajas
		for i, window := range boxSizes {
			count := 0               //La cantidad para la caja de tamanio window es 0 inicialmente
			mtx := make(chan int, 1) //Crear un canal que sirva de mutex
			// Iterar sobre los cuadrados de tamaño s y contar los que contienen al menos un píxel negro
			psync.ParallelWindow(width, height, window, parallel, func(minX, minY, maxX, maxY int) {
				//Recorrer la ventana para ver si tiene un objeto
				for x := minX; x < maxX; x++ { //Iterar en x
					for y := minY; y < maxY; y++ { //Iterar en y
						if bin.At(x, y) { //Si hay un pixel en la coordenada (x,y)
							mtx <- 0 //Bloquear la gorutine
							count++  //Incrementar en uno el contador
							<-mtx    //Desbloquear la gorutine
							return
						}
					}
				}
			})
			counts[i] = float64(count) //Establecer la cantidad de cajas en el slice de cantidades
		}
		//Obtener los logaritmos de las cantidades
		logCounts := make([]float64, len(counts))
		for i := 0; i < len(counts); i++ {
			logCounts[i] = math.Log(counts[i])
		}
		//Obtener los logaritmos de los tamanios
		logSizes := make([]float64, len(boxSizes))
		for i := 0; i < len(boxSizes); i++ {
			logSizes[i] = math.Log(1.0 / float64(boxSizes[i]))
		}
		// Realizar la interpolación polinomial y calcular los logaritmos de la cantidad sobre la longitud
		fit := polyfit.NewFit(logSizes, logCounts, 1.0).Solve()
		//Retornar el objeto que representa la dimension fractal
		return pipes.NewFD(logSizes, logCounts, fit[1])
	}
}
