package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"image"

	arch "github.com/stellviaproject/pipfil-arch"
)

func NewObjectFilter(cfg cfg.Configuration, binInPipe, objOutPipe arch.Pipe) arch.Filter {
	return arch.NewFilterWithPipes(
		"object",
		func(bin pipes.ImageBin) []*pipes.ImageObject {
			ratio := cfg.Ratio
			minArea := cfg.MinArea
			maxArea := cfg.MaxArea

			// Crear una copia de la imagen para marcar los píxeles visitados
			visited := pipes.NewImageBin(bin.Width(), bin.Height())

			// Definir la lista de objetos encontrados
			var objects []*pipes.ImageObject

			// Recorrer todos los píxeles de la imagen
			for y := 0; y < bin.Height(); y++ {
				for x := 0; x < bin.Width(); x++ {
					// Si el pixel ya ha sido visitado o no es parte de un objeto, saltar al siguiente pixel
					if visited.At(x, y) || !bin.At(x, y) {
						continue
					}

					// Crear un nuevo objeto para almacenar los píxeles conectados
					obj := pipes.NewImageBin(bin.Width(), bin.Height())
					area := 0

					// Definir la cola de píxeles adyacentes que se deben visitar
					queue := []image.Point{{x, y}}

					// Mientras la cola no esté vacía, visitar todos los píxeles adyacentes y agregarlos al objeto
					for len(queue) > 0 {
						// Obtener el siguiente pixel de la cola
						p := queue[0]
						queue = queue[1:]

						// Si el pixel ya ha sido visitado saltar al siguiente pixel
						if visited.At(p.X, p.Y) {
							continue
						}

						// Marcar el pixel como visitado y agregarlo al objeto
						visited.Set(p.X, p.Y, true)
						// Si no es parte de un objeto, saltar al siguiente pixel
						if !bin.At(p.X, p.Y) {
							continue
						}
						obj.Set(p.X, p.Y, true)
						area++

						// Agregar todos los píxeles adyacentes que estén dentro del ratio de distancia especificado a la cola
						for dy := -ratio; dy <= ratio; dy++ {
							for dx := -ratio; dx <= ratio; dx++ {
								// Saltar los píxeles diagonales
								if dx == 0 && dy == 0 {
									continue
								}

								// Obtener la coordenada del pixel adyacente
								px := p.X + dx
								py := p.Y + dy

								// Si el pixel adyacente está dentro de la imagen y no ha sido visitado, agregarlo a la cola
								if px >= 0 && px < bin.Width() && py >= 0 && py < bin.Height() && !visited.At(px, py) {
									queue = append(queue, image.Point{px, py})
								}
							}
						}
					}
					if area >= minArea && area <= maxArea {
						// Agregar el objeto encontrado a la lista de objetos
						objects = append(objects, &pipes.ImageObject{
							Image: obj,
							Area:  area,
						})
					}
				}
			}
			return objects
		},
		arch.WithPipes(binInPipe),
		arch.WithPipes(objOutPipe),
		arch.WithLens(),
	)
}
