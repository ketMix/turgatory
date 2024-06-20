package game

import "github.com/kettek/ebijam24/internal/render"

type Level struct {
	towers []*Tower // There should probably only be 1 tower, but we'll let it slide for now...
}

func NewLevel() *Level {
	return &Level{}
}

// Updates level stuff
func (l *Level) Update() {
	for _, t := range l.towers {
		t.Update()
	}
}

// Draw draws the level.
func (l *Level) Draw(o *render.Options) {
	for _, t := range l.towers {
		t.Draw(o)
	}
}

// AddTower adds a forbidden tower.
func (l *Level) AddTower(t *Tower) {
	l.towers = append(l.towers, t)
}
