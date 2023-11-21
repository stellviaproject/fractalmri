package main

import (
	"bytes"
	"encoding/json"
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
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"golang.org/x/image/tiff"
)

func main() {
	d := flag.String("d", "", "directory with a folder called tumors and a folder called notumors.")
	o := flag.String("o", "", "output directory for saveing exel files")
	s := flag.String("s", "", "load sample.json with files for testing")
	flag.Parse()
	ratios := []int{2, 5, 10, 15, 20, 25, 30, 35, 40}
	var tumors, notumors []string
	if *s == "" {
		tumors, notumors = LoadAll(*d)
	} else {
		sample, err := lib.LoadSample(*s)
		if err != nil {
			log.Fatalln(err)
		}
		tumors, notumors, testTumors, testNoTumors := sample.Tumors, sample.NoTumors, sample.TumorsTest, sample.NoTumorsTest
		evaluateWithCfg := func(config cfg.Configuration, tumors, notumors, testtumors, testnotumors []string) *PresitionEvaluation {
			model := lib.NewFDModelNoLogLog(config)
			evaluate := func(pts *[]*lib.DataPoint, msg, label string, files []string) {
				for i := 0; i < len(files); i++ {
					data, err := os.ReadFile(files[i])
					if err != nil {
						log.Fatalln(err)
					}
					imgs, err := ReadImage(data)
					if err != nil {
						log.Fatalln(err)
					}
					for k := 0; k < len(imgs); k++ {
						ev, err := model.Eval(imgs[k].(*image.RGBA))
						if err != nil {
							log.Fatalln(err)
						}
						*pts = append(*pts, &lib.DataPoint{
							MRILabel: label,
							FD:       ev.GetPoints(),
						})
					}
					log.Printf("evaluating-%s %d/%d %.1f %%", msg, i, len(files), 100.0*float64(i)/float64(len(files)))
				}
			}
			pev := &PartitionEvaluation{}
			evaluate(&pev.Tumors, "tumors", "tumors", tumors)
			evaluate(&pev.NoTumors, "notumors", "notumors", notumors)
			evaluate(&pev.Tumors, "testtumors", "tumors", tumors)
			evaluate(&pev.NoTumors, "testnotumors", "notumors", notumors)
			return pev.EvaluatePresition(config)
		}
		for i := 0; i < len(ratios); i++ {
			ratio := ratios[i]
			config := cfg.GetCFG()
			config.WindowRatio = ratio
			pres := evaluateWithCfg(config, tumors, notumors, testTumors, testNoTumors)
			file, err := os.Create(path.Join(*o, fmt.Sprintf("window-ratio-optimize-%d.json", ratio)))
			if err != nil {
				log.Fatalln(err)
			}
			data, err := json.MarshalIndent(pres, "", "    ")
			if err != nil {
				log.Fatalln(err)
			}
			_, err = file.Write(data)
			if err != nil {
				log.Fatalln(err)
			}
			file.Close()
		}
		return
	}
	Dis(tumors)
	Dis(notumors)
	for i := 0; i < len(ratios); i++ {
		EvalWindowRatio(ratios[i], *o, tumors, notumors)
	}
}

func Dis(set []string) {
	for i := range set {
		j := rand.Intn(len(set))
		set[i], set[j] = set[j], set[i]
	}
}

func EvalWindowRatio(ratio int, out string, tumors, notumors []string) {
	config := cfg.GetCFG()
	config.WindowRatio = ratio
	pres := EvalForConfig(config, tumors, notumors)
	file, err := os.Create(path.Join(out, fmt.Sprintf("window-ratio-%d.json", ratio)))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	data, err := json.MarshalIndent(pres, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
}

func LoadAll(dir string) (tumors, notumors []string) {
	tumors = LoadFiles(path.Join(dir, "tumors"))
	notumors = LoadFiles(path.Join(dir, "notumors"))
	return
}

func EvalForConfig(config cfg.Configuration, tumors, notumors []string) []*PresitionEvaluation {
	pte := GeneratePoints(config, tumors, notumors)
	partitions := []*PartitionEvaluation{
		pte.PartitionFiles(0.2),
		pte.PartitionFiles(0.3),
		pte.PartitionFiles(0.4),
		pte.PartitionFiles(0.5),
	}
	presitions := []*PresitionEvaluation{}
	for i := 0; i < len(partitions); i++ {
		presitions = append(presitions, partitions[i].EvaluatePresition(config))
	}
	return presitions
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

type PresitionEvaluation struct {
	FailsForTumors    int     `json:"failsForTumors"`
	FailsForNoTumors  int     `json:"failsForNoTumors"`
	TumorsAll         int     `json:"TumorsAll"`
	NoTumorsAll       int     `json:"NoTumorsAll"`
	PresitionTumors   float64 `json:"PresitionTumors"`
	PresitionNoTumors float64 `json:"PresitionNoTumors"`
	FailsAll          int     `json:"FailsAll"`
	PresitionAll      float64 `json:"PresitionAll"`
	All               int     `json:"All"`
}

type PartitionEvaluation struct {
	PointEvaluation
	TestTumors   []*lib.DataPoint
	TestNoTumors []*lib.DataPoint
}

func (pev *PartitionEvaluation) EvaluatePresition(config cfg.Configuration) *PresitionEvaluation {
	knn := lib.NewKNNFractal(config, append(pev.Tumors, pev.NoTumors...))
	evaluate := func(test []*lib.DataPoint) (fails int) {
		for i := 0; i < len(test); i++ {
			log.Printf("evaluating presition %d/%d %.1f %%\n", i, len(test), 100.0*float64(i)/float64(len(test)))
			result := knn.FitPoint(test[i].Point())
			if result != test[i].MRILabel {
				fails++
			}
		}
		log.Printf("evaluating presition %d/%d %.1f %%\n", len(test), len(test), 100.0)
		return
	}
	failsForTumors := evaluate(pev.TestTumors)
	failsForNoTumors := evaluate(pev.TestNoTumors)
	return &PresitionEvaluation{
		FailsForTumors:    failsForTumors,
		FailsForNoTumors:  failsForNoTumors,
		TumorsAll:         len(pev.TestTumors),
		NoTumorsAll:       len(pev.TestNoTumors),
		PresitionTumors:   100.0 * float64(len(pev.TestTumors)-failsForTumors) / float64(len(pev.TestTumors)),
		PresitionNoTumors: 100.0 * float64(len(pev.TestNoTumors)-failsForNoTumors) / float64(len(pev.TestNoTumors)),
		FailsAll:          failsForNoTumors + failsForTumors,
		PresitionAll:      100.0 * (float64(len(pev.TestTumors) + len(pev.TestNoTumors) - failsForTumors - failsForNoTumors)) / float64(len(pev.TestTumors)+len(pev.TestNoTumors)),
		All:               len(pev.TestTumors) + len(pev.TestNoTumors),
	}
}

type PointEvaluation struct {
	Tumors   []*lib.DataPoint
	NoTumors []*lib.DataPoint
}

func (pte *PointEvaluation) PartitionFiles(part float64) *PartitionEvaluation {
	log.Println("partitioning files...")
	tumorsPart := int(float64(len(pte.Tumors)) * part)
	notumorsPart := int(float64(len(pte.NoTumors)) * part)
	pev := &PartitionEvaluation{
		PointEvaluation: PointEvaluation{
			Tumors:   pte.Tumors[tumorsPart:],
			NoTumors: pte.NoTumors[notumorsPart:],
		},
		TestTumors:   pte.Tumors[:tumorsPart],
		TestNoTumors: pte.NoTumors[:tumorsPart],
	}
	return pev
}

func GeneratePoints(config cfg.Configuration, tumors, noTumors []string) *PointEvaluation {
	model := lib.NewFDModelNoLogLog(config)
	evaluate := func(output *[]*lib.DataPoint, label string, files []string) error {
		for i := 0; i < len(files); i++ {
			log.Printf("evaluating-%s %d/%d %.1f %%\n", label, i, len(files), 100.0*float64(i)/float64(len(files)))
			data, err := os.ReadFile(files[i])
			if err != nil {
				return err
			}
			imgs, err := ReadImage(data)
			if err != nil {
				return err
			}
			for _, img := range imgs {
				ev, err := model.Eval(img.(*image.RGBA))
				if err != nil {
					log.Fatalln(err)
				}
				*output = append(*output, &lib.DataPoint{
					FD:       ev.GetPoints(),
					MRILabel: label,
				})
			}
		}
		log.Printf("evaluating-%s %d/%d %.1f\n", label, len(files), len(files), 100.0)
		return nil
	}
	pte := new(PointEvaluation)
	if err := evaluate(&pte.Tumors, "tumors", tumors); err != nil {
		log.Fatalln(err)
	}
	if err := evaluate(&pte.NoTumors, "notumors", noTumors); err != nil {
		log.Fatalln(err)
	}
	return pte
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
