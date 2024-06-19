package assets

import (
	"bytes"
	"fmt"

	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var stacks = make(map[string]*Stack)

type Stack struct {
	Staxie
	Image *ebiten.Image
}

func LoadStack(name string) (*Stack, error) {
	if stack, ok := stacks[name]; ok {
		return stack, nil
	}

	b, err := FS.ReadFile(name + ".png")
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	i, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	// Convert the image to an Ebiten image.
	eimg := ebiten.NewImageFromImage(i)

	// Read our staxie PNG data.
	staxie := &Staxie{}
	err = staxie.FromBytes(b)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", staxie)

	stack := &Stack{
		Staxie: *staxie,
		Image:  eimg,
	}

	stacks[name] = stack

	return stack, nil
}
