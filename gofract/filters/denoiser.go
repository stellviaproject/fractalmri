package filters

import (
	"fractalmri/gofract/cfg"
	"image"
	"image/color"
	"math"

	psync "github.com/stellviaproject/image-psync"

	arch "github.com/stellviaproject/pipfil-arch"
)

// El filtro bilateral para reducir y eliminar el ruido de la imagen
func NewDenoiserFilter(cfg cfg.Configuration, imageInPipe, imageOutPipe arch.Pipe) arch.Filter {
	//Crear el filtro
	return arch.NewFilterWithPipes(
		"denoiser", //El nombre del filtro
		func(img image.Image) image.Image { //La funcion del filtro
			//Obtener la configuracion del filtro
			parallel := cfg.Parallel             //Nivel de paralelismo
			diameter := cfg.DenoiserDiameter     //Diametro del filtro bilateral
			sigmaColor := cfg.DenoiserSigmaColor //SigmaColor del filtro bilateral
			sigmaSpace := cfg.DenoiserSigmaSpace //SigmaSpace del filtro bilateral
			minColor := cfg.DenoiserUmbralColor  //Un umbral minimo para eliminar el ruido
			// Estructura para representar un píxel en la imagen
			type Pixel struct {
				R, G, B uint8
			}
			// Obtener las dimensiones de la imagen
			bounds := img.Bounds()
			width := bounds.Max.X
			height := bounds.Max.Y

			// Crear una nueva imagen para el resultado
			output := image.NewRGBA(bounds)

			// Generar una matriz de pesos para el filtro
			weights := make([][]float64, diameter)
			for i := 0; i < diameter; i++ {
				weights[i] = make([]float64, diameter)
			}
			for i := 0; i < diameter; i++ {
				for j := 0; j < diameter; j++ {
					distance := math.Sqrt(float64((i-diameter/2)*(i-diameter/2) + (j-diameter/2)*(j-diameter/2)))
					weights[i][j] = math.Exp(-(distance * distance) / (2 * sigmaSpace * sigmaSpace))
				}
			}
			// Recorrer todos los píxeles de la imagen en paralelo
			psync.ParallelForEach(width, height, parallel, func(x, y, maxX, maxY int) {
				// Obtener el valor de color del píxel
				r, g, b, a := img.At(x, y).RGBA()

				pixel := Pixel{uint8(r / 256), uint8(g / 256), uint8(b / 256)}
				// Calcular el nuevo valor de color del píxel
				sumRw, sumGw, sumBw, sumWeights := 0.0, 0.0, 0.0, 0.0
				sumR, sumG, sumB := 0.0, 0.0, 0.0
				for i := 0; i < diameter; i++ {
					for j := 0; j < diameter; j++ {
						// Obtener el valor de color del píxel vecino
						nx := x + i - diameter/2
						ny := y + j - diameter/2
						if nx < 0 || nx >= width || ny < 0 || ny >= height {
							continue
						}
						r2, g2, b2, _ := img.At(nx, ny).RGBA()
						pixel2 := Pixel{uint8(r2 / 256), uint8(g2 / 256), uint8(b2 / 256)}

						// Calcular el peso del píxel vecino
						dR := int(pixel.R) - int(pixel2.R)
						dG := int(pixel.G) - int(pixel2.G)
						dB := int(pixel.B) - int(pixel2.B)
						colorDiff := math.Sqrt(float64(dR*dR + dG*dG + dB*dB))
						weight := math.Exp(-(colorDiff*colorDiff)/(2*sigmaColor*sigmaColor)) * weights[i][j]

						// Acumular los valores ponderados
						sumRw += float64(pixel2.R) * weight
						sumGw += float64(pixel2.G) * weight
						sumBw += float64(pixel2.B) * weight

						sumR += float64(r2)
						sumG += float64(g2)
						sumB += float64(b2)

						sumWeights += weight
					}
				}
				newR := uint8(math.Round(sumRw / sumWeights))
				newG := uint8(math.Round(sumGw / sumWeights))
				newB := uint8(math.Round(sumBw / sumWeights))

				intensity := (0.299*float64(sumR) + 0.587*float64(sumG) + 0.114*float64(sumB)) / (256.0 * float64(diameter) * float64(diameter))
				if intensity <= minColor {
					output.Set(x, y, color.RGBA{0, 0, 0, 255})
				} else {
					// Establecer el nuevo valor de color del píxel en la imagen de salida
					output.Set(x, y, color.RGBA{newR, newG, newB, uint8(a / 256)})
				}
			})
			return output
		},
		arch.WithPipes(imageInPipe),  //La tuberia de entrada
		arch.WithPipes(imageOutPipe), //La tuberia de salida
		arch.WithLens(),
	)
}
