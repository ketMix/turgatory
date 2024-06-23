package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/assets"
)

type Sprite struct {
	Positionable
	Originable
	Rotateable
	image *ebiten.Image
}

func NewSprite(name string) (*Sprite, error) {
	dataSprite, err := assets.LoadSprite(name)
	if err != nil {
		return nil, err
	}
	sprite := &Sprite{}
	sprite.image = dataSprite.Image
	return sprite, nil
}

func NewSubSprite(dataSprite *Sprite, x, y, w, h int) (*Sprite, error) {
	sprite := &Sprite{}
	sprite.image = dataSprite.image.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	return sprite, nil
}

func (s *Sprite) Draw(o *Options) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Translate(s.Origin())
	opts.GeoM.Translate(s.Position())

	ox, oy := s.Origin()
	opts.GeoM.Translate(-ox, -oy)
	opts.GeoM.Rotate(s.Rotation())
	opts.GeoM.Translate(ox, oy)

	opts.GeoM.Concat(o.DrawImageOptions.GeoM)

	o.Screen.DrawImage(s.image, opts)
}

func (s *Sprite) Update() {
	// ???
}
