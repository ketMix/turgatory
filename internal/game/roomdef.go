package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/internal/render"
)

// RoomDef is used for providing the visuals for placing a room. It renders a stack at a pitch of 1 to a new image and stores that image for use during rendering.
type RoomDef struct {
	kind  RoomKind
	size  RoomSize
	image *ebiten.Image
}

var RoomDefs = make(map[string]*RoomDef)

func GetRoomDef(kind RoomKind, size RoomSize) *RoomDef {
	key := fmt.Sprintf("%s_%s", kind.String(), size.String())
	if r, ok := RoomDefs[key]; ok {
		return r
	}

	stack, err := render.NewStack(fmt.Sprintf("rooms/%s", size.String()), kind.String(), "")
	if err != nil {
		panic(err)
	}
	/*stack.SetOriginToCenter()
	stack.SetRotation(math.Pi / 8)
	stack.SetPosition(10, 4)*/
	stack.SetPosition(0, float64(stack.Height()/4))

	img := ebiten.NewImage(stack.Width() /*+4*/, int(float64(stack.Height())*1.25))
	img.Clear() // Just in case...

	o := render.Options{
		Screen: img,
		Pitch:  1,
	}

	stack.Draw(&o)

	r := &RoomDef{
		kind:  kind,
		size:  size,
		image: img,
	}

	RoomDefs[key] = r
	return r
}
