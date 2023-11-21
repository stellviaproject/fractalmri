package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"math"
	"reflect"
	"sync"

	arch "github.com/stellviaproject/pipfil-arch"
)

// El filtro que calcula la lacunaridad
func NewLacFilter(cfg cfg.Configuration, objInPipe, lacOutPipe arch.Pipe) arch.Filter {
	var fn any
	if objInPipe.CheckType() == reflect.TypeOf(pipes.ImageBin{}) {
		fn = lacunarityFn(cfg)
	} else if objInPipe.CheckType() == reflect.TypeOf(&pipes.ImageObject{}) {
		fn = func(obj *pipes.ImageObject) []float64 {
			return lacunarityFn(cfg)(obj.Image)
		}
	} else {
		panic("unknown type for lacunarity input")
	}
	return arch.NewFilterWithPipes(
		"lacunarity",
		fn,
		arch.WithPipes(objInPipe),
		arch.WithPipes(lacOutPipe),
		arch.WithLens(),
	)
}

func lacunarityFn(cfg cfg.Configuration) func(bin pipes.ImageBin) []float64 {
	return func(bin pipes.ImageBin) []float64 {
		// Get configuration parameters
		boxSizes := cfg.BoxSizes
		parallel := cfg.Parallel
		// Prepare parallelism control
		ch := make(chan int, parallel)
		wg := sync.WaitGroup{}
		// Prepare lacunarities array
		lacs := make([]float64, len(boxSizes))
		// Compute lacunarity for each box-size
		for index, window := range boxSizes {
			// avoid another gorutine more than allowed
			ch <- 0
			// add a new gorutine to wait
			wg.Add(1)
			go func(index, window int) {
				defer func() {
					// continue with other gorutines
					<-ch
					// mark this gorutine as done
					wg.Done()
				}()
				// get image parameters
				width := bin.Width()
				height := bin.Height()
				// parameters to measure in image
				numBoxes := 0
				sum := 0.0
				for y := 0; y < height; y += window {
					for x := 0; x < width; x += window {
						numBoxes++
						boxCount := 0
						for j := y; j < y+window; j++ {
							for i := x; i < x+window; i++ {
								if i < width && j < height {
									if bin.At(i, j) {
										boxCount++
									}
								}
							}
						}
						sum += math.Pow(float64(boxCount), 2)
					}
				}
				mean := sum / float64(numBoxes)
				sum = 0
				for y := 0; y < height; y += window {
					for x := 0; x < width; x += window {
						boxCount := 0
						for j := y; j < y+window; j++ {
							for i := x; i < x+window; i++ {
								if i < width && j < height {
									if bin.At(i, j) {
										boxCount++
									}
								}
							}
						}
						sum += math.Pow(float64(boxCount)-mean, 2)
					}
				}
				variance := sum / float64(numBoxes)
				lacs[index] = variance / math.Pow(mean, 2)
			}(index, window)
		}
		wg.Wait()
		return lacs
	}
}
