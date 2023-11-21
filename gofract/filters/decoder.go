package filters

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	arch "github.com/stellviaproject/pipfil-arch"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"golang.org/x/image/tiff"
)

/*
Este es el filtro decodificador.

Recibe los bytes de una imagen y obtiene una imagen.

inputPipe debe ser creada con NewBufferPipe.

imagePipe debe ser creada con NewImagePipe.

Si se usa una imagen en formato DICOM con varias imagenes dentro, se debe agregar un filtro
en la salida para unir los resultados de cada imagen con tal de evitar un deadlock.
*/
func NewDecoderFilter(inputPipe, imagePipe arch.Pipe) arch.Filter {
	//Crea el filtro.
	return arch.NewFilterWithPipes(
		"decoder", //El nombre del filtro
		func(input []byte) ([]image.Image, error) { //La funcion del filtro
			reader := bytes.NewReader(input) //Crear un lector con los bytes de la entrada
			img, err := tiff.Decode(reader)  //Decodificar la imagen en tiff
			if err == nil {                  //Si no hay error retornar la imagen
				return []image.Image{img}, nil
			}
			//Si hay error es posible que el formato no sea tiff
			reader.Seek(0, io.SeekStart)  //Mover la posicion del lector a 0
			img, err = png.Decode(reader) //Decodificar como png
			if err == nil {               //Si no hay error retornar la imagen
				return []image.Image{img}, nil
			}
			//Si hay error es posible que el formato no sea png
			reader.Seek(0, io.SeekStart)   //Mover la posicion del lector a 0
			img, err = jpeg.Decode(reader) //Decodificar como jpg
			if err == nil {                //Si no hay error retornar la imagen
				return []image.Image{img}, nil
			}
			reader.Seek(0, io.SeekStart)       //Mover la posicion del lector a 0
			img, _, err = image.Decode(reader) //Decodificar como cualquier otro tipo de imagen que acepte el paquete image
			if err == nil {                    //Si no hay error retornar la imagen
				return []image.Image{img}, nil
			}
			//Si hay error es posible que el formato sea DICOM
			reader.Seek(0, io.SeekStart)                                  //Regresar la posicion del lector al inicio del buffer
			dataset, err := dicom.Parse(reader, int64(reader.Len()), nil) //Crear el parseador de DICOM
			if err != nil {                                               //Si hay error ya no se sigue decodificando
				return nil, err
			}
			//Encontrar los pixeles de la imagen
			pixelDataElement, _ := dataset.FindElementByTag(tag.PixelData)
			pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
			//Slice para guardar las imagenes
			imgls := []image.Image{}
			//Recorrer los frames de imagenes en el archivo DICOM
			for _, fr := range pixelDataInfo.Frames {
				img, err := fr.GetImage() //Obtener la imagen de un frame
				if err == nil {           //Si no hay error
					imgls = append(imgls, img) //Agregar la imagen al slice
				}
			}
			//Retornar las imagenes del archivo DICOM
			return imgls, err
		},
		arch.WithPipes(inputPipe), //La tuberia de entrada del filtro
		arch.WithPipes(imagePipe), //La tuberia de salida del filtro
		arch.WithLens(),
	)
}
