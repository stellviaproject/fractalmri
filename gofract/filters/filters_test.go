package filters

import (
	"fmt"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"image"
	"image/png"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stellviaproject/colorize"
	arch "github.com/stellviaproject/pipfil-arch"
)

var (
	TestFolder                = "./testing"
	TestOutputFolder          = path.Join(TestFolder, "output")
	TestFileName              = path.Join(TestFolder, "TCGA_HT_A61A_20000127_28.tif")
	TestDecodedFileName       = path.Join(TestOutputFolder, "decoded.png")
	TestDenoisedFileName      = path.Join(TestOutputFolder, "denoised.png")
	TestGrayConvertedFileName = path.Join(TestOutputFolder, "grayconverted.png")
	TestMeasureFileName       = path.Join(TestOutputFolder, "measure.png")
	TestMFSFileName           = path.Join(TestOutputFolder, "mfs.png")
	TestUmbralsFolder         = path.Join(TestOutputFolder, "umbrals")
	TestBinsFolder            = path.Join(TestOutputFolder, "bins")
	TestObjsFolder            = path.Join(TestOutputFolder, "objs")
	TestFDsFolder             = path.Join(TestOutputFolder, "fds")
	TestLogLogsFolder         = path.Join(TestOutputFolder, "loglogs")
)

func TestInit(t *testing.T) {
	os.Mkdir(TestOutputFolder, os.ModePerm)
	os.Mkdir(TestUmbralsFolder, os.ModePerm)
	os.Mkdir(TestBinsFolder, os.ModePerm)
	os.Mkdir(TestObjsFolder, os.ModePerm)
	os.Mkdir(TestFDsFolder, os.ModePerm)
	os.Mkdir(TestLogLogsFolder, os.ModePerm)
}

func TestDecoder(t *testing.T) {
	input := pipes.NewBufferPipe("input", cfg.GetCFG())
	decoded := arch.NewPipe("decoded", []image.Image{}, 1)
	filter := NewDecoderFilter(input, decoded)
	model := arch.NewModel(arch.WithFilters(filter), arch.WithPipes(input), arch.WithPipes(decoded))
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	img := model.Call(arch.WithInput(data))[0].([]image.Image)[0]
	model.Stop()
	file, err := os.OpenFile(TestDecodedFileName, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func TestDenoiser(t *testing.T) {
	cfg := cfg.GetCFG()
	input := pipes.NewBufferPipe("buffer-bytes", cfg)
	decoded := pipes.NewImagePipe("decoded-image", cfg)
	denoised := pipes.NewImagePipe("denoised-image", cfg)
	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg, decoded, denoised)
	model := arch.NewModel(arch.WithFilters(decoder, denoiser), arch.WithPipes(input), arch.WithPipes(denoised))
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	img := model.Call(arch.WithInput(data))[0].(image.Image)
	file, err := os.OpenFile(TestDenoisedFileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func TestGrayConverter(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
	),
		arch.WithPipes(input),
		arch.WithPipes(grayconverted),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	img := model.Call(arch.WithInput(data))[0].(pipes.Image64).ToImage()
	file, err := os.OpenFile(TestGrayConvertedFileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func TestMeasure(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
	),
		arch.WithPipes(input),
		arch.WithPipes(measure),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	img := model.Call(arch.WithInput(data))[0].(pipes.Image64).ToImage()
	file, err := os.OpenFile(TestMeasureFileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func TestMFS(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
	),
		arch.WithPipes(input),
		arch.WithPipes(mfs),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	mfsImg := model.Call(arch.WithInput(data))[0].(pipes.Image64)
	img := mfsImg.Normalized().Scaled(255).ToImage()
	file, err := os.OpenFile(TestMFSFileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
	file.Close()
	s := 1.0
	e := 2.0
	o := 0.1
	c := int((e-s)/o) + 1
	mfsColor := colorize.Colorize(s, e, o, colorize.GenColorList(c), [][]float64(mfsImg))
	file, err = os.OpenFile(path.Join(TestOutputFolder, "mfs-color.png"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, mfsColor); err != nil {
		panic(err)
	}
	file.Close()
}

func TestUmbralizer(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())
	umbrals := arch.NewPipe("umbrals-image64[]", []pipes.Image64{}, 10)
	umbralList := arch.NewPipe("umbral-list", []*pipes.Umbral{}, 10)

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)
	umbralizer := NewUmbralizerFilter(cfg.GetCFG(), mfs, umbrals, umbralList)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
		umbralizer,
	),
		arch.WithPipes(input),
		arch.WithPipes(umbrals, umbralList),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	output := model.Call(arch.WithInput(data))
	model.Stop()
	umbImgLs := output[0].([]pipes.Image64)
	umbLs := output[1].([]*pipes.Umbral)
	os.Mkdir(TestUmbralsFolder, os.ModePerm)
	wg := sync.WaitGroup{}
	for i := 0; i < len(umbImgLs); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			file, err := os.OpenFile(path.Join(TestUmbralsFolder, fmt.Sprintf("u-%v.png", umbLs[i])), os.O_WRONLY|os.O_CREATE, 0664)
			if err != nil {
				panic(err)
			}
			if err := png.Encode(file, umbImgLs[i].Normalized().Scaled(255).ToImage()); err != nil {
				panic(err)
			}
			fmt.Println(umbLs[i])
		}(i)
	}
	wg.Wait()
}

func TestBinarizer(t *testing.T) {
	//cfg := config.Config
	//cfg.Umbral = []float64{1.16, 1.4}
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())
	umbrals := pipes.NewImage64Pipe("umbrals-image64", cfg.GetCFG())
	umbralList := pipes.NewUmbralListPipe("umbral-items", cfg.GetCFG())
	binarized := pipes.NewImageBinPipe("binary-umbral", cfg.GetCFG())

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)
	umbralizer := NewUmbralizerFilter(cfg.GetCFG(), mfs, umbrals, umbralList)
	binarizer := NewBinarizeFilter(cfg.GetCFG(), umbrals, binarized)
	type BinUmb struct {
		Bin pipes.ImageBin
		Umb *pipes.Umbral
	}
	joinerOut := arch.NewPipe("joinerOut", []*BinUmb{}, 10)
	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(bins []pipes.ImageBin, umbs []*pipes.Umbral) []*BinUmb {
			binUmbLs := []*BinUmb{}
			for i := 0; i < len(bins); i++ {
				binUmbLs = append(binUmbLs, &BinUmb{
					Bin: bins[i],
					Umb: umbs[i],
				})
			}
			return binUmbLs
		},
		arch.WithPipes(binarized, umbralList),
		arch.WithPipes(joinerOut),
		arch.WithLens(arch.NewLen(binarized, umbrals)),
	)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
		umbralizer,
		binarizer,
		joiner,
	),
		arch.WithPipes(input),
		arch.WithPipes(joinerOut),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	output := model.Call(arch.WithInput(data))[0].([]*BinUmb)
	model.Stop()
	os.Mkdir(TestBinsFolder, os.ModePerm)
	wg := sync.WaitGroup{}
	for i := 0; i < len(output); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filePath := path.Join(TestBinsFolder, fmt.Sprintf("b-%v.png", output[i].Umb))
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0664)
			if err != nil {
				panic(err)
			}
			if err := png.Encode(file, output[i].Bin.ToImage()); err != nil {
				panic(err)
			}
			fmt.Println(output[i].Umb)
		}(i)
	}
	wg.Wait()
}

func TestObjectProd(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())
	umbrals := pipes.NewImage64Pipe("umbrals-image64", cfg.GetCFG())
	umbralList := pipes.NewUmbralListPipe("umbral-items", cfg.GetCFG())
	binarized := pipes.NewImageBinPipe("binary-umbral", cfg.GetCFG())
	objectSlice := arch.NewPipe("object-slice", []*pipes.ImageObject{}, 10)

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)
	umbralizer := NewUmbralizerFilter(cfg.GetCFG(), mfs, umbrals, umbralList)
	binarizer := NewBinarizeFilter(cfg.GetCFG(), umbrals, binarized)
	objectprod := NewObjectFilter(cfg.GetCFG(), binarized, objectSlice)

	type ObjUmb struct {
		Objs []*pipes.ImageObject
		Umb  *pipes.Umbral
	}
	joinerOut := arch.NewPipe("joinerOut", []*ObjUmb{}, 10)
	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(objs [][]*pipes.ImageObject, umbs []*pipes.Umbral) []*ObjUmb {
			binUmbLs := []*ObjUmb{}
			for i := 0; i < len(objs); i++ {
				binUmbLs = append(binUmbLs, &ObjUmb{
					Objs: objs[i],
					Umb:  umbs[i],
				})
			}
			return binUmbLs
		},
		arch.WithPipes(objectSlice, umbralList),
		arch.WithPipes(joinerOut),
		arch.WithLens(arch.NewLen(objectSlice, umbrals)),
	)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
		umbralizer,
		binarizer,
		objectprod,
		joiner,
	),
		arch.WithPipes(input),
		arch.WithPipes(joinerOut),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	output := model.Call(arch.WithInput(data))[0].([]*ObjUmb)
	model.Stop()
	os.Mkdir(TestObjsFolder, os.ModePerm)
	wg := sync.WaitGroup{}
	for i := 0; i < len(output); i++ {
		objs := output[i].Objs
		umb := output[i].Umb
		objsFolderPath := path.Join(TestObjsFolder, fmt.Sprintf("%v", umb))
		os.Mkdir(objsFolderPath, os.ModePerm)
		for j := 0; j < len(objs); j++ {
			wg.Add(1)
			go func(j int, obj *pipes.ImageObject, umb *pipes.Umbral) {
				defer wg.Done()
				filePath := path.Join(TestObjsFolder, fmt.Sprintf("%v", umb), fmt.Sprintf("0-%d.png", j))
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0664)
				if err != nil {
					panic(err)
				}
				if err := png.Encode(file, obj.Image.ToImage()); err != nil {
					panic(err)
				}
				fmt.Println(umb)
			}(j, objs[j], umb)
		}
	}
	wg.Wait()
}

func TestFractalDIM(t *testing.T) {
	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())
	umbrals := pipes.NewImage64Pipe("umbrals-image64", cfg.GetCFG())
	umbralList := pipes.NewUmbralListPipe("umbral-items", cfg.GetCFG())
	binarized := pipes.NewImageBinPipe("binary-umbral", cfg.GetCFG())
	objects := pipes.NewObjectPipe("objects", cfg.GetCFG())
	fractaldim := pipes.NewFractalPipe("fractaldim", cfg.GetCFG())

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)
	umbralizer := NewUmbralizerFilter(cfg.GetCFG(), mfs, umbrals, umbralList)
	binarizer := NewBinarizeFilter(cfg.GetCFG(), umbrals, binarized)
	objectprod := NewObjectFilter(cfg.GetCFG(), binarized, objects)
	fractalizer := NewFractalFilter(cfg.GetCFG(), objects, fractaldim)

	type FracUmb struct {
		Fracs []*pipes.FractalDim
		Objs  []*pipes.ImageObject
		Umb   *pipes.Umbral
	}

	objJoinerOut := arch.NewPipe("object-slice", []*pipes.ImageObject{}, 10)
	objectJoiner := arch.NewFilterWithPipes(
		"objectJoiner",
		func(objs []*pipes.ImageObject) []*pipes.ImageObject {
			return objs
		},
		arch.WithPipes(objects),
		arch.WithPipes(objJoinerOut),
		arch.WithLens(arch.NewLen(objects, objects)),
	)

	fracJoinerOut := arch.NewPipe("frac-slice", []*pipes.FractalDim{}, 10)
	fracJoiner := arch.NewFilterWithPipes(
		"fracJoiner",
		func(fracs []*pipes.FractalDim) []*pipes.FractalDim {
			return fracs
		},
		arch.WithPipes(fractaldim),
		arch.WithPipes(fracJoinerOut),
		arch.WithLens(arch.NewLen(fractaldim, objects)),
	)

	joinerOut := arch.NewPipe("joinerOut", []*FracUmb{}, 10)
	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(fracs [][]*pipes.FractalDim, objs [][]*pipes.ImageObject, umbs []*pipes.Umbral) []*FracUmb {
			binUmbLs := []*FracUmb{}
			for i := 0; i < len(fracs); i++ {
				binUmbLs = append(binUmbLs, &FracUmb{
					Fracs: fracs[i],
					Umb:   umbs[i],
					Objs:  objs[i],
				})
			}
			return binUmbLs
		},
		arch.WithPipes(fracJoinerOut, objJoinerOut, umbralList),
		arch.WithPipes(joinerOut),
		arch.WithLens(
			arch.NewLen(fracJoinerOut, umbrals),
			arch.NewLen(objJoinerOut, umbrals),
		),
	)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
		umbralizer,
		binarizer,
		objectprod,
		fractalizer,
		objectJoiner,
		fracJoiner,
		joiner,
	),
		arch.WithPipes(input),
		arch.WithPipes(joinerOut),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	output := model.Call(arch.WithInput(data))[0].([]*FracUmb)
	model.Stop()
	os.Mkdir(TestFDsFolder, os.ModePerm)
	wg := sync.WaitGroup{}
	for i := 0; i < len(output); i++ {
		objs := output[i].Objs
		umb := output[i].Umb
		fds := output[i].Fracs
		os.Mkdir(path.Join(TestFDsFolder, fmt.Sprintf("%v", umb)), os.ModePerm)
		for j := 0; j < len(objs); j++ {
			wg.Add(1)
			go func(j int, obj *pipes.ImageObject, umb *pipes.Umbral, fd *pipes.FractalDim) {
				defer wg.Done()
				file, err := os.OpenFile(path.Join(TestFDsFolder, fmt.Sprintf("%v", umb), fmt.Sprintf("o-%d-%f.png", j, fd.FD)), os.O_WRONLY|os.O_CREATE, 0664)
				if err != nil {
					panic(err)
				}
				if err := png.Encode(file, obj.Image.ToImage()); err != nil {
					panic(err)
				}
				fmt.Printf("U:(%f, %f) FD:%f\n", umb.Min, umb.Max, fd.FD)
			}(j, objs[j], umb, fds[j])
		}
	}
	wg.Wait()
}

func TestLogLog(t *testing.T) {

	input := pipes.NewBufferPipe("buffer-bytes", cfg.GetCFG())
	decoded := pipes.NewImagePipe("decoded-image", cfg.GetCFG())
	denoised := pipes.NewImagePipe("denoised-image", cfg.GetCFG())
	grayconverted := pipes.NewImage64Pipe("gray-image", cfg.GetCFG())
	measure := pipes.NewImage64Pipe("measure-image64", cfg.GetCFG())
	mfs := pipes.NewImage64Pipe("mfs-image64", cfg.GetCFG())
	umbrals := pipes.NewImage64Pipe("umbrals-image64", cfg.GetCFG())
	umbralList := pipes.NewUmbralListPipe("umbral-items", cfg.GetCFG())
	binarized := pipes.NewImageBinPipe("binary-umbral", cfg.GetCFG())
	objects := pipes.NewObjectPipe("objects", cfg.GetCFG())
	fractaldim := pipes.NewFractalPipe("fractaldim", cfg.GetCFG())
	loglog := pipes.NewImagePipe("loglog", cfg.GetCFG())

	decoder := NewDecoderFilter(input, decoded)
	denoiser := NewDenoiserFilter(cfg.GetCFG(), decoded, denoised)
	grayconverter := NewGrayConverterFilter(cfg.GetCFG(), denoised, grayconverted)
	measureftr := NewMeasureFilter(cfg.GetCFG(), grayconverted, measure)
	mfsftr := NewMultiFractal(cfg.GetCFG(), measure, mfs)
	umbralizer := NewUmbralizerFilter(cfg.GetCFG(), mfs, umbrals, umbralList)
	binarizer := NewBinarizeFilter(cfg.GetCFG(), umbrals, binarized)
	objectprod := NewObjectFilter(cfg.GetCFG(), binarized, objects)
	fractalizer := NewFractalFilter(cfg.GetCFG(), objects, fractaldim)
	loglogftr := NewLogLogFilter(fractaldim, loglog)

	type FracUmb struct {
		Fracs   []*pipes.FractalDim
		LogLogs []image.Image
		Objs    []*pipes.ImageObject
		Umb     *pipes.Umbral
	}

	loglogJoinerOut := arch.NewPipe("loglogOut", []image.Image{}, 10)
	loglogJoiner := arch.NewFilterWithPipes(
		"loglogJoiner",
		func(loglogs []image.Image) []image.Image {
			return loglogs
		},
		arch.WithPipes(loglog),
		arch.WithPipes(loglogJoinerOut),
		arch.WithLens(arch.NewLen(loglog, objects)),
	)

	objJoinerOut := arch.NewPipe("object-slice", []*pipes.ImageObject{}, 10)
	objectJoiner := arch.NewFilterWithPipes(
		"objectJoiner",
		func(objs []*pipes.ImageObject) []*pipes.ImageObject {
			return objs
		},
		arch.WithPipes(objects),
		arch.WithPipes(objJoinerOut),
		arch.WithLens(arch.NewLen(objects, objects)),
	)

	fracJoinerOut := arch.NewPipe("frac-slice", []*pipes.FractalDim{}, 10)
	fracJoiner := arch.NewFilterWithPipes(
		"fracJoiner",
		func(fracs []*pipes.FractalDim) []*pipes.FractalDim {
			return fracs
		},
		arch.WithPipes(fractaldim),
		arch.WithPipes(fracJoinerOut),
		arch.WithLens(arch.NewLen(fractaldim, objects)),
	)

	joinerOut := arch.NewPipe("joinerOut", []*FracUmb{}, 10)
	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(fracs [][]*pipes.FractalDim, objs [][]*pipes.ImageObject, loglogs [][]image.Image, umbs []*pipes.Umbral) []*FracUmb {
			binUmbLs := []*FracUmb{}
			for i := 0; i < len(fracs); i++ {
				binUmbLs = append(binUmbLs, &FracUmb{
					Fracs:   fracs[i],
					Umb:     umbs[i],
					Objs:    objs[i],
					LogLogs: loglogs[i],
				})
			}
			return binUmbLs
		},
		arch.WithPipes(fracJoinerOut, objJoinerOut, loglogJoinerOut, umbralList),
		arch.WithPipes(joinerOut),
		arch.WithLens(
			arch.NewLen(fracJoinerOut, umbrals),
			arch.NewLen(objJoinerOut, umbrals),
			arch.NewLen(loglogJoinerOut, umbrals),
		),
	)

	model := arch.NewModel(arch.WithFilters(
		decoder,
		denoiser,
		grayconverter,
		measureftr,
		mfsftr,
		umbralizer,
		binarizer,
		objectprod,
		fractalizer,
		loglogftr,
		objectJoiner,
		fracJoiner,
		loglogJoiner,
		joiner,
	),
		arch.WithPipes(input),
		arch.WithPipes(joinerOut),
	)
	model.Run()
	data, err := os.ReadFile(TestFileName)
	if err != nil {
		panic(err)
	}
	output := model.Call(arch.WithInput(data))[0].([]*FracUmb)
	model.Stop()
	os.Mkdir(TestLogLogsFolder, os.ModePerm)
	wg := sync.WaitGroup{}
	for i := 0; i < len(output); i++ {
		objs := output[i].Objs
		umb := output[i].Umb
		fds := output[i].Fracs
		loglogs := output[i].LogLogs
		os.Mkdir(path.Join(TestLogLogsFolder, fmt.Sprintf("%v", umb)), os.ModePerm)
		for j := 0; j < len(objs); j++ {
			wg.Add(1)
			go func(j int, obj *pipes.ImageObject, umb *pipes.Umbral, fd *pipes.FractalDim, loglog image.Image) {
				defer wg.Done()
				folder := path.Join(TestLogLogsFolder, fmt.Sprintf("%v", umb))
				file, err := os.OpenFile(path.Join(folder, fmt.Sprintf("o-%d-%f.png", j, fd.FD)), os.O_WRONLY|os.O_CREATE, 0664)
				if err != nil {
					panic(err)
				}
				if err := png.Encode(file, obj.Image.ToImage()); err != nil {
					panic(err)
				}
				file.Close()
				file, err = os.OpenFile(path.Join(folder, fmt.Sprintf("o-%d-%f-log.png", j, fd.FD)), os.O_CREATE|os.O_WRONLY, 0664)
				if err != nil {
					panic(err)
				}
				if err := png.Encode(file, loglog); err != nil {
					panic(err)
				}
				file.Close()
				fmt.Printf("U:(%f, %f) FD:%f\n", umb.Min, umb.Max, fd.FD)
			}(j, objs[j], umb, fds[j], loglogs[j])
		}
	}
	wg.Wait()
}
