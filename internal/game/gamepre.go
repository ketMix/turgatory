package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameStatePre struct {
}

func (s *GameStatePre) Begin(g *Game) {
}
func (s *GameStatePre) End(g *Game) {
}
func (s *GameStatePre) Update(g *Game) GameState {
	return &GameStateStart{}
}
func (s *GameStatePre) Draw(g *Game, screen *ebiten.Image) {
}
