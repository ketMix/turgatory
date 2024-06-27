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
	image        *ebiten.Image
	Scale        float64
	Transparency float32
}

func (s *Sprite) Size() (float64, float64) {
	return float64(s.image.Bounds().Dx()) * s.Scale, float64(s.image.Bounds().Dy()) * s.Scale
}

func (s *Sprite) Width() float64 {
	return float64(s.image.Bounds().Dx()) * s.Scale
}

func (s *Sprite) Height() float64 {
	return float64(s.image.Bounds().Dy()) * s.Scale
}

func (s *Sprite) SetStaxie(name, stackName string) error {
	staxie, err := assets.LoadStaxie(name)
	if err != nil {
		return err
	}
	stack, ok := staxie.GetStack(stackName)
	if !ok {
		return assets.ErrStackNotFound
	}
	anim, ok := stack.GetAnimation("base")
	if !ok {
		return assets.ErrAnimationNotFound
	}
	frame, ok := anim.GetFrame(0)
	if !ok {
		return assets.ErrFrameNotFound
	}
	slice, ok := frame.GetSlice(0)
	if !ok {
		return assets.ErrSliceNotFound
	}
	s.image = slice.Image
	return nil
}

func (s *Sprite) SetStaxieAnimation(name, stackName, animName string) error {
	staxie, err := assets.LoadStaxie(name)
	if err != nil {
		return err
	}
	stack, ok := staxie.GetStack(stackName)
	if !ok {
		return assets.ErrStackNotFound
	}
	anim, ok := stack.GetAnimation(animName)
	if !ok {
		return assets.ErrAnimationNotFound
	}
	frame, ok := anim.GetFrame(0)
	if !ok {
		return assets.ErrFrameNotFound
	}
	slice, ok := frame.GetSlice(0)
	if !ok {
		return assets.ErrSliceNotFound
	}
	s.image = slice.Image
	return nil
}

func NewSpriteFromStaxie(name string, stackName string) (*Sprite, error) {
	sprite := &Sprite{
		Scale: 1,
	}

	err := sprite.SetStaxie(name, stackName)
	if err != nil {
		return nil, err
	}

	return sprite, nil
}

func NewSpriteFromStaxieAnimation(name string, stackName string, animName string) (*Sprite, error) {
	sprite := &Sprite{
		Scale: 1,
	}

	err := sprite.SetStaxieAnimation(name, stackName, animName)
	if err != nil {
		return nil, err
	}

	return sprite, nil
}

func NewSprite(name string) (*Sprite, error) {
	dataSprite, err := assets.LoadSprite(name)
	if err != nil {
		return nil, err
	}
	sprite := &Sprite{
		Scale: 1,
	}
	sprite.image = dataSprite.Image
	return sprite, nil
}

func NewSubSprite(dataSprite *Sprite, x, y, w, h int) (*Sprite, error) {
	sprite := &Sprite{
		Scale: 1,
	}
	sprite.image = dataSprite.image.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	return sprite, nil
}

func (s *Sprite) Draw(o *Options) {
	opts := &ebiten.DrawImageOptions{}

	ox, oy := s.Origin()
	opts.GeoM.Translate(-ox, -oy)
	opts.GeoM.Rotate(s.Rotation())
	opts.GeoM.Translate(ox, oy)

	opts.GeoM.Scale(s.Scale, s.Scale)

	opts.GeoM.Translate(s.Origin())
	opts.GeoM.Translate(s.Position())

	opts.GeoM.Concat(o.DrawImageOptions.GeoM)

	if s.Transparency != 0 {
		opts.ColorScale.ScaleAlpha(1.0 - s.Transparency)
	}

	o.Screen.DrawImage(s.image, opts)
}

func (s *Sprite) Update() {
	// ???
}
