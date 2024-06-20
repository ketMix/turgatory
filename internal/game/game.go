package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type Game struct {
	renderables []render.Renderable
	camera      render.Camera
	level       *Level
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	// Move this stuff elsewhere, probs.
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

	// Update the level, yo.
	g.level.Update()

	// Update other stuff
	for _, r := range g.renderables {
		r.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	options := render.Options{Screen: screen, Camera: &g.camera}

	// Transform our options via the camera.
	g.camera.Transform(&options)

	// Draw that level -> tower -> story -> room -> ???
	g.level.Draw(&options)

	// Render stuff
	for _, r := range g.renderables {
		r.Draw(&options)
	}
}

func (g *Game) Init() {
	lvl := NewLevel()
	tower := NewTower()
	tower.AddStory(NewStory(8))
	lvl.AddTower(tower)

	g.level = lvl

	g.camera = *render.NewCamera(0, 0)
}

func New() *Game {
	return &Game{}
}
