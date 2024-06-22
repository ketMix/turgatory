package game

import (
	"fmt"

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
	t.stairs.Draw(o) // We draw the stairs first and allow the stories to be drawn overtop.
}

// AddStory does as it says.
func (t *Tower) AddStory(s *Story) {
	t.Stories = append(t.Stories, s)
}
