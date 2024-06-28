package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type GameStateBuild struct {
	availableRooms   []*RoomDef
	wobbler          float64
	titleTimer       int
	nextStory        *Story
	focusedRoom      *Room
	highlightedRooms []*Room
	placingRoom      *RoomDef
	placingIndex     int
	readyAttempts    int
	//
	selectedEquipment int
}

func (s *GameStateBuild) Begin(g *Game) {
	// g.audioController.PlayRoomTracks()
	// On build phase, full heal all dudes and restore uses
	for _, d := range g.dudes {
		d.FullHeal()
		d.RestoreUses()
	}

	for _, st := range g.tower.Stories {
		if st.level == len(g.tower.Stories)-1 {
			// If it's the last story, automatically remove the door.
			s.nextStory = st
			st.RemoveDoor()
			break
		} else if st.doorStack != nil {
			s.nextStory = st
			break
		}
	}
	// This shouldn't happen.
	if s.nextStory == nil {
		panic("No next story found!")
	}

	// Make camera reset and focus story.
	g.camera.SetStory(s.nextStory.level)
	g.camera.SetRotation(math.Pi / 8)

	// Just in case
	g.ui.bossPanel.hidden = true
	g.selectedDude = nil

	// Generate our new rooms.
	s.RollRooms(g)
	// Add onClick handler.
	g.ui.roomPanel.onItemClick = func(which int) {
		g.ui.roomInfoPanel.hidden = false
		selected := s.availableRooms[which]
		s.placingRoom = selected
		s.placingIndex = which
		g.ui.roomInfoPanel.title.SetText(fmt.Sprintf("%s %s room", selected.size.String(), selected.kind.String()))
		g.ui.roomInfoPanel.description.SetText(selected.GetDescription())
		g.ui.roomInfoPanel.cost.SetText(fmt.Sprintf("Cost: %d", GetRoomCost(selected.kind, selected.size, s.nextStory.level)))
	}

	g.ui.equipmentPanel.buyButton.onClick = func() {
		s.BuyEquipment(g)
	}
	g.ui.equipmentPanel.buyButton.text.SetText(fmt.Sprintf("Random Loot\n%dgp", s.EquipmentCost()))
	g.ui.equipmentPanel.buyButton.disabled = false

	g.ui.equipmentPanel.onItemClick = func(which int) {
		g.ui.equipmentPanel.list.selected = which
		s.selectedEquipment = which
		g.ui.equipmentPanel.details.SetEquipment(g.equipment[which])
		g.ui.equipmentPanel.showDetails = true
	}
	g.ui.equipmentPanel.details.onSellClick = func() {
		if s.selectedEquipment >= len(g.equipment) {
			return
		}
		eq := g.equipment[s.selectedEquipment]
		s.SellEquipment(g, eq)
		g.ui.equipmentPanel.SetEquipment(g.equipment)
		if s.selectedEquipment >= len(g.equipment) {
			s.selectedEquipment--
		}
		if s.selectedEquipment >= 0 {
			g.ui.equipmentPanel.onItemClick(s.selectedEquipment)
		} else {
			g.ui.equipmentPanel.details.SetEquipment(nil)
		}
	}
	g.ui.equipmentPanel.details.onSwapClick = func() {
		if g.selectedDude == nil {
			g.ui.feedback.Msg(FeedbackBad, "select a dude to swap equipment with!")
			return
		}

		if s.selectedEquipment >= len(g.equipment) {
			return
		}
		eq := g.equipment[s.selectedEquipment]

		fmt.Println("swap", eq.Name(), "to", g.selectedDude.Name())
	}
	g.ui.dudeInfoPanel.equipmentDetails.sellButton.hidden = false
	g.ui.dudeInfoPanel.equipmentDetails.onSellClick = func() {
		fmt.Println("argh, sell equipment")
	}
	g.ui.dudeInfoPanel.equipmentDetails.swapButton.hidden = false
	g.ui.dudeInfoPanel.equipmentDetails.onSwapClick = func() {
		fmt.Println("argh, swap equipment")
	}

	g.ui.dudePanel.buyButton.onClick = func() {
		s.BuyDude(g)
	}
	g.ui.dudePanel.buyButton.text.SetText(fmt.Sprintf("Random Dude\n%dgp", s.DudeCost()))
	g.ui.dudePanel.buyButton.disabled = false

	g.ui.roomPanel.buyButton.onClick = func() {
		s.RerollRooms(g)
	}
	g.ui.roomPanel.buyButton.text.SetText(fmt.Sprintf("Reroll Rooms\n%dgp", s.RerollCost()))
	g.ui.roomPanel.buyButton.disabled = false

	// I guess we can allow the player to yeet whenever.
	g.ui.buttonPanel.Enable()
	g.ui.buttonPanel.onClick = func() {
		emptyRooms := 0
		requiredRooms := 0
		for _, room := range s.nextStory.rooms {
			if room.kind == Empty {
				emptyRooms++
			}
		}
		for _, room := range s.availableRooms {
			if room.required {
				requiredRooms++
			}
		}

		// If there are required rooms, prevent the player from leaving.
		if requiredRooms > 0 {
			g.ui.feedback.Msg(FeedbackBad, "you must place all required rooms!")
			return
		}

		// If there are empty rooms, and we have rooms to place, ask for confirmation.
		if emptyRooms > 0 && len(s.availableRooms) > 0 {
			if s.readyAttempts > 0 {
				g.ui.feedback.Msg(FeedbackBad, "...so be it")
				s.readyAttempts = 2
				return
			} else {
				g.ui.feedback.Msg(FeedbackWarning, fmt.Sprintf("%d empty rooms remain, are you sure you want to proceed?", emptyRooms))
				s.readyAttempts++
				return
			}
		} else {
			s.readyAttempts = 2
		}
	}
	g.ui.buttonPanel.text.SetText("adventure forth!")
	g.ui.buttonPanel.hidden = false

	// Update info.
	g.UpdateInfo()
}
func (s *GameStateBuild) End(g *Game) {
	g.ui.roomPanel.onItemClick = nil
	g.ui.roomPanel.SetRoomDefs(nil)
	// Make sure rooms ain't highlighted
	for _, room := range s.nextStory.rooms {
		room.highlight = false
	}
	g.ui.buttonPanel.hidden = true
	g.ui.roomInfoPanel.hidden = true
	g.ui.roomPanel.buyButton.disabled = true
	g.ui.dudePanel.buyButton.disabled = true
	g.ui.equipmentPanel.buyButton.disabled = true
	g.ui.dudeInfoPanel.equipmentDetails.sellButton.hidden = true
	g.ui.dudeInfoPanel.equipmentDetails.swapButton.hidden = true
}
func (s *GameStateBuild) Update(g *Game) GameState {
	if s.readyAttempts >= 2 {
		return &GameStatePlay{}
	}

	// Some cancel garbo.
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.ui.roomInfoPanel.hidden = true
	}

	s.wobbler += 0.05
	s.titleTimer++

	if handled, kind := g.CheckUI(); !handled {
		if kind == UICheckHover {
			mx, my := IntToFloat2(ebiten.CursorPosition())

			// Center of screen.
			cx := float64(g.lastWidth) / 2
			cy := float64(g.lastHeight) / 2

			// if mouse is not within a bounds, unhighlight.
			buffer := 100.0 * g.camera.Zoom()
			if mx < cx-buffer || mx > cx+buffer || my < cy-buffer || my > cy+buffer {
				if s.focusedRoom != nil {
					s.focusedRoom.highlight = false
				}
				s.focusedRoom = nil
				return nil
			}

			// FIXME: This ain't right.
			cy -= float64(s.nextStory.level) * g.camera.GetMultiplier() * g.camera.Zoom()

			r := math.Atan2(my-cy, mx-cx) - g.camera.Rotation()
			roomIndex := s.nextStory.RoomIndexFromAngle(r)

			// Highlight all rooms equal to size of placing.
			for _, room := range s.highlightedRooms {
				room.highlight = false
			}
			if s.placingRoom != nil {
				s.highlightedRooms = nil
				for i := roomIndex; i < roomIndex+int(s.placingRoom.size) && i < 7; i++ {
					s.nextStory.rooms[i].highlight = true
					s.highlightedRooms = append(s.highlightedRooms, s.nextStory.rooms[i])
				}
			}

			// Ensure focusing our actual target root room.
			if s.focusedRoom != nil {
				s.focusedRoom.highlight = false
			}
			s.focusedRoom = s.nextStory.rooms[roomIndex]
			s.focusedRoom.highlight = true

		} else if kind == UICheckClick {
			s.TryPlaceRoom(g)
			s.readyAttempts = 0
		}
	} else {
		if s.focusedRoom != nil {
			s.focusedRoom.highlight = false
			s.focusedRoom = nil
		}
		for _, room := range s.highlightedRooms {
			room.highlight = false
		}
		s.highlightedRooms = nil
	}

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

func adjustSelectionIndex(index int, max int) int {
	// I'm lazy.
	if index >= max {
		index--
		if index < 0 {
			index = 0
		}
	}
	return index
}

func (s *GameStateBuild) TryPlaceRoom(g *Game) {
	if s.focusedRoom != nil {
		if s.placingRoom == nil {
			g.ui.feedback.Msg(FeedbackGeneric, "select a room to place :)")
		} else {
			// If it's not stairs or empty, allow the player to sell it back at full value.
			if s.focusedRoom.kind != Stairs && s.focusedRoom.kind != Empty {
				g.gold += GetRoomCost(s.focusedRoom.kind, s.focusedRoom.size, s.nextStory.level)
				s.availableRooms = append(s.availableRooms, GetRoomDef(s.focusedRoom.kind, s.focusedRoom.size, s.focusedRoom.required))
				s.availableRooms = SortRooms(s.availableRooms)
				// Reselect after sort
				s.placingIndex = adjustSelectionIndex(s.placingIndex, len(s.availableRooms))
				if s.placingIndex < len(s.availableRooms) {
					g.ui.roomPanel.onItemClick(s.placingIndex)
				}
				g.ui.roomPanel.SetRoomDefs(s.availableRooms)
				g.UpdateInfo()
				g.ui.feedback.Msg(FeedbackGood, fmt.Sprintf("%s %s sold!", s.focusedRoom.size.String(), s.focusedRoom.kind.String()))
				s.nextStory.RemoveRoom(s.focusedRoom.index)
			} else {
				if g.gold-GetRoomCost(s.placingRoom.kind, s.placingRoom.size, s.nextStory.level) < 0 {
					g.ui.feedback.Msg(FeedbackBad, "ur broke lol")
				} else {
					room := NewRoom(s.placingRoom.size, s.placingRoom.kind, s.placingRoom.required)
					if err := s.nextStory.PlaceRoom(room, s.focusedRoom.index); err != nil {
						g.ui.feedback.Msg(FeedbackBad, err.Error())
					} else {
						// it worked!11!
						g.gold -= GetRoomCost(s.placingRoom.kind, s.placingRoom.size, s.nextStory.level)
						g.UpdateInfo()
						g.ui.feedback.Msg(FeedbackGood, fmt.Sprintf("%s %s placed!", s.placingRoom.size.String(), s.placingRoom.kind.String()))

						// Remove from rooms.
						s.availableRooms = append(s.availableRooms[:s.placingIndex], s.availableRooms[s.placingIndex+1:]...)
						s.availableRooms = SortRooms(s.availableRooms)

						// Resync UI, I guess.
						g.ui.roomPanel.SetRoomDefs(s.availableRooms)
						g.ui.roomInfoPanel.hidden = true

						// I'm lazy.
						s.placingIndex = adjustSelectionIndex(s.placingIndex, len(s.availableRooms))
						if s.placingIndex < len(s.availableRooms) {
							g.ui.roomPanel.onItemClick(s.placingIndex)
						}
					}
				}
			}
		}
	}
}

func (s *GameStateBuild) RollRooms(g *Game) {
	s.availableRooms = nil
	numOptional := 5 + s.nextStory.level/3
	required := GetRequiredRooms(s.nextStory.level, 2)
	optional := GetOptionalRooms(s.nextStory.level, numOptional) // 6 is minimum, but let's given 3 more for fun.
	s.availableRooms = append(s.availableRooms, required...)
	s.availableRooms = append(s.availableRooms, optional...)
	s.availableRooms = SortRooms(s.availableRooms)

	// Update room panel.
	g.ui.roomPanel.SetRoomDefs(s.availableRooms)
}

func (s *GameStateBuild) RerollCost() int {
	return 25 + 75*(s.nextStory.level+1)
}

func (s *GameStateBuild) RerollRooms(g *Game) {
	cost := s.RerollCost()
	if g.gold < cost {
		g.ui.feedback.Msg(FeedbackBad, fmt.Sprintf("need more gold to reroll rooms! (%d)", cost))
		return
	}
	g.gold -= cost
	s.RollRooms(g)
	g.UpdateInfo()
}

func (s *GameStateBuild) DudeCost() int {
	return 100 + 50*(s.nextStory.level+1)
}

func (s *GameStateBuild) BuyDude(g *Game) {
	// COST?
	cost := s.DudeCost()
	if g.gold < cost {
		g.ui.feedback.Msg(FeedbackBad, fmt.Sprintf("need more gold to purchase a dude! (%d)", cost))
		return
	}
	g.gold -= cost
	level := len(g.tower.Stories) + 1
	if level < 1 {
		level = 1
	}

	// Random profession ??
	profession := RandomProfessionKind()
	dude := NewDude(profession, level)
	g.dudes = append(g.dudes, dude)
	g.UpdateInfo()
}

func (s *GameStateBuild) EquipmentCost() int {
	return 25 + 50*(s.nextStory.level+1)
}

func (s *GameStateBuild) BuyEquipment(g *Game) {
	// COST?
	cost := s.EquipmentCost()
	if g.gold < cost {
		g.ui.feedback.Msg(FeedbackBad, fmt.Sprintf("need more gold to purchase a equipment! (%d)", cost))
		return
	}
	g.gold -= cost

	level := len(g.tower.Stories)
	e := GetRandomEquipment(level)
	g.equipment = append(g.equipment, e)
	g.ui.equipmentPanel.SetEquipment(g.equipment)
	g.UpdateInfo()
}

func (s *GameStateBuild) SellEquipment(g *Game, e *Equipment) {
	if e == nil {
		return
	}
	found := false
	for i, eq := range g.equipment {
		if eq == e {
			g.equipment = append(g.equipment[:i], g.equipment[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return
	}
	value := int(e.GoldValue())
	g.gold += value
	AddMessage(MessageLoot, fmt.Sprintf("Sold %s for %d gold.", e.Name(), value))
	// Trigger on sell event
	e.Activate(EventSell{equipment: e})
	g.UpdateInfo()
}
