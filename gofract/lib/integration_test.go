package lib

import (
	"fractalmri/gofract/cfg"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestKNN(t *testing.T) {
	tumors := LoadImages("./testing/tumors")
	notumors := LoadImages("./testing/notumors")
	var testTumors, testNoTumors []*image.RGBA
	tumors, notumors, testTumors, testNoTumors = Partition(0.8, tumors, notumors)
	knn := NewKNNFractal(cfg.GetCFG(), []*DataPoint{})
	//Error on [59:] Umbral=Range(0.1, 1.5, 0.05)
	onErr(knn.TrainWithImages(GetImageSet(tumors, notumors)))
	tmOk := 0
	for i := range testTumors {
		if label, err := knn.Fit(tumors[i]); err != nil {
			panic(err)
		} else if label == "tumor" {
			tmOk++
		}
	}
	noTmOk := 0
	for i := range testNoTumors {
		if label, err := knn.Fit(testNoTumors[i]); err != nil {
			panic(err)
		} else if label == "notumor" {
			noTmOk++
		}
	}
	tmOkP := float64(tmOk) / float64(len(testTumors))
	noTmOkP := float64(noTmOk) / float64(len(testNoTumors))
	presition := ((1 - tmOkP) + (1 - noTmOkP)) / 2.0
	log.Println(presition)
	//0.08904109589041098
}

func Partition(part float64, tumors, notumors []*image.RGBA) (tms, noTms, testTms, testNoTms []*image.RGBA) {
	notestPart := int(float64(len(tumors)) * part)
	tms = tumors[:notestPart]
	noTms = notumors[:notestPart]
	testTms = tumors[notestPart:]
	testNoTms = notumors[notestPart:]
	return
}

func GetImageSet(tumors, notumors []*image.RGBA) []*ImageSetItem {
	set := make([]*ImageSetItem, 0, len(tumors)+len(notumors))
	for i := range tumors {
		set = append(set, NewImageSetItem(tumors[i], "tumor"))
	}
	for i := range notumors {
		set = append(set, NewImageSetItem(notumors[i], "notumor"))
	}
	return set
}

func LoadImages(folder string) []*image.RGBA {
	files, err := os.ReadDir(folder)
	onErr(err)
	images := make([]*image.RGBA, 0, len(files)/2)
	p := make(chan int, 20)
	wg := sync.WaitGroup{}
	for i := range files {
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-p
				wg.Done()
			}()
			p <- 0
			if files[i].IsDir() {
				return
			}
			fileName := files[i].Name()
			if strings.HasSuffix(fileName, "mask.png") {
				return
			}
			fileName = filepath.Join(folder, fileName)
			file, err := os.Open(fileName)
			onErr(err)
			img, err := png.Decode(file)
			onErr(err)
			images = append(images, img.(*image.RGBA))
		}(i)
	}
	wg.Wait()
	return images
}

func onErr(err error) {
	log.Fatalln(err)
}
