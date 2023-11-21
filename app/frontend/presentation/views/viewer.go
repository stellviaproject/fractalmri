package views

import (
	a "fractalmri/app/frontend/application"
	"fractalmri/app/frontend/domain/model"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Viewer struct {
	app.Compo
	IsWaitingMFS bool
	MfsLoadErr   string
}

func (vr *Viewer) OnMount(ctx app.Context) {
	ctx.ObserveState("/is-waiting-mfs").Value(&vr.IsWaitingMFS)
	ctx.ObserveState("/mfs-load-err").Value(&vr.MfsLoadErr)
	ctx.ObserveState("/current-result").Value(&model.M.CurrentResult)
	a.Ctrl.DrawMFS(model.M.CurrentResult, ctx)
}

// func (vr *Viewer) OnUpdate(ctx app.Context) {
// 	var currResult int
// 	ctx.GetState("/current-result", &currResult)
// 	if currResult != vr.resultID {
// 		canvas := app.Window().GetElementByID(vr.CanvasID)
// 		a.Ctrl.DrawMFS(canvas, ctx)
// 		vr.resultID = currResult
// 	}
// }

func (vr *Viewer) Render() app.UI {
	hideCanvasClass, hidePreloaders := "", ""
	if vr.IsWaitingMFS || vr.MfsLoadErr != "" {
		hideCanvasClass += "hide"
	} else {
		hidePreloaders = "hide"
	}
	hideWaitingClass := ""
	if !vr.IsWaitingMFS {
		hideWaitingClass += "hide"
	}
	hideMfsLdErrClass := ""
	if vr.MfsLoadErr == "" {
		hideMfsLdErrClass += "hide"
	}

	return app.Div().Class("viewer col s12 m12 l6 center scrollable").Body(
		&SegmentationTool{},
		app.Canvas().ID("canvas-vr").Class(hideCanvasClass, "mfs-canvas center"),
		app.Div().Class("row valign-wrapper", hidePreloaders).Body(
			app.Div().Class("col s4 m4 l4 center", hideWaitingClass).Body(
				app.Div().Class("mfs-preloader preloader-wrapper big active").Body(
					app.Div().Class("spinner-layer spinner-blue-only").Body(
						app.Div().Class("circle-clipper left").Body(
							app.Div().Class("circle"),
						),
						app.Div().Class("gap-patch").Body(
							app.Div().Class("circle"),
						),
						app.Div().Class("circle-clipper right").Body(
							app.Div().Class("circle"),
						),
					)),
			),
			app.H4().Class("col s8 m8 l8 center", hideWaitingClass).Text("Cargando resultados..."),
			app.Div().Class("col s4 m4 l4 center", hideMfsLdErrClass).Body(
				app.Div().Class("mfs-preloader preloader-wrapper big active", hideMfsLdErrClass).Body(
					app.Div().Class("spinner-layer spinner-red-only").Body(
						app.Div().Class("circle-clipper left").Body(
							app.Div().Class("circle"),
						),
						app.Div().Class("gap-patch").Body(
							app.Div().Class("circle"),
						),
						app.Div().Class("circle-clipper right").Body(
							app.Div().Class("circle"),
						),
					)),
			),
			app.H4().Class("col s8 m8 l8 center", hideMfsLdErrClass).Text(vr.MfsLoadErr),
		),
	)
}
