package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type GameState interface {
	Begin(g *Game)
	End(g *Game)
	Update(g *Game) GameState
	Draw(g *Game, screen *ebiten.Image)
}

type GameStatePreBuild struct {
	newDudes []*Dude
}

func (s *GameStatePreBuild) Begin(g *Game) {
	professions := []ProfessionKind{Knight, Vagabond, Ranger, Cleric}
	dudeLimit := len(professions)
	for i := 0; i < dudeLimit; i++ {
		pk := professions[i%len(professions)]
		dude := NewDude(pk, 1)
		dude.stats.agility += i * 5
		s.newDudes = append(s.newDudes, dude)
	}
	// Add some more randomized dudes.
	for i := 0; i < 3; i++ {
		pk := professions[rand.Intn(len(professions))]
		dude := NewDude(pk, 1)
		dude.stats.agility += i * 5
		s.newDudes = append(s.newDudes, dude)
	}
	g.camera.SetMode(render.CameraModeTower)
}
func (s *GameStatePreBuild) End(g *Game) {
	g.dudes = append(g.dudes, s.newDudes...)
	g.tower.AddDudes(s.newDudes...)
}
func (s *GameStatePreBuild) Update(g *Game) GameState {
	return &GameStateBuild{}
}
func (s *GameStatePreBuild) Draw(g *Game, screen *ebiten.Image) {
}

type GameStateBuild struct {
	availableRooms []RoomDef
}

func (s *GameStateBuild) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeStack)
}
func (s *GameStateBuild) End(g *Game) {
}
func (s *GameStateBuild) Update(g *Game) GameState {
	return &GameStatePlay{}
}
func (s *GameStateBuild) Draw(g *Game, screen *ebiten.Image) {
}

type GameStatePlay struct {
	paused       bool
	pauseWobbler float64
}

func (s *GameStatePlay) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeStack)
	// TODO: Set up dude state to spawn outside first story?
	g.ui.dudePanel.SyncDudes(g.dudes)
}
func (s *GameStatePlay) End(g *Game) {
	// TODO: Create a portal at highest story's last room and issue dudes to walk into it?
}
func (s *GameStatePlay) Update(g *Game) GameState {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.paused = !s.paused
	}
	if s.paused {
		s.pauseWobbler += 0.05
	}

	// Update the game!
	if !s.paused {
		g.tower.Update()
	}
	// TODO: Periodically sync dudes with panel??? Or mark dudes as dirty if armor changes then refresh?

	g.ui.Update(&g.uiOptions)

	return nil
}
func (s *GameStatePlay) Draw(g *Game, screen *ebiten.Image) {
	options := render.Options{Screen: screen, Overlay: g.overlay, Camera: &g.camera}

	// Transform our options via the camera.
	g.camera.Transform(&options)

	// Draw that tower -> story -> room -> ???
	g.tower.Draw(&options)

	// Render stuff
	for _, r := range g.renderables {
		r.Draw(&options)
	}

	// Draw UI
	options.DrawImageOptions.GeoM.Reset()
	options.DrawImageOptions.ColorScale.Reset()
	g.ui.Draw(&options)

	// Draw pause
	if s.paused {
		geom := ebiten.GeoM{}
		geom.Scale(4, 4)
		opts := &render.TextOptions{
			Screen: screen,
			Font:   assets.DisplayFont,
			Color:  color.Black,
			GeoM:   geom,
		}

		w, h := text.Measure("PAUSED", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(-w/2, -h/2)
		opts.GeoM.Rotate(math.Sin(s.pauseWobbler) * 0.05)
		opts.GeoM.Translate(w/2, h/2)
		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

		opts.GeoM.Translate(-10, -10)
		opts.Color = color.NRGBA{0, 0, 0, 200}
		render.DrawText(opts, "PAUSED")
		opts.GeoM.Translate(20, 20)
		render.DrawText(opts, "PAUSED")
		opts.Color = color.NRGBA{200, 200, 200, 200}
		opts.GeoM.Translate(-10, -10)
		render.DrawText(opts, "PAUSED")
	}
}

type GameStateLose struct {
}

type GameStateWin struct {
}
