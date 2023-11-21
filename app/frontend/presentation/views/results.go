package views

import (
	"fmt"
	ctrl "fractalmri/app/frontend/application"
	"fractalmri/app/frontend/domain/model"
	"sort"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Results struct {
	app.Compo
}

func (r *Results) OnMount(ctx app.Context) {
	ctx.ObserveState("/results").Value(&model.M.Results)
	ctx.ObserveState("/current-result").Value(&model.M.CurrentResult)
	ctx.ObserveState("/show-result").Value(&model.M.ShowResult)
}

func (r *Results) Render() app.UI {
	resultsMap := model.M.Results
	results := make([]*model.ResultModel, 0, len(resultsMap))
	for _, result := range resultsMap {
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})
	id := model.M.CurrentResult
	var current *model.ResultModel
	if id != -1 {
		current = model.M.Results[id]
	} else {
		current = model.NewResultModel(0)
	}
	loglogURLs := make([]string, 0, current.LogLogCount)
	for i := 0; i < current.LogLogCount; i++ {
		loglogURLs = append(loglogURLs, fmt.Sprintf("/loglog?id=%d&index=%d", current.ID, i))
	}
	return app.Div().Class("results col offset-l6").Body(
		app.Div().Class("results-side hide-on-large-only").Body(
			app.Button().Class("waves-effect waves-light btn modal-trigger").Attr("data-target", "results-modal").Body(
				app.I().Class("material-icons center").Text("view_list"),
			),
			app.Div().ID("results-modal").Class("modal").Body(
				app.Div().Class("modal-content").Body(
					app.Ul().Class("collection scrollable blue lighten-5").Body(
						app.If(model.M.ShowResult,
							app.Li().Class("modal-item collection-item valign-wrapper").Body(
								app.Button().ID("button-back").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowResultList).Body(
									app.I().Class("material-icons left").Text("arrow_back"),
									app.Text("Ver Resultados"),
								),
							),
							app.Li().Class("modal-item collection-item valign-wrapper").Body(
								app.Button().ID("button-segtool").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowSegTool).Body(
									app.Text("Segmentación"),
									app.I().Class("material-icons right").Text("content_cut"),
								),
							),
							app.Li().Class("modal-item collection-item valign-wrapper").Body(
								app.Button().ID("button-mfs").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowMFS).Body(
									app.Text("Espectro Multifractal"),
									app.I().Class("material-icons right").Text("image"),
								),
							),
							app.Range(loglogURLs).Slice(func(i int) app.UI {
								return app.Li().Class("modal-item collection-item valign-wrapper").Body(
									app.Button().ID(fmt.Sprintf("button-loglog-%d", i)).Class("btn-result waves-effect waves-light white btn").OnClick(ctrl.Ctrl.HandleShowFDChart).Body(
										app.Text(fmt.Sprintf("Umbral [%f, %f]", current.Umbrals[i].Min, current.Umbrals[i].Max)),
										app.I().Class("material-icons right").Text("show_chart"),
									),
								)
							}),
						).Else(
							app.Range(results).Slice(func(i int) app.UI {
								return app.Li().Class("collection-item valign-wrapper").Body(
									app.Label().Body(
										app.Input().Type("checkbox").Class("filled-in").ID(fmt.Sprintf("result-%d", results[i].ID)),
										app.Span().Text(fmt.Sprintf("Resultado %d", results[i].ID+1)),
									),
									app.Button().ID(fmt.Sprintf("button-%d", results[i].ID)).Class("btn-show waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowResult).Body(
										app.I().Class("material-icons center").Text("arrow_forward"),
									),
								)
							}),
						),
					),
				),
			),
		),
		app.Ul().Class("collection scrollable hide-on-med-and-down blue lighten-5").Body(
			app.If(model.M.ShowResult,
				app.Li().Class("collection-item valign-wrapper").Body(
					app.Button().ID("button-back").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowResultList).Body(
						app.I().Class("material-icons left").Text("arrow_back"),
						app.Text("Ver Resultados"),
					),
				),
				app.Li().Class("collection-item valign-wrapper").Body(
					app.Button().ID("button-segtool").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowSegTool).Body(
						app.Text("Segmentación"),
						app.I().Class("material-icons right").Text("content_cut"),
					),
				),
				app.Li().Class("collection-item valign-wrapper").Body(
					app.Button().ID("button-mfs").Class("btn-result waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowMFS).Body(
						app.Text("Espectro Multifractal"),
						app.I().Class("material-icons right").Text("image"),
					),
				),
				app.Range(loglogURLs).Slice(func(i int) app.UI {
					return app.Li().Class("collection-item valign-wrapper").Body(
						app.Button().ID(fmt.Sprintf("button-loglog-%d", i)).Class("btn-result waves-effect waves-light white btn").OnClick(ctrl.Ctrl.HandleShowFDChart).Body(
							app.Text(fmt.Sprintf("Umbral [%f, %f]", current.Umbrals[i].Min, current.Umbrals[i].Max)),
							app.I().Class("material-icons right").Text("show_chart"),
						),
					)
				}),
			).Else(
				app.Range(results).Slice(func(i int) app.UI {
					return app.Li().Class("collection-item valign-wrapper").Body(
						app.Label().Body(
							app.Input().Type("checkbox").Class("filled-in").ID(fmt.Sprintf("result-%d", results[i].ID)),
							app.Span().Text(fmt.Sprintf("Resultado %d", results[i].ID+1)),
						),
						app.Button().ID(fmt.Sprintf("button-%d", results[i].ID)).Class("btn-show waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleShowResult).Body(
							app.I().Class("material-icons center").Text("arrow_forward"),
						),
					)
				}),
			),
		),
	)
}
