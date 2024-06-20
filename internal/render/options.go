package render

import "github.com/hajimehoshi/ebiten/v2"

type Options struct {
	Screen           *ebiten.Image
	Camera           *Camera
	Pitch            float64
	DrawImageOptions ebiten.DrawImageOptions
}
