package main

import (
	"github.com/kettek/ebijam24/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()
	g.Init()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("ebijam24")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
