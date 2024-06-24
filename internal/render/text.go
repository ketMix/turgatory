package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
)

type TextOptions struct {
	Screen *ebiten.Image
	Font   assets.Font
	Color  color.Color
	GeoM   ebiten.GeoM
}

func DrawText(o *TextOptions, txt string) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(o.Color)
	op.LineSpacing = o.Font.LineHeight
	op.GeoM.Concat(o.GeoM)
	text.Draw(o.Screen, txt, o.Font.Face, op)
}

func NewTextDrawer(o TextOptions) *TextDrawer {
	return &TextDrawer{Options: o}
}

type TextDrawer struct {
	Options TextOptions
}

func (t *TextDrawer) Draw(txt string, x, y float64) {
	t.Options.GeoM.Reset()
	t.Options.GeoM.Translate(x, y)
	DrawText(&t.Options, txt)
}
