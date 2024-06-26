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
		case StoryEnterNextActivity:
			// We can always allow this to happen since the logic is triggered in RoomEnterActivity with index 7.
			nextStory := t.Stories[u.story.level+1]

			if u.dude.room != nil {
				u.dude.room.RemoveDudeFromCenter(u.dude)
				u.dude.room.RemoveDude(u.dude)
				u.dude.room = nil
			}
			u.story.RemoveDude(u.dude)
			nextStory.AddDude(u.dude)
		case RoomCombatActivity:
			u.dude.Trigger(EventCombatRoom{room: u.room, dude: u.dude})
		case RoomEnterActivity:
			if u.dude.room != nil {
				u.dude.Trigger(EventLeaveRoom{room: u.dude.room, dude: u.dude})
				u.dude.room = nil
			}
			u.dude.room = u.room
			u.dude.Trigger(EventEnterRoom{room: u.room, dude: u.dude})
			// If it's the last room, then move upwards and go poof (unless we're coming from stairs or are entering the tower for the first time).
			if u.dude.activity != StairsFromDown && u.dude.activity != FirstEntering && u.room.index == 7 {
				u.dude.activity = StairsToUp
			}
		case RoomCenterActivity:
			u.dude.Trigger(EventCenterRoom{room: u.room, dude: u.dude})
		case RoomEndActivity:
			u.dude.Trigger(EventEndRoom{room: u.room, dude: u.dude})
			// Check if the given room is the last room in the story.
			if u.room.story.rooms[6] == u.room {
				if u.room.story.level == len(t.Stories) {
					fmt.Println("final level!! we made it")
					u.dude.activity = Idle
				} else if u.room.story.level < len(t.Stories)-1 {
					nextStory := t.Stories[u.room.story.level+1]
					if !nextStory.open {
						fmt.Println("next not open! set tower to teleporting!")
						u.dude.activity = Idle
					} else {
						fmt.Println("next is open, let's goooo")
					}
				}

				//u.dude.activity = Idle

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
		o.DrawImageOptions.GeoM.Translate(0, -StoryHeight*o.Camera.Zoom()) // Transform our rendering, ofc
	}
}

// AddStory does as it says.
func (t *Tower) AddStory(s *Story) {
	s.level = len(t.Stories)
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
	//d.stack.HeightOffset = 20
	enter := TowerEntrance - d.variation
	if enter < TowerEntrance+0.1 {
		enter = TowerEntrance
	}
	d.SetPosition(story.PositionFromCenter(math.Pi/2, enter))
}

func (t *Tower) AddDudes(dudes ...*Dude) {
	for _, d := range dudes {
		t.AddDude(d)
	}
}
