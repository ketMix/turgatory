package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type Game struct {
	renderables []render.Renderable
	camera      render.Camera
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Pitch += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Pitch -= 0.01
	}
	if g.camera.Pitch < 0 {
		g.camera.Pitch = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.SetRotation(g.camera.Rotation() - 0.01)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.SetRotation(g.camera.Rotation() + 0.01)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.camera.Zoom += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.camera.Zoom -= 0.01
	}

	for _, r := range g.renderables {
		r.Update()
		r.SetRotation(r.Rotation() + 0.01)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	options := render.Options{Screen: screen, Camera: &g.camera}

	g.camera.Transform(&options)

	for _, r := range g.renderables {
		r.Draw(&options)
	}
}

func (g *Game) Init() {
	// creates 4 slices of pie with positions to evenly spread them
	for i := 0; i < 4; i++ {
		stack, err := render.NewStack("floors/base", "", "")
		if err != nil {
			panic(err)
		}
		if i%2 == 0 {
			stack.SetStack("rocky")
		}

		rotationAngle := math.Pi / 2 * float64(i)
		stack.SetRotation(rotationAngle)

		stack.SetRotationDistance(0)

		// Append the stack to the renderables
		g.renderables = append(g.renderables, stack)
	}
	g.camera = *render.NewCamera(0, 0)
}

func New() *Game {
	return &Game{}
}
