package filters

import (
	"fractalmri/gofract/cfg"
	"fractalmri/gofract/pipes"
	"image"
	"math"

	arch "github.com/stellviaproject/pipfil-arch"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"

	"gonum.org/v1/plot/vg/draw"
)

func NewLogLogFilter(fdInPipe, loglogImgOutPipe arch.Pipe) arch.Filter {
	return arch.NewFilterWithPipes(
		"loglog",
		func(fd *pipes.FractalDim) (image.Image, error) {
			defer func() {
				if err := recover(); err != nil {
					panic(err)
				}
			}()
			cfg := cfg.GetCFG()
			width, height := cfg.LogLogWidth, cfg.LogLogHeight

			logSizes, logCounts := fd.LogSizes, fd.LogMeasure
			// Crear un nuevo plot y establecer sus propiedades
			p := plot.New()
			p.Title.Text = "Gráfico Log-Log"
			p.X.Label.Text = "Log(1/Size)"
			p.Y.Label.Text = "Log(Count)"

			// Agregar un gráfico de dispersión utilizando los datos proporcionados
			var pts plotter.XYs
			if math.IsNaN(fd.FD) {
				pts = plotter.XYs{}
			} else {
				pts = make(plotter.XYs, len(logSizes))
			}
			for i := range pts {
				pts[i].X = logSizes[i]
				pts[i].Y = logCounts[i]
			}
			s, err := plotter.NewScatter(pts)
			if err != nil {
				return nil, err
			}
			s.GlyphStyle.Color = plotutil.Color(0)
			s.GlyphStyle.Radius = vg.Points(3)

			// Agregar el gráfico de dispersión a la trama
			p.Add(s)

			// Generar la imagen y regresarla
			c := vgimg.New(font.Length(width), font.Length(height))
			p.Draw(draw.New(c))
			return c.Image(), nil
		},
		arch.WithPipes(fdInPipe),
		arch.WithPipes(loglogImgOutPipe),
		arch.WithLens(),
	)
}
