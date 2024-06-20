package game

import "github.com/kettek/ebijam24/internal/render"

// Tower is our glorious tower :o
type Tower struct {
	render.Positionable // I guess it's okay to re-use this in such a fashion.
	Stories             []*Story
}

// NewTower creates a new tower.
func NewTower() *Tower {
	return &Tower{}
}

// Update the tower.
func (t *Tower) Update() {
	// TODO: Should this only update "active" stories?
	for _, s := range t.Stories {
		s.Update()
	}
}

// Draw our glorious tower.
func (t *Tower) Draw(o *render.Options) {
	for _, s := range t.Stories {
		s.Draw(o)
		o.DrawImageOptions.GeoM.Translate(0, -StoryHeight*o.Camera.Zoom) // Transform our rendering, ofc
	}
}

// AddStory does as it says.
func (t *Tower) AddStory(s *Story) {
	t.Stories = append(t.Stories, s)
}
