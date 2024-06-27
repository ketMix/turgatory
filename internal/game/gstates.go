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
	// Always remove door from last story(?)
	tower.Stories[len(tower.Stories)-1].RemoveDoor()

	g.tower = tower
}
func (s *GameStatePreBuild) End(g *Game) {
	g.dudes = append(g.dudes, s.newDudes...)
}
func (s *GameStatePreBuild) Update(g *Game) GameState {
	//return &GameStateWin{}
	return &GameStateBuild{}
}
func (s *GameStatePreBuild) Draw(g *Game, screen *ebiten.Image) {
}

type GameStateBuild struct {
	availableRooms []RoomDef
	wobbler        float64
	titleTimer     int
}

func (s *GameStateBuild) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)

	// On build phase, full heal all dudes and restore uses
	for _, d := range g.dudes {
		d.FullHeal()
		d.FullUseRestore()
	}
}
func (s *GameStateBuild) End(g *Game) {
}
func (s *GameStateBuild) Update(g *Game) GameState {
	s.wobbler += 0.05
	s.titleTimer++
	if s.titleTimer > 120 {
		return &GameStatePlay{}
	}
	return nil
}
func (s *GameStateBuild) Draw(g *Game, screen *ebiten.Image) {
	if s.titleTimer < 240 {
		opts := render.TextOptions{
			Screen: screen,
			Font:   assets.DisplayFont,
			Color:  color.NRGBA{184, 152, 93, 200},
		}
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("BUILD", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(-w/2, -h/2)
		opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.05)
		opts.GeoM.Translate(w/2, h/2)
		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

		render.DrawText(&opts, "BUILD")
	}
}

type GameStatePlay struct {
	titleTimer     int
	wobbler        float64
	updateTicker   int
	returningDudes []*Dude
}

func (s *GameStatePlay) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeStack)

	// Make sure our dudes are visible.
	for _, d := range g.dudes {
		d.stack.Transparency = 0
		d.shadow.Transparency = 0
	}

	// Add our dudes to the tower!
	g.tower.AddDudes(g.dudes...)
	// TODO: Set up dude state to spawn outside first story?
	g.ui.dudePanel.SyncDudes(g.dudes)
}
func (s *GameStatePlay) End(g *Game) {
	// Replace our current dudes!
	g.dudes = nil
	g.dudes = append(g.dudes, s.returningDudes...)
	for i, s := range g.tower.Stories {
		if !s.open {
			s.Open()
			if i-1 >= 0 {
				// Remove door, ofc
				g.tower.Stories[i-1].RemoveDoor()
				// Remove old portal
				g.tower.Stories[i-1].RemovePortal()
				g.tower.portalOpen = false
			}
			break
		}
	}

	// Clear out any floating text.
	g.tower.ClearTexts()

	// Collect gold
	g.CollectGold()
	// Collect inventory
	g.CollectInventory()
}
func (s *GameStatePlay) Update(g *Game) GameState {
	s.titleTimer++

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.TogglePause()
	}
	s.wobbler += 0.05

	// Update the game!
	if !g.paused {
		s.updateTicker++
		if s.updateTicker > g.speed {
			var req ActivityRequests
			g.tower.Update(&req)
			for _, u := range req {
				switch u := u.(type) {
				case DudeDeadActivity:
					if !g.tower.HasAliveDudes() {
						return &GameStateLose{}
					}
				case TowerCompleteActivity:
					return &GameStateWin{}
				case TowerLeaveActivity:
					s.returningDudes = append(s.returningDudes, u.dude)
					if !g.tower.HasAliveDudes() {
						// No more alive dudes! Switch game state, yo.
						g.tower.ClearBodies()
						return &GameStateBuild{}
					}
				}
			}

			s.updateTicker = 0
		}
	}
	// TODO: Periodically sync dudes with panel??? Or mark dudes as dirty if armor changes then refresh?

	g.ui.Update(&g.uiOptions)

	return nil
}
func (s *GameStatePlay) Draw(g *Game, screen *ebiten.Image) {

	if s.titleTimer < 240 {
		opts := render.TextOptions{
			Screen: screen,
			Font:   assets.DisplayFont,
			Color:  color.NRGBA{184, 152, 93, 200},
		}
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("ADVENTURE", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(-w/2, -h/2)
		opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.05)
		opts.GeoM.Translate(w/2, h/2)
		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

		render.DrawText(&opts, "ADVENTURE")
	}

	// Draw pause
	if g.paused {
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
		opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.05)
		opts.GeoM.Translate(w/2, h/2)
		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

		opts.GeoM.Translate(-10, -10)
		opts.Color = color.NRGBA{10, 0, 0, 200}
		render.DrawText(opts, "PAUSED")
		opts.GeoM.Translate(20, 20)
		render.DrawText(opts, "PAUSED")
		opts.Color = color.NRGBA{184, 152, 93, 200}
		opts.GeoM.Translate(-10, -10)
		render.DrawText(opts, "PAUSED")
	}
}

type GameStateLose struct {
	wobbler float64
}

func (s *GameStateLose) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)
	g.audioController.PlaySfx("loss", 0.5, 0.0)
}
func (s *GameStateLose) End(g *Game) {
}
func (s *GameStateLose) Update(g *Game) GameState {
	s.wobbler += 0.05
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return &GameStatePreBuild{}
	}
	return nil
}
func (s *GameStateLose) Draw(g *Game, screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	opts := &render.TextOptions{
		Screen: screen,
		Font:   assets.DisplayFont,
		Color:  color.Black,
		GeoM:   geom,
	}

	w, h := text.Measure("GAME OVER", opts.Font.Face, opts.Font.LineHeight)
	w *= 4
	h *= 4

	opts.GeoM.Translate(-w/2, -h/2)
	opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.1)
	opts.GeoM.Translate(w/2, h/2)
	opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2+h)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*20)

	opts.GeoM.Translate(-10, -10)
	opts.Color = color.NRGBA{50, 0, 0, 200}
	render.DrawText(opts, "GAME OVER")
	opts.GeoM.Translate(20, 20)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*2)
	render.DrawText(opts, "GAME OVER")
	opts.Color = color.NRGBA{255, 52, 33, 200}
	opts.GeoM.Translate(-10, -10)
	opts.GeoM.Translate(0, math.Cos(s.wobbler)*2)
	render.DrawText(opts, "GAME OVER")

	opts.Font = assets.BodyFont
	y := 0.0
	{
		opts.GeoM.Reset()
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("all ur dudes r ded :((", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/2)-h/2)
		y = float64(screen.Bounds().Dy()/2) - h/2
		render.DrawText(opts, "all ur dudes r ded :((")
	}

	{
		opts.GeoM.Reset()
		opts.GeoM.Scale(4, 4)

		w, h := text.Measure("Press SPACE to try again!", opts.Font.Face, opts.Font.LineHeight)
		w *= 4
		h *= 4

		opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, y+h)
		render.DrawText(opts, "Press SPACE to try again!")
	}
}

type GameStateWin struct {
	wobbler float64
}

func (s *GameStateWin) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)
	g.audioController.PlaySfx("win", 0.5, 0.0)
}
func (s *GameStateWin) End(g *Game) {
}

func (s *GameStateWin) Update(g *Game) GameState {
	s.wobbler += 0.05
	g.camera.SetRotation(g.camera.Rotation() + 0.005)
	return nil
}

func (s *GameStateWin) Draw(g *Game, screen *ebiten.Image) {
	s.DrawRainbow(screen, "U   WINNE!1!")
}

func (s *GameStateWin) DrawRainbow(screen *ebiten.Image, t string) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	opts := &render.TextOptions{
		Screen: screen,
		Font:   assets.DisplayFont,
		Color:  color.Black,
		GeoM:   geom,
	}

	w, h := text.Measure(t, opts.Font.Face, opts.Font.LineHeight)
	w *= 4
	h *= 4

	opts.GeoM.Translate(-w/2, -h/2)
	opts.GeoM.Rotate(math.Sin(s.wobbler) * 0.05)
	opts.GeoM.Translate(w/2, h/2)
	opts.GeoM.Translate(float64(screen.Bounds().Dx()/2)-w/2, float64(screen.Bounds().Dy()/4)-h/2)

	opts.GeoM.Translate(0, h*2) // Uh... this *2 is odd.

	// RAINBOW
	var colors = []color.NRGBA{
		{255, 0, 0, 255},
		{255, 127, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
		{75, 0, 130, 255},
		{238, 130, 238, 255},
	}
	ci := 0

	// Draw "shadow"
	oldGeoM := ebiten.GeoM{}
	oldGeoM.Concat(opts.GeoM)
	for _, r := range t {
		tx, _ := text.Measure(string(r), opts.Font.Face, opts.Font.LineHeight)
		tx *= 4

		opts.GeoM.Translate(-10, -10)
		opts.Color = color.NRGBA{10, 0, 0, 200}
		render.DrawText(opts, string(r))
		opts.GeoM.Translate(20, 20)
		render.DrawText(opts, string(r))
		opts.GeoM.Translate(-10, -10)

		opts.GeoM.Translate(tx, 0)
		opts.GeoM.Rotate(math.Cos(s.wobbler) * 0.003)
	}

	// Draw main text
	opts.GeoM = oldGeoM
	for _, r := range t {
		tx, _ := text.Measure(string(r), opts.Font.Face, opts.Font.LineHeight)
		tx *= 4

		opts.Color = colors[ci%len(colors)]
		render.DrawText(opts, string(r))

		opts.GeoM.Translate(tx, 0)
		opts.GeoM.Rotate(math.Cos(s.wobbler) * 0.003)

		ci++
	}
}
