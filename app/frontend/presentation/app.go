package presentation

import (
	//ctrl "fractalmri/app/frontend/application"
	"fractalmri/app/frontend/application/client"
	"fractalmri/app/frontend/presentation/views"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func NewApp(resources ...string) *app.Handler {
	client.Init(app.Window().URL())
	//No eliminar esto
	home := &views.HomeView{}
	app.Route("/", home)
	app.RunWhenOnBrowser()
	a := &app.Handler{
		Title: "FractalMRI",
		Styles: []string{
			"/web/assets/materialize.css",
			"/web/assets/materialize.min.css",
			"/web/assets/styles.css",
		},
		Scripts: []string{
			"/web/assets/materialize.js",
			"/web/assets/materialize.min.js",
			"/web/assets/jquery-3.7.1.min.js",
			"/web/assets/jquery-3.7.1.js",
			"/web/assets/chart.js",
		},
		CacheableResources: resources,
	}
	return a
}
