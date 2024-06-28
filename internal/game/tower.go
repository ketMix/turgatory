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
	portalOpen          bool
	dudes               []*Dude
}

// NewTower creates a new tower.
func NewTower() *Tower {
	t := &Tower{}

	return t
}

// Update the tower.
func (t *Tower) Update(req *ActivityRequests) {
	// TODO: Should this only update "active" stories?
	var storyUpdates ActivityRequests
	for _, s := range t.Stories {
		s.Update(&storyUpdates)
	}
	for _, u := range storyUpdates {
		switch u := u.(type) {
		case DudeDeadActivity:
			// Propagate the dead dude up to the game.
			req.Add(u)
		case TowerLeaveActivity:
			t.RemoveDude(u.dude)
			// Send it up to the game so it can do logic to see if all dudes are gone.
			req.Add(u)
		case StoryEnterNextActivity:
			// We can always allow this to happen since the logic is triggered in RoomEnterActivity with index 7.
			if u.story.level == len(t.Stories)-1 {
				req.Add(TowerCompleteActivity{dude: u.dude})
			} else {
				nextStory := t.Stories[u.story.level+1]

				if u.dude.room != nil {
					u.dude.room.RemoveDudeFromCenter(u.dude)
					u.dude.room.RemoveDude(u.dude)
					u.dude.room = nil
				}
				u.story.RemoveDude(u.dude)
				nextStory.AddDude(u.dude)
			}
		case RoomCombatActivity:
			if act := u.dude.Trigger(EventCombatRoom{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
		case RoomEnterActivity:
			if u.dude.room != nil {
				if act := u.dude.Trigger(EventLeaveRoom{room: u.dude.room, dude: u.dude}); act != nil {
					req.Add(act)
				}
				u.dude.room = nil
			}
			u.dude.room = u.room
			if act := u.dude.Trigger(EventEnterRoom{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
			// If it's the last room, then move upwards and go poof (unless we're coming from stairs or are entering the tower for the first time).
			if u.dude.activity != StairsFromDown && u.dude.activity != FirstEntering && u.room.index == 7 {
				u.dude.SetActivity(StairsToUp)
			}
		case RoomCenterActivity:
			if act := u.dude.Trigger(EventCenterRoom{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
		case RoomWaitActivity:
			if act := u.dude.Trigger(EventWaitRoom{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
			if u.dude != nil {
				u.dude.SetActivity(Idle)
			}
			// Set wait activity so we can notify the player of their dudes' impending doom
			req.Add(u)
		case RoomStartBossActivity:
			if act := u.dude.Trigger(EventStartBoss{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
			if u.dude != nil {
				u.dude.SetActivity(FightBoss)
			}
			// Set boss start up to tower so we can show health bar.
			req.Add(u)
		case RoomEndBossActivity:
			if act := u.dude.Trigger(EventEndBoss{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
			if u.dude != nil {
				// Keep it movin bub
				u.dude.SetActivity(Centering)
			}
			// Set boss end up to tower so we can show health bar.
			req.Add(u)
		case RoomEndActivity:
			if act := u.dude.Trigger(EventEndRoom{room: u.room, dude: u.dude}); act != nil {
				req.Add(act)
			}
			// Check if the given room is the last room in the story.
			if u.room.story.rooms[6] == u.room {
				if u.room.story.level == len(t.Stories) {
					u.dude.SetActivity(Idle)
				} else if u.room.story.level < len(t.Stories)-1 {
					nextStory := t.Stories[u.room.story.level+1]
					if !nextStory.open {
						if !t.portalOpen {
							t.portalOpen = true
							u.room.story.AddPortal()
						}
						u.dude.SetActivity(EnterPortal)
					}
				}
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
	t.dudes = append(t.dudes, d)
}

func (t *Tower) RemoveDude(d *Dude) {
	// Remove dude from any rooms.
	if d.room != nil {
		d.room.RemoveDudeFromCenter(d)
		d.room.RemoveDude(d)
		d.room = nil
	}
	// Remove from story.
	if d.story != nil {
		d.story.RemoveDude(d)
		d.story = nil
	}
	// Finally remove from our own slice.
	for i, dude := range t.dudes {
		if dude == d {
			t.dudes = append(t.dudes[:i], t.dudes[i+1:]...)
			return
		}
	}
}

func (t *Tower) AddDudes(dudes ...*Dude) {
	for _, d := range dudes {
		t.AddDude(d)
	}
}

func (t *Tower) HasAliveDudes() bool {
	b := false
	for _, d := range t.dudes {
		if !d.IsDead() {
			fmt.Println("Dude is alive: ", d.Name())
			fmt.Println("Stats: ", d.stats)
			b = true
			break
		} else {
			fmt.Println("Dude is dead: ", d.Name())
		}
	}
	return b
}

func (t *Tower) ClearBodies() {
	// Remove dude from any rooms.
	for _, d := range t.dudes {
		if d.room != nil {
			d.room.RemoveDudeFromCenter(d)
			d.room.RemoveDude(d)
		}
		// Remove from story.
		if d.story != nil {
			d.story.RemoveDude(d)
		}
	}
	t.dudes = nil
}

func (t *Tower) ClearTexts() {
	for _, s := range t.Stories {
		s.texts = nil
	}
}

func (t *Tower) Reset() {
	for _, s := range t.Stories {
		s.Reset()
	}
}
