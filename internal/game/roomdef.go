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

func (r *RoomDef) GetDescription() string {
	switch r.kind {
	case Boss:
		return "Fight a tremendous enemy!"
	case Stairs:
		return "Stairs"
	case Armory:
		switch r.size {
		case Medium:
			return "Increases your equipment levels by 1"
		case Large:
			return "Increases your equipment levels by 5"
		}
	case HealingShrine:
		switch r.size {
		case Small:
			return "Heals 25% of your health"
		case Medium:
			return "Heals 75% of your health"
		case Large:
			return "Heals 100% of your health"
		}
	case Combat:
		return "Engage with enemies to gain gold and XP!"
	case Well:
		return "Restores your equipment uses"
	case Treasure:
		switch r.size {
		case Small:
			return "Contains a peasant's pittance"
		case Medium:
			return "Contains a squire's savings"
		case Large:
			return "Contains a baron's bounty"
		case Huge:
			return "Contains wealth beyond measure"
		}
	case Library:
		return "A chance to enchant your equipment"
	case Curse:
		return "A chance to lose gold, equipment level, perk level, or dude level"
	case Trap:
		return "Watch your steppie!"
	}
	return "Unknown"
}

// Cost of the room scales with the story level
// With each level, the cost of the room increases by 25%
func GetRoomCost(kind RoomKind, size RoomSize, level int) int {
	cost := 0
	perLevelMultiplier := 0.25

	switch kind {
	case Armory:
		switch size {
		case Medium:
			cost = 100
		case Large:
			cost = 500
		}
	case HealingShrine:
		switch size {
		case Small:
			cost = 50
		case Medium:
			cost = 200
		case Large:
			cost = 500
		}
	case Well:
		cost = 250
	case Treasure:
		switch size {
		case Small:
			cost = 100
		case Medium:
			cost = 250
		case Large:
			cost = 500
		case Huge:
			cost = 1000
		}
	case Library:
		cost = 250
	}

	return cost + int(float64(cost)*float64(level)*perLevelMultiplier)
}
