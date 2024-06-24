package game

import (
	"fmt"
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Tower is our glorious tower :o
type Tower struct {
	render.Positionable // I guess it's okay to re-use this in such a fashion.
	Stories             []*Story
	stairs              *Prefab // Stairs at the bottom of the tower
}

// NewTower creates a new tower.
func NewTower() *Tower {
	t := &Tower{}

	// Create the stairs prefab.
	t.stairs = NewPrefab(Must(render.NewStack("walls/stairs", "", "")))
	t.stairs.SetPosition(0, 60)
	t.stairs.vgroup.Debug = true

	return t
}

// Update the tower.
func (t *Tower) Update() {
	t.stairs.Update()
	// TODO: Should this only update "active" stories?
	var storyUpdates ActivityRequests
	for _, s := range t.Stories {
		s.Update(&storyUpdates)
	}
	for _, u := range storyUpdates {
		switch u := u.(type) {
		case RoomCombatActivity:
			u.dude.Trigger(EventCombatRoom{room: u.room, dude: u.dude})
		case RoomEnterActivity:
			if u.room == nil {
				fmt.Printf("%s is in an empty room\n", u.initiator.Name())
			} else {
				var level int
				if story := u.room.story; story != nil {
					level = story.level
				}
				fmt.Printf("%s in story %d is moving to %s %s\n", u.initiator.Name(), level, u.room.size.String(), u.room.kind.String())
			}
			if dude, ok := u.initiator.(*Dude); ok {
				if dude.room != nil {
					dude.Trigger(EventLeaveRoom{room: dude.room, dude: dude})
				}
				dude.Trigger(EventEnterRoom{room: u.room, dude: dude})
			}
		case RoomCenterActivity:
			if dude, ok := u.initiator.(*Dude); ok {
				dude.Trigger(EventCenterRoom{room: u.room, dude: dude})
			}
		case RoomEndActivity:
			if dude, ok := u.initiator.(*Dude); ok {
				dude.Trigger(EventEndRoom{room: u.room, dude: dude})
			}
			if u.room.index == 7 {
				// NOTE: triggering a portal should only happen _once_, need to assign some sort of tower state.
				if !u.room.story.open {
					// TODO: If the next story is not open, then create a portal stack and have the lil dudes walk into + dematerialize (fade out).
				} else {
					// TODO: If the next story is open, then move to center stairs, set dude state to GoUpStairs, upon which the success of will cause an Activity of "EnterStory", which will move the dude to the next story and set the ComeFromStairs state.
				}
				fmt.Println("END OF STORY, OH GOSH")
			}
		}
		u.Apply()
		if cb := u.Cb(); cb != nil {
			cb(true)
		}
	}
}

// Draw our glorious tower.
func (t *Tower) Draw(o *render.Options) {
	for _, s := range t.Stories {
		s.Draw(o)
		o.DrawImageOptions.GeoM.Translate(0, -StoryHeight*o.Camera.Zoom) // Transform our rendering, ofc
	}
	//t.stairs.Draw(o) // We draw the stairs first and allow the stories to be drawn overtop.
}

// AddStory does as it says.
func (t *Tower) AddStory(s *Story) {
	t.Stories = append(t.Stories, s)
}

// AddDude adds a new dude at the lowest story of the tower and assigns the dude's appropriate activity state.
func (t *Tower) AddDude(d *Dude) {
	if len(t.Stories) == 0 {
		return
	}
	story := t.Stories[0]
	story.AddDude(d)
	d.activity = FirstEntering
	d.stack.HeightOffset = 20
	d.SetPosition(story.PositionFromCenter(math.Pi/2, TowerEntrance+d.variation))
}

func (t *Tower) AddDudes(dudes ...*Dude) {
	for _, d := range dudes {
		t.AddDude(d)
	}
}
