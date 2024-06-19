package game

import (
	"math"

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
	// creates 4 slices of pie with positions to evenly spread them
	for i := 0; i < 4; i++ {
		stack, err := render.NewStack("walls/pie", "", "")
		if err != nil {
			panic(err)
		}

		rotationAngle := math.Pi / 2 * float64(i)
		stack.SetRotation(rotationAngle)

		stack.SetRotationDistance(0)

		// Append the stack to the renderables
		g.renderables = append(g.renderables, stack)
	}
}

func New() *Game {
	return &Game{}
}
