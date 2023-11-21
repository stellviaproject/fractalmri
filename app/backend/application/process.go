package application

import (
	"fmt"
	"fractalmri/app/backend/domain/model"
	md "fractalmri/app/frontend/domain/model"
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/lib"
	"image"
	"path"
)

type Process struct {
	Image      image.Image
	Params     *md.RunModel
	Evaluation *lib.ModelEvaluation
	Err        error
}

func NewProcess(img image.Image, params *md.RunModel) *Process {
	return &Process{
		Params: params,
		Image:  img,
	}
}

func (p *Process) Run(u *model.UserModel) {
	config := cfg.GetCFG()
	config.WindowRatio = p.Params.Window
	config.Umbral = cfg.Range(p.Params.MinUmbral, p.Params.MaxUmbral, p.Params.StepUmbral)
	config.MinUmbral = p.Params.MinUmbral
	config.MaxUmbral = p.Params.MaxUmbral
	fdModel := lib.NewFDModel(config)
	p.Evaluation, p.Err = fdModel.Eval(p.Image.(*image.RGBA))
	p.Evaluation.Save(path.Join(model.C.StorePath, fmt.Sprintf("%d/%d", u.ID, p.Params.ImageID)))
	result := &md.ResultModel{
		ID:          p.Params.ImageID,
		FDs:         p.Evaluation.FDs,
		Umbrals:     p.Evaluation.Umbrals,
		LogLogCount: len(p.Evaluation.LogLogs),
		IsCompleted: true,
	}
	u.AddResult(result)
}
