package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type GameStateLose struct {
	wobbler float64
}

func (s *GameStateLose) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)
	g.audioController.PauseRoomTracks()
	g.audioController.PlaySfx("loss", 0.5, 0.0)
}
func (s *GameStateLose) End(g *Game) {
}
func (s *GameStateLose) Update(g *Game) GameState {
	s.wobbler += 0.05
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || (len(g.releasedTouchIDs) > 0 && inpututil.IsTouchJustReleased(g.releasedTouchIDs[0])) {
		return &GameStatePre{}
	}
	return nil
}
func (s *GameStateLose) Draw(g *Game, screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	opts := &render.TextOptions{
		Screen: screen,
		Font:   assets.DisplayFont,
		Color:  color.Black,
		GeoM:   geom,
	}

	w, h := text.Measure("GAME OVER", opts.Font.Face, opts.Font.LineHeight)
	w *= 4
	h *= 4

	opts.GeoM.Translate(-w/2, -h/2)
	opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.1)
	opts.GeoM.Translate(w/2, h/2)
	opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2+h)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*20)

	opts.GeoM.Translate(-10, -10)
	opts.Color = color.NRGBA{50, 0, 0, 200}
	render.DrawText(opts, "GAME OVER")
	opts.GeoM.Translate(20, 20)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*2)
	render.DrawText(opts, "GAME OVER")
	opts.Color = assets.ColorGameOver
	opts.GeoM.Translate(-10, -10)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*2)
	render.DrawText(opts, "GAME OVER")

	opts.Font = assets.BodyFont
	y := 0.0
	{
		opts.GeoM.Reset()
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("all ur dudes r ded :((", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/2)-h/2)
		y = float64(screen.Bounds().Dy()/2) - h/2
		render.DrawText(opts, "all ur dudes r ded :((")
	}

	{
		opts.GeoM.Reset()
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("Press SPACE to try again!", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, y+h)
		render.DrawText(opts, "Press SPACE to try again!")
	}
}
