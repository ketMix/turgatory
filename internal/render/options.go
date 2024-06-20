package render

import "github.com/hajimehoshi/ebiten/v2"

type Options struct {
	Screen           *ebiten.Image
	VGroup           *VGroup
	Camera           *Camera
	Pitch            float64
	DrawImageOptions ebiten.DrawImageOptions
}
