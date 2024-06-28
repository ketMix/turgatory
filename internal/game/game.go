package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type Game struct {
	ui                    *UI
	dudes                 []*Dude
	renderables           []render.Renderable
	camera                render.Camera
	mouseX, mouseY        int
	cursorX, cursorY      float64
	tower                 *Tower
	lastWidth, lastHeight int
	uiOptions             UIOptions
	state                 GameState
	audioController       *AudioController
	overlay               *ebiten.Image
	selectedDude          *Dude
	hoveredDude           *Dude
	paused                bool
	speed                 int
	gold                  int
	equipment             []*Equipment
}

type GameState interface {
	Begin(g *Game)
	End(g *Game)
	Update(g *Game) GameState
	Draw(g *Game, screen *ebiten.Image)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != g.lastWidth && outsideHeight != g.lastHeight {
		// Always set the camera's origin to be half the size of the screen.
		g.camera.SetOrigin(float64(outsideWidth/2), float64(outsideHeight/2))
		g.lastWidth, g.lastHeight = outsideWidth, outsideHeight
		g.uiOptions.Width, g.uiOptions.Height = outsideWidth, outsideHeight
		if g.overlay != nil {
			g.overlay.Deallocate()
		}
		g.overlay = ebiten.NewImage(outsideWidth, outsideHeight)
		g.ui.Layout(&g.uiOptions)
	}

	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	// Transform mouse coordinates by camera.
	g.cursorX, g.cursorY = g.camera.ScreenToWorld(float64(g.mouseX), float64(g.mouseY))

	g.camera.Update()

	// Move this stuff elsewhere, probs.
	/*if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Pitch += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Pitch -= 0.01
	}
	if g.camera.Pitch < 0 {
		g.camera.Pitch = 0
	}*/

	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.SetRotation(g.camera.Rotation() - 0.01)
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.SetRotation(g.camera.Rotation() + 0.01)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.camera.SetStory(g.camera.Story() + 1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.camera.SetStory(g.camera.Story() - 1)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.camera.SetMode(render.CameraModeTower)
		g.ui.speedPanel.cameraButton.SetImage("tower")
		g.ui.speedPanel.cameraButton.tooltip = "camera: tower"
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.camera.SetMode(render.CameraModeStack)
		g.ui.speedPanel.cameraButton.SetImage("story")
		g.ui.speedPanel.cameraButton.tooltip = "camera: story"
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.camera.SetMode(render.CameraModeSuperZoom)
		g.ui.speedPanel.cameraButton.SetImage("room")
		g.ui.speedPanel.cameraButton.tooltip = "camera: room"
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.camera.ZoomIn()
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.camera.ZoomOut()
	}

	// FIXME: For some reason Layout doesn't position the UI properly...
	g.ui.Layout(&g.uiOptions)

	g.ui.Update(&g.uiOptions)

	if nextState := g.state.Update(g); nextState != nil {
		g.state.End(g)
		g.state = nextState
		g.state.Begin(g)
	}

	// If we have a tower with rooms, synchronize the music.
	if g.tower != nil {
		for _, story := range g.tower.Stories {
			if story.open {
				// Build map of roomkind to max vol
				roomPanVol := make(map[RoomKind]PanVol)
				for i, room := range story.rooms {
					pan, vol := room.getPanVol(g.camera.Rotation(), story.GetRoomCenterRad(i), 1.0) // Replace 1.0 with a calculation based on focused story index vs. current
					if roomPanVol[room.kind].Vol < vol {
						roomPanVol[room.kind] = PanVol{Pan: pan, Vol: vol}
					}
				}
				g.audioController.SetStoryPanVol(roomPanVol)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear overlay.
	g.overlay.Clear()

	g.DrawTower(screen)

	g.state.Draw(g, screen)

	// Draw overlay.
	screen.DrawImage(g.overlay, nil)

	options := render.Options{Screen: screen, Overlay: g.overlay, Camera: &g.camera}

	// Draw UI
	options.DrawImageOptions.GeoM.Reset()
	options.DrawImageOptions.ColorScale.Reset()
	g.ui.Draw(&options)

	// Print fps
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%0.2f", ebiten.CurrentFPS()), 0, 0)
}

func (g *Game) DrawTower(screen *ebiten.Image) {
	options := render.Options{Screen: screen, Overlay: g.overlay, Camera: &g.camera}
	// Transform our options via the camera.
	g.camera.Transform(&options)

	// Draw that tower -> story -> room -> ???
	g.tower.Draw(&options)

	// Render stuff
	for _, r := range g.renderables {
		r.Draw(&options)
	}
}

func (g *Game) CheckUI() (bool, UICheckKind) {
	mx, my := IntToFloat2(ebiten.CursorPosition())
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.ui.Check(mx, my, UICheckClick) {
			return true, UICheckClick
		}
		g.selectedDude = nil
		//g.ui.dudeInfoPanel.HideDetails()
		g.ui.dudeInfoPanel.SetDude(g.hoveredDude)
		if g.hoveredDude == nil {
			g.ui.dudeInfoPanel.showDetails = false
		} else {
			g.ui.dudeInfoPanel.showDetails = true
		}
		g.ui.equipmentPanel.showDetails = false
		g.ui.roomInfoPanel.hidden = true
		return false, UICheckClick
	} else {
		if g.ui.Check(mx, my, UICheckHover) {
			return true, UICheckHover
		}
		g.hoveredDude = nil
		g.ui.dudeInfoPanel.SetDude(g.selectedDude)
		if g.selectedDude != nil {
			g.ui.dudeInfoPanel.showDetails = true
		} else {
			g.ui.dudeInfoPanel.showDetails = false
		}
		/*if g.selectedDude != nil {
			g.ui.dudeInfoPanel.ShowDetails()
		} else {
			g.ui.dudeInfoPanel.HideDetails()
		}*/
		return false, UICheckHover
	}
	return false, UICheckNone
}

func (g *Game) UpdateInfo() {
	var currentStory *Story
	for i := len(g.tower.Stories) - 1; i >= 0; i-- {
		story := g.tower.Stories[i]
		if story.open {
			currentStory = story
			break
		}
	}
	if currentStory != nil {
		g.ui.gameInfoPanel.storyText.SetText(fmt.Sprintf("Stories: %d/%d", currentStory.level, len(g.tower.Stories)))
	}
	g.ui.gameInfoPanel.goldText.SetText(fmt.Sprintf("Gold: %d", g.gold))
	g.ui.gameInfoPanel.dudeText.SetText(fmt.Sprintf("Dudes: %d", len(g.dudes)))
	// Move this if it's too heavy.
	g.ui.dudePanel.SetDudes(g.dudes)
}

func (g *Game) Init() {
	// Init the equipment
	assets.LoadEquipment()

	g.ui = NewUI()
	g.uiOptions = UIOptions{Scale: 2.0}
	g.ui.speedPanel.pauseButton.onCheck = func(kind UICheckKind) {
		if kind == UICheckClick {
			g.TogglePause()
		}
	}
	g.ui.speedPanel.speedButton.onCheck = func(kind UICheckKind) {
		if kind == UICheckClick {
			g.AdjustSpeed()
		}
	}
	g.ui.speedPanel.cameraButton.onCheck = func(kind UICheckKind) {
		if kind == UICheckClick {
			g.AdjustCamera()
		}
	}
	g.ui.speedPanel.musicButton.onCheck = func(kind UICheckKind) {
		if kind == UICheckClick {
			if g.audioController.tracksPaused {
				g.audioController.PlayRoomTracks()
				g.ui.speedPanel.musicButton.SetImage("music")
				g.ui.speedPanel.musicButton.tooltip = "music on"
			} else {
				g.audioController.PauseRoomTracks()
				g.ui.speedPanel.musicButton.SetImage("music-mute")
				g.ui.speedPanel.musicButton.tooltip = "music off"
			}
		}
	}
	g.ui.speedPanel.soundButton.onCheck = func(kind UICheckKind) {
		if kind == UICheckClick {
			if g.audioController.sfxPaused {
				g.audioController.sfxPaused = false
				g.ui.speedPanel.soundButton.SetImage("sound")
				g.ui.speedPanel.soundButton.tooltip = "sound on"
			} else {
				g.audioController.sfxPaused = true
				g.ui.speedPanel.soundButton.SetImage("sound-mute")
				g.ui.speedPanel.soundButton.tooltip = "sound off"
			}
		}
	}
	g.ui.dudePanel.onItemClick = func(index int) {
		if index < 0 || index >= len(g.dudes) {
			return
		}
		dude := g.dudes[index]
		g.selectedDude = dude
		g.ui.dudeInfoPanel.showDetails = true
		//g.ui.dudeInfoPanel.ShowDetails()
		//g.ui.dudeInfoPanel.SetDude(dude)
	}
	g.ui.dudePanel.onItemHover = func(index int) {
		if index < 0 || index >= len(g.dudes) {
			return
		}
		dude := g.dudes[index]
		g.hoveredDude = dude
		g.ui.dudeInfoPanel.SetDude(dude)
	}

	g.camera = *render.NewCamera(0, 0)
	g.audioController = NewAudioController()
	g.gold = 0
	g.equipment = make([]*Equipment, 0)
	g.state = &GameStateStart{}
	g.state.Begin(g)
}

func (g *Game) TogglePause() {
	g.paused = !g.paused
	if g.paused {
		g.ui.speedPanel.pauseButton.SetImage("pause")
		g.ui.speedPanel.pauseButton.tooltip = "paused"
	} else {
		g.ui.speedPanel.pauseButton.SetImage("play")
		g.ui.speedPanel.pauseButton.tooltip = "playing"
	}
}

func (g *Game) AdjustSpeed() {
	g.speed += 2
	if g.speed > 4 {
		g.speed = 0
	}
	switch g.speed {
	case 0:
		g.ui.speedPanel.speedButton.SetImage("fast")
		g.ui.speedPanel.speedButton.tooltip = "fast"
	case 2:
		g.ui.speedPanel.speedButton.SetImage("medium")
		g.ui.speedPanel.speedButton.tooltip = "medium"
	case 4:
		g.ui.speedPanel.speedButton.SetImage("slow")
		g.ui.speedPanel.speedButton.tooltip = "slow"
	}
}

func (g *Game) AdjustCamera() {
	if g.camera.Mode == render.CameraModeTower {
		g.camera.SetMode(render.CameraModeStack)
		g.ui.speedPanel.cameraButton.SetImage("story")
		g.ui.speedPanel.cameraButton.tooltip = "camera: story"
	} else if g.camera.Mode == render.CameraModeStack {
		g.camera.SetMode(render.CameraModeSuperZoom)
		g.ui.speedPanel.cameraButton.SetImage("room")
		g.ui.speedPanel.cameraButton.tooltip = "camera: room"
	} else {
		g.camera.SetMode(render.CameraModeTower)
		g.ui.speedPanel.cameraButton.SetImage("tower")
		g.ui.speedPanel.cameraButton.tooltip = "camera: tower"
	}
}

func New() *Game {
	return &Game{}
}
