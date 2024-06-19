package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type Game struct {
	renderables []render.Renderable
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	for _, r := range g.renderables {
		r.Update()
		r.SetRotation(r.Rotation() + 0.01)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, r := range g.renderables {
		r.Draw(render.Options{Screen: screen})
	}
}

func (g *Game) Init() {
	stack, err := render.NewStack("walls/base", "", "")
	if err != nil {
		panic(err)
	}
	stack.SetPosition(100, 100)
	g.renderables = append(g.renderables, stack)
}

func New() *Game {
	return &Game{}
}
