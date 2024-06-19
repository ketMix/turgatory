package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/assets"
	_ "github.com/kettek/ebijam24/assets"
)

type Game struct {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Init() {
	assets.LoadStack(("walls/base"))
}

func New() *Game {
	return &Game{}
}
