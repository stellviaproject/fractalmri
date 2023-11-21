package lib

import (
	"encoding/json"
	"fractalmri/gofract/cfg"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/stellviaproject/go-ia/knn"
)

const (
	TumorLabel   = "tumor"
	NoTumorLabel = "notumor"
)

func EvalPoints(c cfg.Configuration, files []string) ([]knn.Point, error) {
	model := NewFDModel(c)
	//point list
	list := []knn.Point{}
	//initialize control variables
	var err error
	//parallelism level
	prll := c.GetParallel()
	ch := make(chan int, prll)
	//wait for all finish
	wg := sync.WaitGroup{}
	//lock when use list
	mtx := sync.Mutex{}
	//notify progress
	pg := 0
	for i := range files {
		wg.Add(1)
		go func(i int) {
			ch <- 0
			defer func() { <-ch; wg.Done() }() // continue with others
			img, e := ReadImage(files[i])
			if e != nil {
				err = e
				return
			}
			eval, e := model.Eval(img)
			if e != nil {
				err = e
				return
			}
			mtx.Lock()
			list = append(list, eval.GetPoints())
			pg++
			log.Printf("progress: %d / %d", pg, len(files))
			mtx.Unlock()
		}(i)
	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return list, nil
}

type Sample struct {
	Presition          float64
	Tumors             []string
	NoTumors           []string
	TumorsTest         []string
	NoTumorsTest       []string
	TumorsPoints       []knn.Point
	NoTumorsPoints     []knn.Point
	TumorsTestPoints   []knn.Point
	NoTumorsTestPoints []knn.Point
}

func NewSample(part float64, tumors, notumors []string) *Sample {
	s := &Sample{}
	notestPart := int(float64(len(tumors)) * part)
	s.Tumors = tumors[:notestPart]
	s.NoTumors = notumors[:notestPart]
	s.TumorsTest = tumors[notestPart:]
	s.NoTumorsTest = notumors[notestPart:]
	return s
}

func (s *Sample) GetPoints() []*DataPoint {
	points := make([]*DataPoint, 0, len(s.TumorsPoints)+len(s.NoTumorsPoints))
	for i := 0; i < len(s.TumorsPoints); i++ {
		points = append(points, NewDataPoint(s.TumorsPoints[i], TumorLabel))
	}
	for i := 0; i < len(s.NoTumorsPoints); i++ {
		points = append(points, NewDataPoint(s.NoTumorsPoints[i], NoTumorLabel))
	}
	return points
}

func (s *Sample) EvalSample(c cfg.Configuration) error {
	log.Println("evaluate tumors points")
	var err error
	s.TumorsPoints, err = EvalPoints(c, s.Tumors)
	if err != nil {
		return err
	}
	log.Println("evaluate notumors points")
	s.NoTumorsPoints, err = EvalPoints(c, s.NoTumors)
	if err != nil {
		return err
	}
	log.Println("evaluate test tumors points")
	s.TumorsTestPoints, err = EvalPoints(c, s.TumorsTest)
	if err != nil {
		return err
	}
	log.Println("evaluate test no tumors points")
	s.NoTumorsTestPoints, err = EvalPoints(c, s.NoTumorsTest)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sample) EvalKNN(c cfg.Configuration) (p float64, fakesTumors, fakesNoTumors []int) {
	log.Println("eval knn / train...")
	//Create KNN
	k := NewKNNFractal(c, s.GetPoints())
	//Initialize control variables
	wg := sync.WaitGroup{}
	th := make(chan int, 4)
	mt := sync.Mutex{}
	count := 0
	//Initialize save variables
	fakesTumors = make([]int, 0, 50)
	fakesNoTumors = make([]int, 0, 50)
	pg := 0
	//Start process
	for i := 0; i < len(s.TumorsTestPoints); i++ {
		//Control parallel
		th <- 0
		wg.Add(1)
		//Start gorutine
		go func(i int) {
			defer wg.Done()
			lbl := k.FitPoint(s.TumorsTestPoints[i])
			mt.Lock()
			if lbl == TumorLabel {
				count++
			} else {
				fakesTumors = append(fakesTumors, i)
			}
			//Notify progress
			pg++
			log.Println("eval tumors ", pg, "/", len(s.TumorsTest), " ", 100*pg/len(s.TumorsTest), "%")
			//Unlock gorutine
			mt.Unlock()
			<-th
		}(i)
	}
	pg = 0
	for i := 0; i < len(s.NoTumorsTestPoints); i++ {
		th <- 0
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lbl := k.FitPoint(s.NoTumorsTestPoints[i])
			mt.Lock()
			if lbl == NoTumorLabel {
				count++
			} else {
				fakesNoTumors = append(fakesNoTumors, i)
			}
			//Notify progress
			pg++
			log.Println("eval notumors ", pg, "/", len(s.NoTumorsTest), " ", 100*pg/len(s.NoTumorsTest), "%")
			//Unlock gorutine
			mt.Unlock()
			<-th
		}(i)
	}
	wg.Wait()
	p = float64(count) / float64(len(s.TumorsTest)+len(s.NoTumorsTest))
	return
}

func (s *Sample) Optimize(n int, c cfg.Configuration) (*Sample, error) {
	if err := s.EvalSample(c); err != nil {
		return nil, err
	}
	log.Println("optimize begin...")
	p, fkT, fkN := s.EvalKNN(c)
	curS, curP := s, p
	for i := 0; i < n; i++ {
		nextS := curS.ApplyFakes(fkT, fkN)
		p, fkT, fkN = nextS.EvalKNN(c)
		if p > curP {
			curP = p
			curS = nextS
		}
		log.Println("progress ", i, "/", n, " ", i*100/n, "%", " p=", p, " cp=", curP)
	}
	curS.Presition = p
	return curS, nil
}

func (s *Sample) ApplyFakes(fakeTumors, fakeNoTumors []int) (o *Sample) {
	//create output variable
	o = &Sample{}
	//initialize fields of output variable
	//init file names
	o.Tumors = make([]string, len(s.Tumors))
	o.TumorsTest = make([]string, len(s.TumorsTest))
	o.NoTumors = make([]string, len(s.NoTumors))
	o.NoTumorsTest = make([]string, len(s.NoTumorsTest))
	//init points
	o.TumorsPoints = make([]knn.Point, len(s.TumorsPoints))
	o.TumorsTestPoints = make([]knn.Point, len(s.TumorsTestPoints))
	o.NoTumorsPoints = make([]knn.Point, len(s.NoTumorsPoints))
	o.NoTumorsTestPoints = make([]knn.Point, len(s.NoTumorsTestPoints))
	//copy current to output
	//copy file names
	copy(o.Tumors, s.Tumors)
	copy(o.TumorsTest, s.TumorsTest)
	copy(o.NoTumors, s.NoTumors)
	copy(o.NoTumorsTest, s.NoTumorsTest)
	//copy points
	copy(o.TumorsPoints, s.TumorsPoints)
	copy(o.TumorsTestPoints, s.TumorsTestPoints)
	copy(o.NoTumorsPoints, s.NoTumorsPoints)
	copy(o.NoTumorsTestPoints, s.NoTumorsTestPoints)
	//interchange test fakes with random no test
	//interchange tumors
	for i := 0; i < len(fakeTumors); i++ {
		index := rand.Intn(len(o.Tumors))
		fkIndex := fakeTumors[i]
		o.Tumors[index], o.TumorsTest[fkIndex] = o.TumorsTest[fkIndex], o.Tumors[index]
		o.TumorsPoints[index], o.TumorsTestPoints[fkIndex] = o.TumorsTestPoints[fkIndex], o.TumorsPoints[index]
	}
	//interchange no tumors
	for i := 0; i < len(fakeNoTumors); i++ {
		index := rand.Intn(len(o.NoTumors))
		fkIndex := fakeNoTumors[i]
		o.NoTumors[index], o.NoTumorsTest[fkIndex] = o.NoTumorsTest[fkIndex], o.NoTumors[index]
		o.NoTumorsPoints[index], o.NoTumorsTestPoints[fkIndex] = o.NoTumorsTestPoints[fkIndex], o.NoTumorsPoints[index]
	}
	return
}

func SaveSample(fileName string, s *Sample) error {
	data, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, os.ModePerm|os.ModeDevice)
	return err
}

func LoadSample(fileName string) (*Sample, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	s := new(Sample)
	err = json.Unmarshal(data, s)
	return s, err
}
