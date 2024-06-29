package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type GameStatePre struct {
	title    *UIText
	title2   *UIText
	short    ButtonPanel
	medium   ButtonPanel
	long     ButtonPanel
	infinite ButtonPanel
	info     *UIText

	//
	gameLength   int
	titleVolTick int
}

const AUDIO_FADE_IN_TICK = 60

func (s *GameStatePre) Begin(g *Game) {
	// TODO: Hide UI crap.
	g.ui.Hide()

	// Always turn off autoplay (maybe???)
	g.SetAutoplay(false)

	s.title = NewUIText("Time is your purgatory...", assets.DisplayFont, assets.ColorHeading)
	s.title2 = NewUIText("choose wisely.", assets.DisplayFont, assets.ColorHeading)

	s.short = MakeButtonPanel(assets.DisplayFont, PanelStyleButton)
	s.short.onClick = func() {
		s.gameLength = 1
	}
	s.short.onHover = func() {
		s.info.SetText("Just a lil trip.")
	}
	s.short.text.SetText("short")
	s.medium = MakeButtonPanel(assets.DisplayFont, PanelStyleButton)
	s.medium.onClick = func() {
		s.gameLength = 2
	}
	s.medium.onHover = func() {
		s.info.SetText("The ideal.")
	}
	s.medium.text.SetText("medium")
	s.long = MakeButtonPanel(assets.DisplayFont, PanelStyleButton)
	s.long.onClick = func() {
		s.gameLength = 3
	}
	s.long.onHover = func() {
		s.info.SetText("The long haul.")
	}
	s.long.text.SetText("long")
	s.infinite = MakeButtonPanel(assets.DisplayFont, PanelStyleButton)
	s.infinite.onClick = func() {
		s.gameLength = -1
	}
	s.infinite.onHover = func() {
		s.info.SetText("How long can you last?")
	}
	s.infinite.text.SetText("endless")

	s.info = NewUIText("beep boop", assets.BodyFont, assets.ColorStory)

	// Init inventory
	g.equipment = make([]*Equipment, 0)

	// Init audio tick
	g.audioController.PlayRoomTracks()
	g.audioController.MuteAll()
	g.playedTitleSong = false
	s.titleVolTick = 0
}

func (s *GameStatePre) End(g *Game) {
	g.ui.Reveal()
}
func (s *GameStatePre) Update(g *Game) GameState {
	if s.titleVolTick < AUDIO_FADE_IN_TICK {
		s.titleVolTick++
		g.audioController.SetTitleTrackVolPercent(float64(s.titleVolTick) / AUDIO_FADE_IN_TICK)
	}

	w, h := float64(g.uiOptions.Width), float64(g.uiOptions.Height)

	s.short.Layout(nil, &g.uiOptions)
	s.medium.Layout(nil, &g.uiOptions)
	s.long.Layout(nil, &g.uiOptions)
	s.infinite.Layout(nil, &g.uiOptions)

	panelsWidth := 0.0
	panelsWidth += s.short.Width()
	panelsWidth += s.medium.Width()
	panelsWidth += s.long.Width()
	panelsWidth += s.infinite.Width()

	panelsX := (w - panelsWidth) / 2
	panelsY := h / 2

	y := 0.0

	s.title.SetPosition(32, 32)
	s.title.Layout(nil, &g.uiOptions)
	y = panelsY - s.title.Height()*2 - 4*g.uiOptions.Scale
	s.title.SetPosition(w/2-s.title.Width()/2, y)
	y += s.title.Height()/1.5 + 4*g.uiOptions.Scale
	s.title2.Layout(nil, &g.uiOptions)
	s.title2.SetPosition(w/2-s.title2.Width()/2, y)
	y += s.title2.Height() + 4*g.uiOptions.Scale

	s.short.SetPosition(panelsX, y)
	panelsX += s.short.Width()
	s.medium.SetPosition(panelsX, y)
	panelsX += s.medium.Width()
	s.long.SetPosition(panelsX, y)
	panelsX += s.long.Width()
	s.infinite.SetPosition(panelsX, y)

	y += s.short.Height() + 4*g.uiOptions.Scale

	s.info.Layout(nil, &g.uiOptions)
	s.info.SetPosition(w/2-s.info.Width()/2, y)

	mx, my := IntToFloat2(ebiten.CursorPosition())
	if s.short.Check(mx, my, UICheckHover) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.short.Check(mx, my, UICheckClick)
		}
	} else if s.medium.Check(mx, my, UICheckHover) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.medium.Check(mx, my, UICheckClick)
		}
	} else if s.long.Check(mx, my, UICheckHover) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.long.Check(mx, my, UICheckClick)
		}
	} else if s.infinite.Check(mx, my, UICheckHover) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.infinite.Check(mx, my, UICheckClick)
		}
	} else {
		s.info.SetText("")
	}

	if s.gameLength != 0 {
		return &GameStateStart{
			length: s.gameLength,
		}
	}
	return nil
}
func (s *GameStatePre) Draw(g *Game, screen *ebiten.Image) {
	opts := &render.Options{
		Screen: screen,
	}
	s.title.Draw(opts)
	s.title2.Draw(opts)
	s.info.Draw(opts)
	s.short.Draw(opts)
	s.medium.Draw(opts)
	s.long.Draw(opts)
	s.infinite.Draw(opts)
}
