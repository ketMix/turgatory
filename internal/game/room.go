package game

import (
	"fmt"

	"github.com/kettek/ebijam24/internal/render"
)

// RoomSize represents the different sizes of rooms, with 1 equating to 1/8th of a circle.
type RoomSize int

func (r RoomSize) String() string {
	switch r {
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Large:
		return "large"
	case Huge:
		return "huge"
	default:
		return "Unknown"
	}
}

const (
	Small  RoomSize = 1
	Medium RoomSize = 2
	Large  RoomSize = 3
	Huge   RoomSize = 4
)

// These origins are used to re-position a room "pie" image so that its center is in the appropriate place.
const (
	LargeOriginY = 44
	HugeOriginY  = 64
)

// RoomStairsEntrance is the distance from the center that a room's stairs is expected to be at.
const RoomStairsEntrance = 12
const RoomPath = 53
const TowerStairs = 60
const TowerEntrance = 80

// RoomKind is an enumeration of the different kinds of rooms in za toweru.
type RoomKind int

func (r *RoomKind) String() string {
	switch *r {
	case Armory:
		return "armory"
	case HealingShrine:
		return "healing"
	case Combat:
		//return "combat"
		return "template"
	case Well:
		return "well"
	default:
		return "Unknown"
	}
}

const (
	Empty RoomKind = iota
	// Armory provide... armor up? damage up? Maybe should be different types.
	Armory
	// Healing shrine heals the adventurers over time.
	HealingShrine
	// Combat is where it goes down. $$$ is acquired.
	Combat
	// Well restores magic items?
	Well
)

// Room is a room within a story of za toweru.
type Room struct {
	story *Story
	index int // Reference to the index within the story.
	kind  RoomKind
	size  RoomSize
	power int // ???

	stacks         render.Stacks
	walls          render.Stacks
	actorsInCenter []Actor
}

func NewRoom(size RoomSize, kind RoomKind) *Room {
	r := &Room{
		size: size,
		kind: kind,
	}

	stack, err := render.NewStack(fmt.Sprintf("rooms/%s", size.String()), kind.String(), "")
	if err != nil {
		panic(err)
	}
	r.stacks.Add(stack)

	// Add our walls.
	for j := 0; j < 3; j++ {
		for i := 0; i < 8; i++ {
			stack, err := render.NewStack(fmt.Sprintf("walls/%s", size.String()), "template", "base")
			if err != nil {
				continue
			}
			stack.VgroupOffset = j * StoryWallHeight
			r.walls.Add(stack)
		}
	}

	return r
}

// Update updates the stuff in the room.
func (r *Room) Update() {
	r.stacks.Update()
}

// Draw our room bits and bobs.
func (r *Room) Draw(o *render.Options) {
	r.stacks.Draw(o)
	r.walls.Draw(o)
}

func (r *Room) AddActorToCenter(a Actor) {
	r.actorsInCenter = append(r.actorsInCenter, a)
}

func (r *Room) RemoveActorFromCenter(a Actor) {
	for i, actor := range r.actorsInCenter {
		if actor == a {
			r.actorsInCenter = append(r.actorsInCenter[:i], r.actorsInCenter[i+1:]...)
			return
		}
	}
}

func (r *Room) IsActorInCenter(a Actor) bool {
	for _, actor := range r.actorsInCenter {
		if actor == a {
			return true
		}
	}
	return false
}
