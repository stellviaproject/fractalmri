package lib

import (
	"log"
	"testing"

	"github.com/stellviaproject/go-ia/knn"
)

func TestLoadPoints(t *testing.T) {
	points, err := LoadPoints("./test-points.json")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(points)
}

func TestSavePoints(t *testing.T) {
	points := []*DataPoint{
		NewDataPoint(knn.Point{0.0, 0.0, 0.0}, "dog"),
		NewDataPoint(knn.Point{0.0, 0.0, 1.0}, "cat"),
		NewDataPoint(knn.Point{0.0, 1.0, 0.0}, "bird"),
		NewDataPoint(knn.Point{1.0, 0.0, 0.0}, "horse"),
		NewDataPoint(knn.Point{0.0, 1.0, 1.0}, "fish"),
		NewDataPoint(knn.Point{1.0, 1.0, 0.0}, "lion"),
		NewDataPoint(knn.Point{1.0, 0.0, 1.0}, "cow"),
		NewDataPoint(knn.Point{1.0, 1.0, 1.0}, "human"),
	}
	err := SavePoints("./test-points.json", points)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
