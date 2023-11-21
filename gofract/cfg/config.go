package cfg

import (
	"encoding/json"
	"math"
	"os"
)

/*
Representa la configuracion con que se ejecuta el algoritmo
*/
type Configuration struct {
	KNN                 KNN       //La configuracion para ejecutar KNN
	Buffer              int       //El buffer para las tuberias
	Parallel            int       //El nivel de paralelismo del algoritmo
	WindowRatio         int       //El tamaño de la ventana para el filtro multifractal
	DenoiserSigmaColor  float64   //El valor sigma-color para reducir el ruido
	DenoiserSigmaSpace  float64   //El valor sigma-space para reducir el ruido
	DenoiserDiameter    int       //El diametro del filtro para reducir el ruido
	DenoiserUmbralColor float64   //Un umbral para eliminar el ruido
	MinUmbral           float64   //Umbral minimo para segmentar
	MaxUmbral           float64   //Umbral maximo para segmentar
	MinArea             int       //Area minima para el filtro de separacion de objetos
	MaxArea             int       //Area maxima para el filtro de separacion de objetos
	Ratio               int       //Radio para la busqueda del pixel cercano que forma un objeto
	BoxSizes            []int     //El tamaño de las cajas para calcular la dimension fractal
	LogLogWidth         int       //El ancho de la imagen de la grafica loglog
	LogLogHeight        int       //El alto de la imagen de la grafica loglog
	Umbral              []float64 //Los umbrales para segmentar el espectro multifractal
}

func (cfg Configuration) GetParallel() int {
	prll := cfg.Parallel
	if prll > cfg.Buffer {
		prll = cfg.Buffer
	}
	return prll
}

// Guarda la configuracion
func SaveCFG(fileName string, cfg Configuration) error {
	//codifica la configuracion en un json
	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	//escribe el json en un archivo
	return os.WriteFile(fileName, data, os.ModeDevice|os.ModePerm)
}

// Carga la configuracion desde un archivo
func LoadCFG(fileName string) (Configuration, error) {
	cfg := Configuration{}
	//leer el archivo
	data, err := os.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}
	//decodificar el json en la estructura
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

type Distance string

const (
	EuclideanDist          Distance = "euclidean"
	ManhattanDist          Distance = "manhattan"
	MinkowskiDist          Distance = "minkowski"
	ChebyshevDist          Distance = "chebyshev"
	PearsonCorrelationDist Distance = "pearson-correlation"
	HammingDist            Distance = "hamming"
)

type Selector string

const (
	BinarySelector        Selector = "binary"
	MultiClassSelector    Selector = "multi-class"
	SmoothInverseSelector Selector = "smooth-inverse"
)

type KNN struct {
	K              int      `json:"k"`               //Los K elementos para la clasificacion
	Distance       Distance `json:"distance"`        //El metodo para calcular la distancia entre los puntos del KNN
	Selector       Selector `json:"selector"`        //El metodo para la seleccion de la clase en el KNN
	MinkowskiRatio float64  `json:"minkowski-ratio"` //El radio de minkowski para el calculo de la distancia de minkowski
	WeightParam    float64  `json:"weight-param"`
	SmoothingParam float64  `json:"smoothing-param"`
}

// La configuracion por defecto
var cfg Configuration = Configuration{
	KNN: KNN{
		K:        10,
		Distance: EuclideanDist,
		Selector: BinarySelector,
	},
	BoxSizes:            []int{128, 64, 32, 16, 8, 4, 2},
	Parallel:            10,
	Buffer:              10,
	WindowRatio:         5,
	DenoiserSigmaColor:  0.100710,
	DenoiserSigmaSpace:  42.005347,
	DenoiserDiameter:    13,
	DenoiserUmbralColor: 10.560879,
	MinUmbral:           0.0,
	MaxUmbral:           2.0,
	Umbral:              Range(0.1, 1.4, 0.1),
	MinArea:             20,
	MaxArea:             math.MaxInt,
	Ratio:               5,
	LogLogWidth:         800,
	LogLogHeight:        600,
}

// Obtiene la configuracion
func GetCFG() Configuration {
	return cfg
}

// Crea un slice con los elementos de recorrer el intervalo [s,e] desplazandose en p con cada iteracion
func Range(s, e, p float64) []float64 {
	r := []float64{}
	for s <= e {
		r = append(r, s)
		s += p
	}
	return r
}
