package game

import (
	"image/color"

	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type FloatingText struct {
	render.Originable
	render.Positionable
	YOffset     float64
	text        string
	birthtime   int
	lifetime    int
	speed       float64
	color       color.NRGBA
	origAlpha   uint8
	textOptions render.TextOptions
}

func MakeFloatingText(text string, clr color.NRGBA, lifetime int, speed float64) FloatingText {
	clr2 := color.NRGBA{R: clr.R, G: clr.G, B: clr.B, A: 0}
	return FloatingText{
		text:      text,
		color:     clr2,
		origAlpha: clr.A,
		lifetime:  lifetime,
		birthtime: lifetime,
		speed:     speed,
		textOptions: render.TextOptions{
			Font:  assets.BodyFont,
			Color: clr2,
		},
	}
}

func MakeFloatingTextFromDude(d *Dude, text string, color color.NRGBA, lifetime int, speed float64) FloatingText {
	t := MakeFloatingText(text, color, lifetime, speed)
	dx, dy := d.Position()
	dx -= TowerCenterX
	dy -= TowerCenterY
	t.SetOrigin(dx, dy)
	t.YOffset = float64(d.stack.SliceCount())
	return t
}

func (t *FloatingText) Alive() bool {
	return t.lifetime > 0
}

func (t *FloatingText) Update() {
	t.lifetime--

	// Fade in first 3 ticks.
	if t.birthtime-t.lifetime <= 3 {
		t.color.A = uint8((float64(t.birthtime-t.lifetime) / 3) * float64(t.origAlpha))
	}

	// Fade out in last 5 ticks.
	if t.lifetime <= 5 {
		t.color.A = uint8((float64(t.lifetime) / 5) * float64(t.origAlpha))
	}
	t.textOptions.Color = t.color

	x, y := t.Position()
	t.SetPosition(x, y-t.speed)
}

func (t *FloatingText) Draw(o *render.Options) {
	t.textOptions.GeoM.Reset()
	t.textOptions.GeoM.Translate(t.Origin())
	t.textOptions.GeoM.Rotate(o.TowerRotation)
	t.textOptions.GeoM.Translate(t.Position())
	t.textOptions.GeoM.Scale(o.Camera.Zoom(), o.Camera.Zoom())

	// Get our own x&y without any rotations.
	x1, y1 := t.textOptions.GeoM.Element(0, 2), t.textOptions.GeoM.Element(1, 2)
	y1 /= o.Camera.Zoom()                                       // This ensures the Y offset squashes/stretches with the zoom
	y1 += (o.Camera.TextOffset() - t.YOffset) * o.Camera.Zoom() // Uh... this sorta works
	t.textOptions.GeoM.Reset()
	t.textOptions.GeoM.Translate(x1, y1)

	// Get our passed in draw image geom's x&y.
	x2 := o.DrawImageOptions.GeoM.Element(0, 2)
	y2 := o.DrawImageOptions.GeoM.Element(1, 2)
	t.textOptions.GeoM.Translate(x2, y2)

	t.textOptions.Screen = o.Screen

	render.DrawText(&t.textOptions, t.text)
}
