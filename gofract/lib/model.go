package lib

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/filters"
	"fractalmri/gofract/pipes"
	"image"
	"image/png"
	"math"
	"os"
	"path"

	"github.com/stellviaproject/go-ia/knn"
	arch "github.com/stellviaproject/pipfil-arch"
)

type FDModel struct {
	model arch.Model
}

//export NewFDModel
func NewFDModel(cfg cfg.Configuration) *FDModel {
	input := pipes.NewImagePipe("input-image", cfg)
	denoiserOutput := pipes.NewImagePipe("denoiser-output", cfg)
	grayConverterOutput := pipes.NewImage64Pipe("gray-converter-output", cfg)
	measureOutput := pipes.NewImage64Pipe("measure-output", cfg)
	mfsOutput := pipes.NewImage64Pipe("mfs-output", cfg)
	umbralOutput := pipes.NewImage64Pipe("umbral-output", cfg)
	umbralListOutput := pipes.NewUmbralListPipe("umbral-list-output", cfg)
	binaryOutput := pipes.NewImageBinPipe("binary-output", cfg)
	fdOutput := pipes.NewFractalPipe("fd-output", cfg)
	loglogOutput := pipes.NewImagePipe("loglog-output", cfg)
	output := arch.NewPipe("final-output", &ModelEvaluation{}, cfg.Buffer)

	denoiserFtr := filters.NewDenoiserFilter(cfg, input, denoiserOutput)
	grayConverterFtr := filters.NewGrayConverterFilter(cfg, denoiserOutput, grayConverterOutput)
	measureFtr := filters.NewMeasureFilter(cfg, grayConverterOutput, measureOutput)
	mfsFtr := filters.NewMultiFractal(cfg, measureOutput, mfsOutput)
	umbralFtr := filters.NewUmbralizerFilter(cfg, mfsOutput, umbralOutput, umbralListOutput)
	binaryFtr := filters.NewBinarizeFilter(cfg, umbralOutput, binaryOutput)
	fractalFtr := filters.NewFractalFilter(cfg, binaryOutput, fdOutput)
	loglogFtr := filters.NewLogLogFilter(fdOutput, loglogOutput)

	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(fds []*pipes.FractalDim, loglogs []image.Image, umbrals []*pipes.Umbral, mfs pipes.Image64) *ModelEvaluation {
			var loglogsRGBA []*image.RGBA
			for i := range loglogs {
				loglogsRGBA = append(loglogsRGBA, loglogs[i].(*image.RGBA))
			}
			return &ModelEvaluation{
				FDs:     fds,
				LogLogs: loglogsRGBA,
				Umbrals: umbrals,
				MFS:     mfs,
			}
		},
		arch.WithPipes(
			fdOutput,
			loglogOutput,
			umbralListOutput,
			mfsOutput,
		),
		arch.WithPipes(output),
		arch.WithLens(
			arch.NewLen(fdOutput, umbralOutput),
			arch.NewLen(loglogOutput, umbralOutput),
		),
	)

	model := arch.NewModel(
		arch.WithFilters(
			denoiserFtr,
			grayConverterFtr,
			measureFtr,
			mfsFtr,
			umbralFtr,
			binaryFtr,
			fractalFtr,
			loglogFtr,
			joiner,
		),
		arch.WithPipes(input),
		arch.WithPipes(output),
	)
	model.Run()
	return &FDModel{
		model: model,
	}
}

func NewFDModelNoLogLog(cfg cfg.Configuration) *FDModel {
	input := pipes.NewImagePipe("input-image", cfg)
	denoiserOutput := pipes.NewImagePipe("denoiser-output", cfg)
	grayConverterOutput := pipes.NewImage64Pipe("gray-converter-output", cfg)
	measureOutput := pipes.NewImage64Pipe("measure-output", cfg)
	mfsOutput := pipes.NewImage64Pipe("mfs-output", cfg)
	umbralOutput := pipes.NewImage64Pipe("umbral-output", cfg)
	umbralListOutput := pipes.NewUmbralListPipe("umbral-list-output", cfg)
	binaryOutput := pipes.NewImageBinPipe("binary-output", cfg)
	fdOutput := pipes.NewFractalPipe("fd-output", cfg)
	output := arch.NewPipe("final-output", &ModelEvaluation{}, cfg.Buffer)

	denoiserFtr := filters.NewDenoiserFilter(cfg, input, denoiserOutput)
	grayConverterFtr := filters.NewGrayConverterFilter(cfg, denoiserOutput, grayConverterOutput)
	measureFtr := filters.NewMeasureFilter(cfg, grayConverterOutput, measureOutput)
	mfsFtr := filters.NewMultiFractal(cfg, measureOutput, mfsOutput)
	umbralFtr := filters.NewUmbralizerFilter(cfg, mfsOutput, umbralOutput, umbralListOutput)
	binaryFtr := filters.NewBinarizeFilter(cfg, umbralOutput, binaryOutput)
	fractalFtr := filters.NewFractalFilter(cfg, binaryOutput, fdOutput)

	joiner := arch.NewFilterWithPipes(
		"joiner",
		func(fds []*pipes.FractalDim, umbrals []*pipes.Umbral, mfs pipes.Image64) *ModelEvaluation {
			return &ModelEvaluation{
				FDs:     fds,
				LogLogs: make([]*image.RGBA, 0),
				Umbrals: umbrals,
				MFS:     mfs,
			}
		},
		arch.WithPipes(
			fdOutput,
			umbralListOutput,
			mfsOutput,
		),
		arch.WithPipes(output),
		arch.WithLens(
			arch.NewLen(fdOutput, umbralOutput),
		),
	)

	model := arch.NewModel(
		arch.WithFilters(
			denoiserFtr,
			grayConverterFtr,
			measureFtr,
			mfsFtr,
			umbralFtr,
			binaryFtr,
			fractalFtr,
			joiner,
		),
		arch.WithPipes(input),
		arch.WithPipes(output),
	)
	model.Run()
	return &FDModel{
		model: model,
	}
}

//export Eval
func (md *FDModel) Eval(img *image.RGBA) (*ModelEvaluation, error) {
	output := md.model.Call(arch.WithInput(img))[0]
	if md.model.HasErrs() {
		return nil, md.model.Errs()[0]
	}
	return output.(*ModelEvaluation), nil
}

type ModelEvaluation struct {
	FDs     []*pipes.FractalDim `json:"fds"`
	LogLogs []*image.RGBA       `json:"-"`
	Umbrals []*pipes.Umbral     `json:"umbrals"`
	MFS     pipes.Image64       `json:"-"`
}

//export GetPoints
func (ev *ModelEvaluation) GetPoints() knn.Point {
	point := knn.NewPoint(len(ev.FDs))
	for i, fd := range ev.FDs {
		if !math.IsNaN(fd.FD) {
			point[i] = fd.FD
		}
	}
	return point
}

func (ev *ModelEvaluation) Save(folder string) error {
	os.Mkdir(folder, os.ModePerm)
	type Data struct {
		FDs     []*pipes.FractalDim `json:"fds"`
		LogLogs []string            `json:"loglogs"`
		Umbrals []*pipes.Umbral     `json:"umbrals"`
		MFS     string              `json:"mfs"`
	}
	data := &Data{
		LogLogs: make([]string, 0, len(ev.LogLogs)),
		FDs:     ev.FDs,
		Umbrals: ev.Umbrals,
	}
	loglogDir := path.Join(folder, "loglog")
	os.Mkdir(loglogDir, os.ModePerm)
	for i := 0; i < len(ev.LogLogs); i++ {
		fileName := path.Join("loglog", fmt.Sprintf("%d.png", i))
		data.LogLogs = append(data.LogLogs, fileName)
		file, err := os.OpenFile(path.Join(folder, fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModeDevice|os.ModePerm)
		defer file.Close()
		if err != nil {
			return err
		}
		if err := png.Encode(file, ev.LogLogs[i]); err != nil {
			return err
		}
	}
	mfsFileName := path.Join(folder, "mfs.bin")
	data.MFS = "mfs.bin"
	file, err := os.OpenFile(mfsFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm|os.ModeDevice)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := binary.Write(file, binary.LittleEndian, int32(ev.MFS.Width())); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, int32(ev.MFS.Height())); err != nil {
		return err
	}
	for i := 0; i < ev.MFS.Width(); i++ {
		for j := 0; j < ev.MFS.Height(); j++ {
			if err := binary.Write(file, binary.LittleEndian, ev.MFS.At(i, j)); err != nil {
				return err
			}
		}
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataFile := path.Join(folder, "data.json")
	err = os.WriteFile(dataFile, dataJSON, os.ModeDevice|os.ModePerm)
	return err
}

func (ev *ModelEvaluation) Load(folder string) error {
	type Data struct {
		FDs     []*pipes.FractalDim `json:"fds"`
		LogLogs []string            `json:"loglogs"`
		Umbrals []*pipes.Umbral     `json:"umbrals"`
		MFS     string              `json:"mfs"`
	}
	dataFile := path.Join(folder, "data.json")
	dataJSON, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	var data Data
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return err
	}
	ev.FDs = data.FDs
	ev.Umbrals = data.Umbrals
	ev.LogLogs = make([]*image.RGBA, 0, len(data.LogLogs))
	for i := 0; i < len(data.LogLogs); i++ {
		fileName := path.Join(folder, data.LogLogs[i])
		file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeDevice|os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()
		img, err := png.Decode(file)
		if err != nil {
			return err
		}
		ev.LogLogs = append(ev.LogLogs, img.(*image.RGBA))
	}
	file, err := os.OpenFile(path.Join(folder, data.MFS), os.O_RDONLY, os.ModeDevice|os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	var (
		width  int32
		height int32
	)
	if err := binary.Read(file, binary.LittleEndian, &width); err != nil {
		return err
	}
	if err := binary.Read(file, binary.LittleEndian, &height); err != nil {
		return err
	}
	mfs := pipes.NewImage64(int(width), int(height))
	var pixel float64
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			if err := binary.Read(file, binary.LittleEndian, &pixel); err != nil {
				return err
			}
			mfs.Set(i, j, pixel)
		}
	}
	ev.MFS = mfs
	return nil
}
