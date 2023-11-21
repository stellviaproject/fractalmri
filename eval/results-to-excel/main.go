package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	o := flag.String("o", "", "the file for saving the excel")
	d := flag.String("d", "", "the directory with json files")
	flag.Parse()
	ens, err := os.ReadDir(*d)
	if err != nil {
		log.Fatalln(err)
	}
	files := make([]string, 0, len(ens))
	for i := 0; i < len(ens); i++ {
		name := ens[i].Name()
		if strings.HasSuffix(name, ".json") {
			files = append(files, path.Join(*d, ens[i].Name()))
		}
	}
	presitions := make([]*PresitionEvaluation, 0, len(files))
	ratios := make([]int, 0, len(presitions))
	for i := 0; i < len(files); i++ {
		base := path.Base(files[i])
		ratioStr := strings.TrimRight(base, ".json")
		ratioStr = strings.TrimLeft(ratioStr, "window-ratio-")
		ratio, err := strconv.Atoi(ratioStr)
		if err != nil {
			log.Fatalln(err)
		}
		data, err := os.ReadFile(files[i])
		if err != nil {
			log.Fatalln(err)
		}
		presitionsInFile := []*PresitionEvaluation{}
		if err := json.Unmarshal(data, &presitionsInFile); err != nil {
			log.Fatalln(err)
		}
		presitions = append(presitions, presitionsInFile...)
		for j := 0; j < len(presitionsInFile); j++ {
			ratios = append(ratios, ratio)
		}
	}
	sort.Sort(&Sorter{
		presitions: &presitions,
		ratios:     &ratios,
	})
	SaveToExcel(*o, presitions, ratios)
}

type Sorter struct {
	presitions *[]*PresitionEvaluation
	ratios     *[]int
}

func (s *Sorter) Len() int {
	return len(*s.presitions)
}

func (s *Sorter) Less(i, j int) bool {
	return (*s.ratios)[i] < (*s.ratios)[j]
}

func (s *Sorter) Swap(i, j int) {
	(*s.presitions)[i], (*s.presitions)[j] = (*s.presitions)[j], (*s.presitions)[i]
	(*s.ratios)[i], (*s.ratios)[j] = (*s.ratios)[j], (*s.ratios)[i]
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

func SaveToCSV(fileName string, presitions []*PresitionEvaluation, ratios []int) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	header := []string{
		"FallidosTumor", "FallidosNoTumor", "Tumores", "NoTumores", "PresicionTumores", "PresicionNoTumores",
		"Fallidos", "Precision", "Total", "TamañoVentana",
	}
	table := [][]string{header}
	for i, p := range presitions {
		table = append(table, []string{
			fmt.Sprint(p.FailsForTumors),
			fmt.Sprint(p.FailsForNoTumors),
			fmt.Sprint(p.TumorsAll),
			fmt.Sprint(p.NoTumorsAll),
			fmt.Sprint(p.PresitionTumors),
			fmt.Sprint(p.PresitionNoTumors),
			fmt.Sprint(p.FailsAll),
			fmt.Sprint(p.PresitionAll),
			fmt.Sprint(p.All),
			fmt.Sprint(ratios[i]),
		})
	}
	if err := writer.WriteAll(table); err != nil {
		log.Fatalln(err)
	}
	writer.Flush()
}

func SaveToExcel(fileName string, presitions []*PresitionEvaluation, ratios []int) {
	f := excelize.NewFile()
	const SheetID = "Sheet1"
	sheet1 := f.NewSheet(SheetID)
	f.SetCellValue(SheetID, "A1", "FallidosTumor")
	f.SetCellValue(SheetID, "B1", "FallidosNoTumor")
	f.SetCellValue(SheetID, "C1", "Tumores")
	f.SetCellValue(SheetID, "D1", "NoTumores")
	f.SetCellValue(SheetID, "E1", "PresicionTumores")
	f.SetCellValue(SheetID, "F1", "PresicionNoTumores")
	f.SetCellValue(SheetID, "G1", "Fallidos")
	f.SetCellValue(SheetID, "H1", "Precision")
	f.SetCellValue(SheetID, "I1", "Total")
	f.SetCellValue(SheetID, "J1", "TamañoVentana")
	for i := 0; i < len(presitions); i++ {
		p := presitions[i]
		f.SetCellValue(SheetID, fmt.Sprintf("A%d", i+2), p.FailsForTumors)
		f.SetCellValue(SheetID, fmt.Sprintf("B%d", i+2), p.FailsForNoTumors)
		f.SetCellValue(SheetID, fmt.Sprintf("C%d", i+2), p.TumorsAll)
		f.SetCellValue(SheetID, fmt.Sprintf("D%d", i+2), p.NoTumorsAll)
		f.SetCellValue(SheetID, fmt.Sprintf("E%d", i+2), p.PresitionTumors)
		f.SetCellValue(SheetID, fmt.Sprintf("F%d", i+2), p.PresitionNoTumors)
		f.SetCellValue(SheetID, fmt.Sprintf("G%d", i+2), p.FailsAll)
		f.SetCellValue(SheetID, fmt.Sprintf("H%d", i+2), p.PresitionAll)
		f.SetCellValue(SheetID, fmt.Sprintf("I%d", i+2), p.All)
		f.SetCellValue(SheetID, fmt.Sprintf("J%d", i+2), ratios[i])
	}
	f.SetActiveSheet(sheet1)
	if err := f.SaveAs(fileName); err != nil {
		log.Fatalln(err)
	}
}
