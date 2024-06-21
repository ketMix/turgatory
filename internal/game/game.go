package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type Game struct {
	renderables      []render.Renderable
	camera           render.Camera
	mouseX, mouseY   int
	cursorX, cursorY float64
	level            *Level
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Always set the camera's origin to be half the size of the screen.
	g.camera.SetOrigin(float64(outsideWidth/2), float64(outsideHeight/2))

	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	// Transform mouse coordinates by camera.
	g.cursorX, g.cursorY = g.camera.ScreenToWorld(float64(g.mouseX), float64(g.mouseY))

	// Move this stuff elsewhere, probs.
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Pitch += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Pitch -= 0.01
	}
	if g.camera.Pitch < 0 {
		g.camera.Pitch = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.SetRotation(g.camera.Rotation() - 0.01)
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.SetRotation(g.camera.Rotation() + 0.01)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		x, y := g.camera.Position()
		g.camera.SetPosition(x-1, y)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		x, y := g.camera.Position()
		g.camera.SetPosition(x+1, y)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		x, y := g.camera.Position()
		g.camera.SetPosition(x, y-1)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		x, y := g.camera.Position()
		g.camera.SetPosition(x, y+1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.camera.Zoom += g.camera.Zoom * 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.camera.Zoom -= g.camera.Zoom * 0.01
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

	// Debug render
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%0.1fx%0.1f", g.cursorX, g.cursorY), g.mouseX, g.mouseY-16)
	ox, oy := g.camera.WorldToScreen(0, 0)
	ebitenutil.DebugPrintAt(screen, "0x0", int(ox)-8, int(oy)-8)
}

func (g *Game) Init() {
	lvl := NewLevel()
	tower := NewTower()

	firstStory := NewStory()
	tower.AddStory(firstStory)

	tower.AddStory(NewStoryWithSize(16))
	lvl.AddTower(tower)

	g.level = lvl

	g.camera = *render.NewCamera(0, 0)
}

func New() *Game {
	return &Game{}
}
