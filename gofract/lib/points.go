package lib

import (
	"encoding/json"
	"os"

	"github.com/stellviaproject/go-ia/knn"
)

type DataPoint struct {
	FD       knn.Point `json:"fd" pipe:"fd"`
	MRILabel string    `json:"label" pipe:"-"`
}

//export NewDataPoint
func NewDataPoint(fdPoint knn.Point, mriLabel string) *DataPoint {
	return &DataPoint{
		FD:       fdPoint,
		MRILabel: mriLabel,
	}
}

//export LoadPoints
func LoadPoints(fileName string) ([]*DataPoint, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var points []*DataPoint
	if err := json.Unmarshal(data, &points); err != nil {
		return nil, err
	}
	return points, err
}

//export SavePoints
func SavePoints(fileName string, points []*DataPoint) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Truncate(0); err != nil {
		return err
	}
	buffer, err := json.Marshal(points)
	if err != nil {
		return err
	}
	_, err = file.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

//export Point
func (dp *DataPoint) Point() knn.Point {
	return dp.FD
}

//export Label
func (dp *DataPoint) Label() any {
	return dp.MRILabel
}
