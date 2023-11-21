package views

import (
	"fmt"
	ctrl "fractalmri/app/frontend/application"
	"fractalmri/app/frontend/domain/model"
	"sort"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type HomeView struct {
	app.Compo
	imageCount int
	loaded     bool
}

func (h *HomeView) OnMount(ctx app.Context) {
	//TODO: Puedes eliminar esto si desdea implementar un login en el futuro
	ctrl.Ctrl.Login(ctx) //<- solo esto
	results := ctrl.Ctrl.GetResults()
	for i := 0; i < len(results); i++ {
		model.M.Results[results[i].ID] = results[i]
	}
	ctrl.Ctrl.GetList()
	h.imageCount = len(model.M.Images)
	ctx.ObserveState("/loaded").Value(&h.loaded)
	ctx.ObserveState("/upload-err").Value(&model.M.UploadErr)
	ctx.ObserveState("/is-uploading").Value(&model.M.IsUploading)
	ctx.ObserveState("/allow-retry").Value(&model.M.AllowRetry)
	ctx.ObserveState("/image-list").Value(&model.M.Images)
	ctx.ObserveState("/image-count").Value(&h.imageCount)
	ctx.ObserveState("/selection").Value(&model.M.Images)
	ctx.ObserveState("/last-selection").Value(&model.M.LastSelection)
	ctx.ObserveState("/run-err").Value(&model.M.RunErr)
	ctx.ObserveState("/is-waiting-result").Value(&model.M.IsWaitingResult)
	ctx.ObserveState("/current-result").Value(&model.M.CurrentResult)
	ctx.ObserveState("/step").Value(&model.M.Step)
	ctx.Defer(func(ctx app.Context) {
		app.Window().Get("M").Call("AutoInit")
		elem := app.Window().GetElementByID("image-list")
		onOpenStart := app.FuncOf(func(this app.Value, args []app.Value) any {
			tool := app.Window().GetElementByID("list-tools")
			tool.Get("classList").Call("remove", "tool-out")
			tool.Get("classList").Call("add", "tool-in")
			return nil
		})
		onCloseStart := app.FuncOf(func(this app.Value, args []app.Value) any {
			tool := app.Window().GetElementByID("list-tools")
			tool.Get("classList").Call("remove", "tool-in")
			tool.Get("classList").Call("add", "tool-out")
			return nil
		})
		options := app.ValueOf(map[string]any{
			"onOpenStart":  "",
			"onCloseStart": "",
		})
		options.Set("onOpenStart", onOpenStart)
		options.Set("onCloseStart", onCloseStart)
		app.Window().Get("M").Get("Sidenav").Call("init", elem, options)
	})
	ctx.Dispatch(func(ctx app.Context) {
		ctx.SetState("/loaded", true)
	})
}

func (h *HomeView) Render() app.UI {
	sizes := GenSequence(2, 40, 1)
	imageMp := model.M.Images
	images := make([]*model.ImageModel, 0, len(imageMp))
	for _, img := range imageMp {
		images = append(images, img)
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].ID < images[j].ID
	})
	hideFix := ""
	if !h.loaded {
		hideFix = "hide"
	}
	return app.Div().Body(
		app.Div().ID("modal-upload").Class("modal").Body(
			app.Div().Class("modal-content").Body(
				app.H4().Text("Subiendo Archivo"),
				app.Div().Class("row").Body(
					app.Div().Class("col s3 m3 l3").Body(
						app.If(model.M.IsUploading, app.Div().Class("preloader-wrapper big active").Body(
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
						).Else(app.Div().Class("preloader-wrapper big active").Body(
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
					),
					app.Div().Class("col offset-s1 offset-m1 offset-l1 s8 m8 l8").Body(
						app.If(model.M.IsUploading,
							app.Text("Espera..."),
						).Else(
							app.Text(model.M.UploadErr),
						),
					),
				),
			),
			app.If(model.M.UploadErr != "", app.Div().Class("modal-footer").Body(
				app.If(model.M.AllowRetry,
					app.Button().Class("waves-effect waves-green btn-flat").Text("Reintentar").OnClick(ctrl.Ctrl.HandleUpload),
				),
				app.A().Href("#").Class("modal-close waves-effect waves-green btn-flat").Text("Cancelar"),
			)),
		),
		app.Ul().ID("image-list").Class("sidenav sidenav-fixed teal lighten-5", hideFix).Body(
			app.Li().Class("row center teal lighten-2 white-text").Body(
				app.H4().Class("col s8 m8 l8").Text("Imágenes"),
			),
			app.Range(images).Slice(func(i int) app.UI {
				if model.M.LastSelection == images[i].ID {
					return app.Li().Body(
						app.Label().Class("teal lighten-4").Body(
							app.Input().Type("checkbox").Class("filled-in").ID(fmt.Sprintf("image-%d", images[i].ID)).OnClick(ctrl.Ctrl.HandleSelect),
							app.Span().Text(images[i].Caption),
						),
					)
				}
				return app.Li().Class("white").Body(
					app.Label().Body(
						app.Input().Type("checkbox").Class("filled-in").ID(fmt.Sprintf("image-%d", images[i].ID)).OnClick(ctrl.Ctrl.HandleSelect),
						app.Span().Text(images[i].Caption),
					),
				)
			}),
			app.Li(),
		),
		app.Div().ID("list-tools").Class("tool", hideFix).Body(
			app.Div().Class("row teal lighten-2").Body(
				app.Div().Class("col s3 m3 l3").Body(
					app.Div().Class("waves-effect waves-light btn file-input-wrapper input-field").Body(
						app.I().Class("material-icons").Text("file_upload"),
						app.Input().ID("upload-file-input").Type("file").Disabled(model.M.IsUploading).OnInput(ctrl.Ctrl.HandleUpload),
					),
				),
				app.Div().Class("col s3 m3 l3").Body(
					app.Button().Class("waves-effect waves-light btn").Body(app.I().Class("material-icons center").Text("file_download")).OnClick(ctrl.Ctrl.HandleDownloadResult),
				),
				app.Div().Class("col s3 m3 l3").Body(
					app.Button().Class("waves-effect waves-light btn").Body(app.I().Class("material-icons center").Text("delete")).OnClick(ctrl.Ctrl.HandleOnDeleteImage),
				),
				app.Div().Class("col s3 m3 l3").Body(
					//TODO: ESTE BOTON NO ES FUNCIONAL PORQUE NO HAY CONEXION CON UNA BD POR EL MOMENTO
					app.Button().Class("waves-effect waves-light btn").Body(app.I().Class("material-icons center").Text("save")),
				),
			),
		),
		app.Main().Body(
			app.Div().Class("container").Body(
				app.Nav().Class("nav-extended").Body(
					app.Div().Class("nav-wrapper teal").Body(
						app.A().Href("#").Class("sidenav-trigger hide-on-large-only").Attr("data-target", "image-list").Body(
							app.I().Class("material-icons").Text("menu"),
						),
						app.A().Href("#").Class("brand-logo center").Text("FractalMRI"),
					),
				),
				app.Div().Class("row").Body(
					app.If(model.M.CurrentResult != -1,
						app.If(model.M.IsWaitingResult, app.Div().Class("viewer col s12 m12 l6 center scrollable").Body(
							app.Div().Class("row valign-wrapper").Body(
								app.Div().Class("col s4 m4 l4 center").Body(
									app.Div().Class("preloader-wrapper big active").Body(
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
								app.H4().Class("col s8 m8 l8 center").Text("Procesando imagen..."),
							),
						),
						).Else(
							&Viewer{},
						),
					).ElseIf(model.M.RunErr != "",
						app.Div().Class("valign-wrapper").Body(
							app.Button(),
							app.Button(),
						),
					).ElseIf(model.M.LastSelection != -1, app.Div().Class("viewer col s12 m12 l6 center scrollable").Body(
						app.Div().Class("row valign-wrapper").Body(
							app.Div().Class("viewer col s12 m12 l6 center scrollable").Body(
								app.Img().Class("img-selection center").Src(ctrl.Ctrl.LastSelectionURL()),
							),
						),
					)).Else(
						app.Div().Class("viewer col s12 m12 l6 center").Body(
							app.Div().Class("row valign-wrapper").Body(
								app.H4().Class("col s12 m12 l12 center").Text("No hay nada seleccionado."),
							),
						),
					),
					&Results{},
				),
				app.Div().Class("main-tool").Body(
					app.Div().Class("row cyan darken-2").Body(
						app.Div().Class("col s2 m3 l2").Body(
							/*
								L=LastSelection != -1
								I=IsWaitingResult
								D=Disabled
								I  L  D  L->I
								0  0  1   1
								1  0  1   1
								0  1  0   0
								1  1  1   1
								L -> I = !L || I
								Disabled = LastSelection==-1 || IsWaitingResult
							*/
							app.Button().Disabled(model.M.LastSelection == -1 || model.M.IsWaitingResult).Class("waves-effect waves-light btn").OnClick(ctrl.Ctrl.HandleRun).Body(
								app.Span().Class("hide-on-small-only").Text("Ejecutar"),
								app.I().Class("material-icons center hide-on-med-and-up").Text("play_arrow"),
								app.I().Class("material-icons right hide-on-small-only").Text("play_arrow"),
							),
						),
						app.Div().Class("col s3 m2 l2").Body(
							app.Div().Class("input-field cyan lighten-5").Body(
								app.Select().ID("input-mfs-window").Body(
									app.Range(sizes).Slice(func(i int) app.UI {
										if i == 0 {
											return app.Option().Value(sizes[i]).Text(fmt.Sprintf("MFS %d", sizes[i]))
										}
										return app.Option().Value(sizes[i]).Text(fmt.Sprintf("MFS %d", sizes[i]))
									}),
								),
							),
						),
						app.Div().Class("col s2 m2 l2").Body(
							app.Div().Class("input-field cyan lighten-5").Body(
								app.Input().ID("min-umbral-input").Max(2.0).Min(0.0).Step(model.M.Step).Type("number").Class("validate").Value(model.M.MinUmbral),
								app.Label().For("min-umbral-input").Text("Umbral mínimo").Class("hide-on-med-and-down"),
								app.Label().For("min-umbral-input").Text("U.Mínimo").Class("hide-on-large-only"),
							),
						),
						app.Div().Class("col s2 m2 l2").Body(
							app.Div().Class("input-field cyan lighten-5").Body(
								app.Input().ID("max-umbral-input").Max(2.0).Min(0.0).Step(model.M.Step).Type("number").Class("validate").Value(model.M.MaxUmbral),
								app.Label().For("max-umbral-input").Text("Umbral máximo").Class("hide-on-med-and-down"),
								app.Label().For("max-umbral-input").Text("U.Máximo").Class("hide-on-large-only"),
							),
						),
						app.Div().Class("col s2 m2 l2").Body(
							app.Div().Class("input-field cyan lighten-5").Body(
								app.Input().ID("step-umbral-input").Max(1.0).Min(0.0).Step(0.05).Type("number").Class("validate").Value(model.M.Step).OnChange(ctrl.Ctrl.HandleOnChangeStep),
								app.Label().For("step-umbral-input").Text("Variación"),
							),
						),
					),
				),
			),
		),
	)
}

func GenSequence(start, end, step int) []int {
	slice := []int{}
	for i := start; i < end; i += step {
		slice = append(slice, i)
	}
	return slice
}
