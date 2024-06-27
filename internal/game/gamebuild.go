package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type GameStateBuild struct {
	availableRooms []*RoomDef
	wobbler        float64
	titleTimer     int
}

func (s *GameStateBuild) Begin(g *Game) {
	g.camera.SetMode(render.CameraModeTower)

	// On build phase, full heal all dudes and restore uses
	for _, d := range g.dudes {
		d.FullHeal()
		d.RestoreUses()
	}

	var nextStory *Story
	for _, s := range g.tower.Stories {
		if !s.open {
			nextStory = s
			break
		}
	}
	// This shouldn't happen.
	if nextStory == nil {
		panic("No next story found!")
	}
	// Generate our new rooms.
	required := GetRequiredRooms(nextStory.level, 2)
	optional := GetOptionalRooms(nextStory.level, 9) // 6 is minimum, but let's given 3 more for fun.
	s.availableRooms = append(s.availableRooms, required...)
	s.availableRooms = append(s.availableRooms, optional...)
	// Update room panel.
	g.ui.roomPanel.SetRoomDefs(s.availableRooms)
	// Add onClick handler.
	g.ui.roomPanel.onItemClick = func(which int) {
		g.ui.roomInfoPanel.hidden = false
		selected := s.availableRooms[which]
		g.ui.roomInfoPanel.title.SetText(fmt.Sprintf("%s %s room", selected.size.String(), selected.kind.String()))
		g.ui.roomInfoPanel.description.SetText(selected.GetDescription())
		g.ui.roomInfoPanel.cost.SetText(fmt.Sprintf("Cost: %d", selected.GetCost(nextStory.level)))
	}

	// Update info.
	g.UpdateInfo()
}
func (s *GameStateBuild) End(g *Game) {
	g.ui.roomPanel.onItemClick = nil
	g.ui.roomPanel.SetRoomDefs(nil)
}
func (s *GameStateBuild) Update(g *Game) GameState {
	s.wobbler += 0.05
	s.titleTimer++
	/*if s.titleTimer > 120 {
		return &GameStatePlay{}
	}*/
	return nil
}
func (s *GameStateBuild) Draw(g *Game, screen *ebiten.Image) {
	if s.titleTimer < 240 {
		opts := render.TextOptions{
			Screen: screen,
			Font:   assets.DisplayFont,
			Color:  assets.ColorState,
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

func (s *GameStateBuild) BuyDude(g *Game) {
	// COST?
	cost := 100
	if g.gold < cost {
		AddMessage(MessageError, "Not enough gold to buy a dude.")
		return
	}
	g.gold -= cost
	level := len(g.tower.Stories) / 2
	if level < 1 {
		level = 1
	}

	// Random profession ??
	profession := RandomProfessionKind()
	dude := NewDude(profession, level)
	g.dudes = append(g.dudes, dude)
	g.UpdateInfo()
}

func (s *GameStateBuild) BuyEquipment(g *Game) {
	// COST?
	cost := 50
	if g.gold < cost {
		AddMessage(MessageError, "Not enough gold to buy equipment.")
		return
	}
	g.gold -= cost

	level := len(g.tower.Stories) / 2
	e := GetRandomEquipment(level)
	g.equipment = append(g.equipment, e)
}

func (s *GameStateBuild) SellEquipment(g *Game, e *Equipment) {
	if e == nil {
		return
	}
	g.gold += int(e.GoldValue())
	for i, eq := range g.equipment {
		if eq == e {
			g.equipment = append(g.equipment[:i], g.equipment[i+1:]...)
			break
		}
	}

	// Trigger on sell event
	e.Activate(EventSell{equipment: e})

}
