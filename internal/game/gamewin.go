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

type GameStateWin struct {
	wobbler float64
}

func (s *GameStateWin) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)
	g.audioController.PauseRoomTracks()
	g.audioController.PlaySfx("win", 0.5, 0.0)
}
func (s *GameStateWin) End(g *Game) {
}

func (s *GameStateWin) Update(g *Game) GameState {
	s.wobbler += 0.05
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || (len(g.releasedTouchIDs) > 0 && inpututil.IsTouchJustReleased(g.releasedTouchIDs[0])) {
		return &GameStatePre{}
	}
	g.camera.SetRotation(g.camera.Rotation() + 0.005)
	return nil
}

func (s *GameStateWin) Draw(g *Game, screen *ebiten.Image) {
	s.DrawRainbow(screen, "U   WINNE!1!")

	opts := &render.TextOptions{
		Screen: screen,
		Font:   assets.DisplayFont,
		Color:  color.Black,
	}

	opts.Color = assets.ColorDudeTitle
	opts.Font = assets.BodyFont
	y := 0.0
	{
		opts.GeoM.Reset()
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("ur dudes escaped :))", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/2)+10*4)
		y = float64(screen.Bounds().Dy()/2) + 20*4
		render.DrawText(opts, "ur dudes escaped :))")
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

func (s *GameStateWin) DrawRainbow(screen *ebiten.Image, t string) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	opts := &render.TextOptions{
		Screen: screen,
		Font:   assets.DisplayFont,
		Color:  color.Black,
		GeoM:   geom,
	}

	w, h := text.Measure(t, opts.Font.Face, opts.Font.LineHeight)
	w *= 4
	h *= 4

	opts.GeoM.Translate(-w/2, -h/2)
	opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.05)
	opts.GeoM.Translate(w/2, h/2)
	opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

	opts.GeoM.Translate(0, h*2) // Uh... this *2 is odd.

	// RAINBOW
	var colors = []color.NRGBA{
		{255, 0, 0, 255},
		{255, 127, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
		{75, 0, 130, 255},
		{238, 130, 238, 255},
	}
	ci := 0

	// Draw "shadow"
	oldGeoM := ebiten.GeoM{}
	oldGeoM.Concat(opts.GeoM)
	for _, r := range t {
		tx, _ := text.Measure(string(r), opts.Font.Face, opts.Font.LineHeight)
		tx *= 4

		opts.GeoM.Translate(-10, -10)
		opts.Color = color.NRGBA{10, 0, 0, 200}
		render.DrawText(opts, string(r))
		opts.GeoM.Translate(20, 20)
		render.DrawText(opts, string(r))
		opts.GeoM.Translate(-10, -10)

		opts.GeoM.Translate(tx, 0)
		opts.GeoM.Rotate(math.Cos(s.wobbler) * 0.003)
	}

	// Draw main text
	opts.GeoM = oldGeoM
	for _, r := range t {
		tx, _ := text.Measure(string(r), opts.Font.Face, opts.Font.LineHeight)
		tx *= 4

		opts.Color = colors[ci%len(colors)]
		render.DrawText(opts, string(r))

		opts.GeoM.Translate(tx, 0)
		opts.GeoM.Rotate(math.Cos(s.wobbler) * 0.003)

		ci++
	}
}
