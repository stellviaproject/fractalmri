package views

import (
	ctrl "fractalmri/app/frontend/application"
	"fractalmri/app/frontend/domain/model"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type SegmentationTool struct {
	app.Compo
}

func (sg *SegmentationTool) OnMount(ctx app.Context) {
	ctx.ObserveState("/has-sg-tool").Value(&model.M.HasSgTool)
}

func (sg *SegmentationTool) OnUpdate(ctx app.Context) {
	sgTool := app.Window().Get("document").Call("querySelector", ".segtool")
	if model.M.HasSgTool {
		sgTool.Get("classList").Call("add", "hide")
	} else {
		sgTool.Get("classList").Call("remove", "hide")
	}
}

func (sg *SegmentationTool) Render() app.UI {
	hideClass := "hide"
	if model.M.HasSgTool {
		hideClass = ""
	}
	return app.Div().Class("segtool", hideClass).Body(
		app.Div().Class("row").Body(
			app.Div().Class("col s2 m2 l2").Body(
				app.Button().Class("waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowSegResult).Body(
					app.I().Class("material-icons").Text("play_arrow"),
				),
			),
			app.Div().Class("col s5 m5 l5").Body(
				app.Div().Class("input-field").Body(
					app.Input().ID("min-seg-input").Max(2.0).Min(0.0).Step(0.1).Type("number").Class("validate").Value(model.M.SegToolMin).OnChange(ctrl.Ctrl.HandleSgOnChange),
				),
			),
			app.Div().Class("col s5 m5 l5").Body(
				app.Div().Class("input-field").Body(
					app.Input().ID("max-seg-input").Max(2.0).Min(0.0).Step(0.1).Type("number").Class("validate").Value(model.M.SegToolMax).OnChange(ctrl.Ctrl.HandleSgOnChange),
				),
			),
		),
	)
}
