package main

import (
	"bytes"
	"flag"
	"fmt"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/lib"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime/pprof"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"golang.org/x/image/tiff"
)

func main() {
	d := flag.String("d", "", "directory with a folder called tumors and a folder called notumors.")
	o := flag.String("o", "", "output directory for saveing exel files")
	flag.Parse()
	ratios := []int{2}
	tumors, notumors := LoadAll(*d)
	Dis(tumors)
	Dis(notumors)
	all := append(tumors, notumors...)
	all = all[:100]
	folderOutput := path.Join(*o, "cpu-profile-")
	for i := 0; i < len(ratios); i++ {
		filePts, err := os.Create(folderOutput + fmt.Sprintf("evalpoints-%d.prof", ratios[i]))
		if err != nil {
			log.Fatalln(err)
		}
		if err := pprof.StartCPUProfile(filePts); err != nil {
			log.Fatalln(err)
		}
		config := cfg.GetCFG()
		config.WindowRatio = ratios[i]
		pts := EvalPoints(config, all)
		pprof.StopCPUProfile()
		filePts.Close()
		fileKNN, err := os.Create(folderOutput + fmt.Sprintf("evalknn-%d.prof", ratios[i]))
		if err != nil {
			log.Fatalln(err)
		}
		if err := pprof.StartCPUProfile(fileKNN); err != nil {
			log.Fatalln(err)
		}
		EvaluateKNN(config, pts)
		pprof.StopCPUProfile()
		fileKNN.Close()
	}
}

func Dis(set []string) {
	for i := range set {
		j := rand.Intn(len(set))
		set[i], set[j] = set[j], set[i]
	}
}

func LoadAll(dir string) (tumors, notumors []string) {
	tumors = LoadFiles(path.Join(dir, "tumors"))
	notumors = LoadFiles(path.Join(dir, "notumors"))
	return
}

func EvaluateKNN(config cfg.Configuration, all []*lib.DataPoint) {
	knn := lib.NewKNNFractal(config, all)
	evaluate := func(test []*lib.DataPoint) (fails int) {
		for i := 0; i < len(test); i++ {
			result := knn.FitPoint(test[i].Point())
			if result != test[i].MRILabel {
				fails++
			}
		}
		return
	}
	evaluate(all)
}

func LoadFiles(dir string) []string {
	ens, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	files := []string{}
	for i := 0; i < len(ens); i++ {
		log.Println("loading files ", i, "...")
		if !strings.HasSuffix(ens[i].Name(), "-mask.png") {
			files = append(files, path.Join(dir, ens[i].Name()))
		}
	}
	log.Println("loading files ", len(files), "...")
	return files
}

func EvalPoints(config cfg.Configuration, all []string) []*lib.DataPoint {
	model := lib.NewFDModelNoLogLog(config)
	pts := make([]*lib.DataPoint, 0, 100)
	for i := 0; i < len(all); i++ {
		data, err := os.ReadFile(all[i])
		if err != nil {
			log.Fatalln(err)
		}
		imgs, err := ReadImage(data)
		if err != nil {
			log.Fatalln(err)
		}
		for _, img := range imgs {
			pt, err := model.Eval(img.(*image.RGBA))
			if err != nil {
				log.Fatalln(err)
			}
			pts = append(pts, &lib.DataPoint{FD: pt.GetPoints(), MRILabel: "none"})
		}
	}
	return pts
}

func ReadImage(input []byte) ([]image.Image, error) { //La funcion del filtro
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
}
