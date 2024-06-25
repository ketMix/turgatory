package game

import (
	"github.com/kettek/ebijam24/internal/render"
)

type Prefab struct {
	X, Y   float64
	ox, oy float64
	stack  *render.Stack
	vgroup *render.VGroup
	width  float64
	depth  float64
}

func NewPrefab(stack *render.Stack) *Prefab {
	p := &Prefab{}

	d := stack.SliceCount()
	h := stack.Height() + d*2
	w := stack.Width() + d*2
	p.width = float64(w)
	p.depth = float64(h)

	p.stack = stack
	p.stack.SetOriginToCenter()
	p.stack.SetPosition(p.width/4, p.depth/4) // Render stairs at offset within vgroup framebuffers
	p.vgroup = render.NewVGroup(w, h, d)

	return p
}

func (p *Prefab) Update() {
	p.stack.Update()
}

func (p *Prefab) Draw(o *render.Options) {
	p.vgroup.Clear()

	opts := &render.Options{
		Screen: o.Screen,
		Pitch:  o.Pitch,
		VGroup: p.vgroup,
	}

	// We can't use the camera's own functionality, so we do it ourselves here.
	opts.DrawImageOptions.GeoM.Translate(-p.width/2, -p.depth/2)
	opts.DrawImageOptions.GeoM.Rotate(o.TowerRotation)
	opts.DrawImageOptions.GeoM.Translate(p.width/2, p.depth/2)
	p.stack.Draw(opts)

	// Reset and draw our vgroup.
	opts.DrawImageOptions.GeoM.Reset()
	o.Camera.Transform(opts)

	// Get our rotated vgroup...
	fakeOpts := render.Options{}
	fakeOpts.DrawImageOptions.GeoM.Scale(o.Camera.Zoom(), o.Camera.Zoom())
	fakeOpts.DrawImageOptions.GeoM.Translate(-p.X, -p.Y)
	fakeOpts.DrawImageOptions.GeoM.Rotate(o.TowerRotation)
	fakeOpts.DrawImageOptions.GeoM.Translate(p.X, p.Y)

	tx := fakeOpts.DrawImageOptions.GeoM.Element(0, 2)
	ty := fakeOpts.DrawImageOptions.GeoM.Element(1, 2)

	// Offset our final render position
	opts.DrawImageOptions.GeoM.Translate(tx, ty)

	p.vgroup.Draw(opts)
}

func (p *Prefab) Position() (float64, float64) {
	return p.X, p.Y
}

func (p *Prefab) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

func (p *Prefab) Origin() (float64, float64) {
	return p.ox, p.oy
}

func (p *Prefab) SetOrigin(x, y float64) {
	p.ox = x
	p.oy = y
}
