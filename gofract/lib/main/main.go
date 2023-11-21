package main

import (
	"encoding/json"
	"fmt"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/lib"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
)

func main() {
	MainTestSample()
}

func MainTestSample() {
	c := cfg.GetCFG()

	s, err := lib.LoadSample("./sample-p=0.983732.json")
	onErr(err)

	tumors, notumors, testTumors, testNoTumors := s.Tumors, s.NoTumors, s.TumorsTest, s.NoTumorsTest

	var k *lib.KNNFractal

	points, err := lib.LoadPoints("./points.json")
	if err != nil {
		k = lib.NewKNNFractal(c, []*lib.DataPoint{})
		onErr(k.TrainWithFiles(GetFilesSet(tumors, notumors)))
		lib.SavePoints("./points.json", k.GetPoints())
	} else {
		k = lib.NewKNNFractal(c, points)
	}
	failsTumors := 0
	for i := range testTumors {
		tumor, err := lib.ReadImage(testTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(tumor); err != nil {
			panic(err)
		} else if label != "tumor" {
			failsTumors++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testTumors), int(float64(i+1)*100/float64(len(testTumors))))
	}
	log.Printf("tumor(%d,%d)", failsTumors, len(testTumors))
	failsNoTumors := 0
	for i := range testNoTumors {
		notumor, err := lib.ReadImage(testNoTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(notumor); err != nil {
			panic(err)
		} else if label != "notumor" {
			failsNoTumors++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testNoTumors), int(float64(i+1)*100/float64(len(testNoTumors))))
	}
	// (resultado con tumores + resultado sin tumores) / (tumores + sin tumores)
	presition := float64(failsTumors+failsNoTumors) / float64(len(testTumors)+len(testNoTumors))
	log.Printf("pressition %f tumor(%d,%d) notumor(%d,%d)", presition, failsTumors, len(testTumors), failsNoTumors, len(testNoTumors))
	//2023/11/10 15:46:32
	//2023/11/10 16:22:01 pressition 0.982775 tumor(151,152) notumor(876,893)
	//2023/08/04 22:19:05 progress  9 / 10   90 %  p= 0.9751196172248804  cp= 0.9837320574162679
}

func MainGenetic() {
	tumors, _ := GetFiles("./tumors")
	notumors, _ := GetFiles("./notumors")
	s := &lib.Sample{}
	s.Tumors, s.NoTumors, s.TumorsTest, s.NoTumorsTest = PartitionFiles(0.8, tumors, notumors)
	s = Reduce(s, 10)
	ns, p := s.Optimize(10, cfg.GetCFG())
	lib.SaveSample(fmt.Sprintf("./sample-p=%f.json", p), ns)
	cr := FindBestCfg(s, 50, 10, 100, 0.2)
	cr.Save("./cromosome.json")
}

func MainBestSample() {
	err := SetOutputToFileAndTerminal("./output.log")
	onErr(err)
	tumors, _ := GetFiles("../tumors")
	notumors, _ := GetFiles("../notumors")
	s := lib.NewSample(0.8, tumors, notumors)
	ns, p := s.Optimize(10, cfg.GetCFG())
	lib.SaveSample(fmt.Sprintf("sample-p=%f.json", p), ns)
}

func SetOutputToFileAndTerminal(filename string) error {
	// Open the file for writing
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Set the log output to the file and the terminal
	log.SetOutput(io.MultiWriter(file, os.Stdout))

	return nil
}

func Main() {
	//2023/08/03 15:06:03 presition 0.553110048 tumor(91,152) notumor(487,893)
	//c.KNN.K = 100
	//c.WindowRatio = 10
	//c.Umbral = cfg.Range(0.1, 1.4, 0.07)
	//presition 0.632535885 tumor(140,152) notumor(521,893)
	//2023/08/04 03:20:23 presition 0.271837 tumor(140,152) notumor(478,893)
	c := cfg.GetCFG()
	c.KNN.K = 10
	c.WindowRatio = 40
	c.Umbral = cfg.Range(0.1, 1.4, 0.05)

	tumors, _ := GetFiles("./tumors")
	notumors, _ := GetFiles("./notumors")
	Dis(tumors)
	Dis(notumors)
	var testTumors, testNoTumors []string
	tumors, notumors, testTumors, testNoTumors = PartitionFiles(0.2, tumors, notumors)

	var k *lib.KNNFractal

	points, err := lib.LoadPoints("./points.json")
	if err != nil {
		k = lib.NewKNNFractal(c, []*lib.DataPoint{})
		//Error on [59:] Umbral=Range(0.1, 1.5, 0.05)
		onErr(k.TrainWithFiles(GetFilesSet(tumors, notumors)))
		lib.SavePoints("./points.json", k.GetPoints())
	} else {
		k = lib.NewKNNFractal(c, points)
	}
	tmOk := 0
	for i := range testTumors {
		tumor, err := lib.ReadImage(testTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(tumor); err != nil {
			panic(err)
		} else if label == "tumor" {
			tmOk++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testTumors), int(float64(i+1)*100/float64(len(testTumors))))
	}
	log.Printf("tumor(%d,%d)", tmOk, len(testTumors))
	noTmOk := 0
	for i := range testNoTumors {
		notumor, err := lib.ReadImage(testNoTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(notumor); err != nil {
			panic(err)
		} else if label == "notumor" {
			noTmOk++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testNoTumors), int(float64(i+1)*100/float64(len(testNoTumors))))
	}
	tmOkP := float64(tmOk) / float64(len(testTumors))
	noTmOkP := float64(noTmOk) / float64(len(testNoTumors))
	presition := ((1 - tmOkP) + (1 - noTmOkP)) / 2.0
	log.Printf("pressition %f tumor(%d,%d) notumor(%d,%d)", presition, tmOk, len(testTumors), noTmOk, len(testNoTumors))
	//0.08904109589041098
}

func onErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Cromosome struct {
	Step float64
	WR   int
	P    float64
}

func (cr Cromosome) Save(fileName string) {
	data, err := json.Marshal(&cr)
	onErr(err)
	err = os.WriteFile(fileName, data, os.ModePerm|os.ModeDevice)
	onErr(err)
}

func (crm Cromosome) GetCFG() cfg.Configuration {
	conf := cfg.GetCFG()
	conf.Umbral = cfg.Range(0.1, 1.4, crm.Step)
	conf.WindowRatio = crm.WR
	return conf
}

func GenPopulation(size int, s *lib.Sample) (pop []Cromosome) {
	for i := 0; i < size; i++ {
		cromo := Cromosome{
			Step: 0.5 + rand.Float64()*2.0,
			WR:   5 + rand.Intn(45),
		}
		pop = append(pop, cromo)
		s.EvalSample(cromo.GetCFG())
		cromo.P, _, _ = s.EvalKNN(cromo.GetCFG())
	}
	return
}

func Selection(list []Cromosome) (result []Cromosome) {
	result = make([]Cromosome, len(list)/2)
	for i := 0; i < len(result); i++ {
		n := rand.Intn(len(list))
		result = append(result, list[n])
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].P > list[j].P
	})
	maxs := len(result) / 10
	if maxs == 0 {
		maxs++
	}
	for i := 0; i < maxs; i++ {
		result[rand.Intn(len(result))] = list[i]
	}
	return result
}

func Cross(pop []Cromosome, s *lib.Sample) (crs []Cromosome) {
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].P > pop[j].P
	})
	for i := 0; i < len(pop)/2; i++ {
		p1 := rand.Intn(len(pop))
		p2 := rand.Intn(len(pop))
		cmax, cmin := pop[p1], pop[p2]
		if cmax.P < cmin.P {
			cmax, cmin = cmin, cmax
		}
		r2 := (cmax.P - cmin.P) / 2
		wrR := cmax.WR - cmin.WR
		stepR := cmax.Step - cmin.Step

		sig := 1 - rand.Intn(2)*2
		wr := (cmax.WR+cmin.WR)/2 + sig*int(float64(wrR)*r2)
		sig = 1 - rand.Intn(2)*2
		step := (cmax.Step+cmin.Step)/2 + float64(sig)*stepR*r2
		if step < 0.05 {
			step = 0.5
		}
		if step > 0.2 {
			step = 0.2
		}
		if wr < 5 {
			wr = 5
		}
		if wr > 50 {
			wr = 50
		}
		cromo := Cromosome{
			Step: step,
			WR:   wr,
		}
		s.EvalKNN(cromo.GetCFG())
		cromo.P, _, _ = s.EvalKNN(cromo.GetCFG())
		crs = append(crs, cromo)
	}
	return
}

func Mutate(mut float64, cromo *Cromosome, s *lib.Sample) {
	if rand.Float64() < mut {
		newCromo := Cromosome{
			Step: 0.5 + rand.Float64()*2.0,
			WR:   5 + rand.Intn(45),
		}
		cromo.Step = (newCromo.Step + cromo.Step) / 2
		cromo.WR = (newCromo.WR + cromo.WR) / 2
		s.EvalKNN(cromo.GetCFG())
		cromo.P, _, _ = s.EvalKNN(cromo.GetCFG())
	}
}

func GetBestCromo(pop []Cromosome) (best Cromosome) {
	best = pop[0]
	for i := 1; i < len(pop); i++ {
		if best.P < pop[i].P {
			best = pop[i]
		}
	}
	return
}

func FindBestCfg(s *lib.Sample, p, g, nopg int, mut float64) Cromosome {
	pop := GenPopulation(p, s)
	i := 0
	best := GetBestCromo(pop)
	nop := 0
	for i < g && nop < nopg {
		log.Println("selection...")
		pop = Selection(pop)
		log.Println("cross...")
		pop = Cross(pop, s)
		log.Println("mutate...")
		for j := 0; j < len(pop); j++ {
			Mutate(mut, &pop[j], s)
		}
		log.Println("get best...")
		bestPop := GetBestCromo(pop)
		if bestPop.P > best.P {
			best = bestPop
			nop = 0
		} else {
			nop++
		}
		log.Println("population ", len(pop), "/", p, " generation ", i, "/", g, " wr: ", best.WR, " step: ", best.Step, " p: ", best.P)
		i++
	}
	return best
}

func Reduce(s *lib.Sample, size int) (o *lib.Sample) {
	o.Tumors = ReduceOne(size, s.Tumors)
	o.NoTumors = ReduceOne(size, s.NoTumors)
	o.TumorsTest = ReduceOne(size, s.TumorsTest)
	o.NoTumorsTest = ReduceOne(size, s.NoTumorsTest)
	return
}

func ReduceOne(size int, in []string) (out []string) {
	out = make([]string, size)
	for i := 0; i < size; i++ {
		n := rand.Intn(len(in))
		out[i] = in[n]
		in = append(in[:n], in[n+1:]...)
	}
	return
}

func Optimize() {
	c := cfg.GetCFG()

	tumors, _ := GetFiles("./tumors")
	notumors, _ := GetFiles("./notumors")
	Dis(tumors)
	Dis(notumors)
	var testTumors, testNoTumors []string
	tumors, notumors, testTumors, testNoTumors = PartitionFiles(0.2, tumors, notumors)

	var k *lib.KNNFractal

	points, err := lib.LoadPoints("./points.json")
	if err != nil {
		k = lib.NewKNNFractal(c, []*lib.DataPoint{})
		//Error on [59:] Umbral=Range(0.1, 1.5, 0.05)
		onErr(k.TrainWithFiles(GetFilesSet(tumors, notumors)))
		lib.SavePoints("./points.json", k.GetPoints())
	} else {
		k = lib.NewKNNFractal(c, points)
	}
	tmOk := 0
	for i := range testTumors {
		tumor, err := lib.ReadImage(testTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(tumor); err != nil {
			panic(err)
		} else if label == "tumor" {
			tmOk++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testTumors), int(float64(i+1)*100/float64(len(testTumors))))
	}
	log.Printf("tumor(%d,%d)", tmOk, len(testTumors))
	noTmOk := 0
	for i := range testNoTumors {
		notumor, err := lib.ReadImage(testNoTumors[i])
		if err != nil {
			log.Fatalln(err)
		}
		if label, err := k.Fit(notumor); err != nil {
			panic(err)
		} else if label == "notumor" {
			noTmOk++
		}
		log.Printf("eval tumors %d/%d %d%%\n", i+1, len(testNoTumors), int(float64(i+1)*100/float64(len(testNoTumors))))
	}
	tmOkP := float64(tmOk) / float64(len(testTumors))
	noTmOkP := float64(noTmOk) / float64(len(testNoTumors))
	presition := ((1 - tmOkP) + (1 - noTmOkP)) / 2.0
	log.Printf("pressition %f tumor(%d,%d) notumor(%d,%d)", presition, tmOk, len(testTumors), noTmOk, len(testNoTumors))
	//0.08904109589041098
}

func GetFilesSet(tumors, notumors []string) []*lib.FileSetItem {
	set := make([]*lib.FileSetItem, 0, len(tumors)+len(notumors))
	for i := range tumors {
		set = append(set, lib.NewFileSetItem(tumors[i], "tumor"))
	}
	for i := range notumors {
		set = append(set, lib.NewFileSetItem(notumors[i], "notumor"))
	}
	return set
}

func Dis(set []string) {
	for i := range set {
		j := rand.Intn(len(set))
		set[i], set[j] = set[j], set[i]
	}
}

func GetFiles(folder string) (imageFiles, maskFiles []string) {
	files, err := os.ReadDir(folder)
	onErr(err)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "mask.png") {
			maskFiles = append(maskFiles, path.Join(folder, file.Name()))
		} else {
			imageFiles = append(imageFiles, path.Join(folder, file.Name()))
		}
	}
	return
}

func PartitionFiles(part float64, tumors, notumors []string) (tms, noTms, testTms, testNoTms []string) {
	notestPart := int(float64(len(tumors)) * part)
	tms = tumors[:notestPart]
	noTms = notumors[:notestPart]
	testTms = tumors[notestPart:]
	testNoTms = notumors[notestPart:]
	return
}
