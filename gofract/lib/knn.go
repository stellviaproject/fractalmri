package lib

import (
	"fractalmri/gofract/cfg"
	"image"
	"log"
	"os"
	"sync"

	"github.com/stellviaproject/go-ia/knn"
)

type KNNFractal struct {
	knnFD     *knn.KNN
	fdclassMd *FDModel
	cfg       cfg.Configuration
}

//export NewKNNFractal
func NewKNNFractal(cfg cfg.Configuration, points []*DataPoint) *KNNFractal {
	//Initialize points evaluation model
	//Initialize knn
	knnPoints := make([]knn.DataPoint, len(points))
	for i := range points {
		knnPoints[i] = points[i]
	}
	knnFD := knn.NewKNN(
		cfg.KNN.K,
		knn.NewEuclideanDist(),
		knn.NewMultiClassSelector(),
		knnPoints,
	)

	return &KNNFractal{
		knnFD:     knnFD,
		fdclassMd: NewFDModel(cfg),
		cfg:       cfg,
	}
}

func (knf *KNNFractal) TrainWithImages(dataset []*ImageSetItem) error {
	var err error
	prll := knf.cfg.GetParallel()
	ch := make(chan int, prll)
	wg := make(chan int, len(dataset))
	mtx := sync.Mutex{}
	pg := 0
	for i := range dataset {
		ch <- 0
		go func(i int) {
			defer func() {
				<-ch
				wg <- 0
			}()
			eval, e := knf.fdclassMd.Eval(dataset[i].Image)
			if e != nil {
				err = e
				return
			}
			mtx.Lock()
			knf.knnFD.Append(NewDataPoint(eval.GetPoints(), dataset[i].MRILabel))
			pg++
			log.Printf("progress: %d / %d", pg, len(dataset))
			mtx.Unlock()
		}(i)
	}
	if err != nil {
		return err
	}
	for i := 0; i < len(dataset); i++ {
		<-wg
	}
	return nil
}

//export TrainWithFiles
func (knf *KNNFractal) TrainWithFiles(dataset []*FileSetItem) error {
	var err error
	prll := knf.cfg.GetParallel()
	ch := make(chan int, prll)
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	pg := 0
	for i := range dataset {
		wg.Add(1)
		go func(i int) {
			ch <- 0
			defer func() { <-ch; wg.Done() }() // continue with others
			file, e := os.Open(dataset[i].FileName)
			if e != nil {
				err = e
				return
			}
			img, _, e := image.Decode(file)
			if e != nil {
				err = e
				return
			}
			file.Close()
			eval, e := knf.fdclassMd.Eval(img.(*image.RGBA))
			if e != nil {
				err = e
				return
			}
			mtx.Lock()
			knf.knnFD.Append(NewDataPoint(eval.GetPoints(), dataset[i].MRILabel))
			pg++
			log.Printf("progress: %d / %d", pg, len(dataset))
			mtx.Unlock()
		}(i)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

//export Fit
func (knf *KNNFractal) Fit(img *image.RGBA) (string, error) {
	ev, err := knf.fdclassMd.Eval(img)
	if err != nil {
		return "", err
	}
	return knf.knnFD.Fit(ev.GetPoints()).(string), nil
}

func (knf *KNNFractal) FitPoint(point knn.Point) string {
	return knf.knnFD.Fit(point).(string)
}

//export GetPoints
func (knf *KNNFractal) GetPoints() []*DataPoint {
	points := knf.knnFD.GetDataPoints()
	pts := make([]*DataPoint, len(points))
	for i := 0; i < len(points); i++ {
		pts[i] = points[i].(*DataPoint)
	}
	return pts
}
