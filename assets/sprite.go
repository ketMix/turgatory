package assets

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var sprites = make(map[string]*Sprite)

type Sprite struct {
	Image *ebiten.Image
}

func LoadSprite(name string) (*Sprite, error) {
	if sprite, ok := sprites[name]; ok {
		return sprite, nil
	}

	b, err := FS.ReadFile(name + ".png")
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	eimg := ebiten.NewImageFromImage(img)

	sprite := &Sprite{
		Image: eimg,
	}

	sprites[name] = sprite

	return sprite, nil
}
