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

// AddStory does as it says.
func (t *Tower) AddStory(s *Story) {
	t.Stories = append(t.Stories, s)
}
