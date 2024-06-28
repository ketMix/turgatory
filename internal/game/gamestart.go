package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/internal/render"
)

type GameStateStart struct {
	newDudes []*Dude
}

func (s *GameStateStart) Begin(g *Game) {
	// Give the player a reasonable amount of GOLD
	g.gold = 700

	professions := []ProfessionKind{Knight, Vagabond, Ranger, Cleric}
	dudeLimit := len(professions)
	for i := 0; i < dudeLimit; i++ {
		pk := professions[i%len(professions)]
		dude := NewDude(pk, 1)
		s.newDudes = append(s.newDudes, dude)
	}
	// Add some more randomized dudes.
	for i := 0; i < 3; i++ {
		pk := professions[rand.Intn(len(professions))]
		dude := NewDude(pk, 1)
		s.newDudes = append(s.newDudes, dude)
	}

	g.camera.SetMode(render.CameraModeTower)

	// Create a new tower, yo.
	tower := NewTower()

	firstStory := NewStory()
	firstStory.Open()
	tower.AddStory(firstStory)
	tower.AddStory(NewStory())
	tower.AddStory(NewStory())
	tower.AddStory(NewStory())
	tower.AddStory(NewStory())

	g.tower = tower
}
func (s *GameStateStart) End(g *Game) {
	g.dudes = append(g.dudes, s.newDudes...)
}
func (s *GameStateStart) Update(g *Game) GameState {
	//return &GameStateWin{}
	return &GameStateBuild{}
}
func (s *GameStateStart) Draw(g *Game, screen *ebiten.Image) {
}
