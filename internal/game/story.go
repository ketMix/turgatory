package game

import (
	"image/color"
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Story is a single story in the tower. It contains rooms.
type Story struct {
	rooms     []*Room // Rooms represent the counter-clockwise "pie" of rooms. This field is sized according to the capacity of the story (which is assumed to always be 8, but not necessarily).
	dudes     []*Dude
	stacks    render.Stacks
	walls     render.Stacks
	vgroup    *render.VGroup
	level     int
	open      bool
	text      []FloatingText
	textTimer int
}

// StoryHeight is the height of a story in da tower.
const StoryHeight = 28                       // StoryHeight is used to space stories apart from each other vertically.
const StorySlices = 28                       // The amount of slices used for the frame buffers, should be equal to maximum staxie slice count used in a story.
const StoryWallHeight = 9                    // The height of the wall stack -- this is repeated 3 times to get the full height (roughly)
const StoryVGroupWidth = 256                 // Framebuffer's maximum width for rendering.
const StoryVGroupHeight = 256                // Framebuffer's maximum height for rendering.
const TowerCenterX = StoryVGroupWidth/2 - 5  // Center of the tower. Have to offset lightly for some dumb reason...
const TowerCenterY = StoryVGroupHeight/2 - 5 // Center of the tower.

// NewStory creates a grand new spankin' story.
func NewStory() *Story {
	return NewStoryWithSize(8)
}
func NewStoryWithSize(size int) *Story {
	story := &Story{}
	story.rooms = make([]*Room, size)

	for i := 0; i < 4; i++ {
		stack := Must(render.NewStack("walls/pie", "template", ""))
		stack.SetRotation(float64(i) * (math.Pi / 2))

		// This feels hacky atm, but position from the center of our vgroup.
		x := float64(StoryVGroupWidth) / 2
		y := float64(StoryVGroupHeight) / 2
		stack.SetPosition(x, y)

		story.stacks.Add(stack)
	}

	// Add our walls.
	for j := 0; j < 3; j++ {
		for i := 0; i < 8; i++ {
			stack := Must(render.NewStack("walls/exterior", "", ""))
			x := float64(StoryVGroupWidth) / 2
			y := float64(StoryVGroupHeight) / 2
			stack.VgroupOffset = j * StoryWallHeight
			stack.SetPosition(x, y)
			stack.SetRotation(float64(i) * (math.Pi / 4))
			story.walls.Add(stack)
		}
	}

	room2 := NewRoom(Medium, Armory)
	PanicIfErr(story.PlaceRoom(room2, 0))

	room3 := NewRoom(Small, Combat)
	PanicIfErr(story.PlaceRoom(room3, 2))

	room4 := NewRoom(Small, Treasure)
	PanicIfErr(story.PlaceRoom(room4, 3))

	room5 := NewRoom(Small, HealingShrine)
	PanicIfErr(story.PlaceRoom(room5, 4))

	room6 := NewRoom(Small, Well)
	PanicIfErr(story.PlaceRoom(room6, 5))

	room7 := NewRoom(Medium, Curse)
	PanicIfErr(story.PlaceRoom(room7, 6))

	{
		center := Must(render.NewStack("rooms/center", "", ""))
		center.SetPosition(float64(StoryVGroupWidth)/2-16, float64(StoryVGroupHeight)/2-16)
		center.SetOriginToCenter()
		story.stacks.Add(center)
	}

	story.vgroup = render.NewVGroup(StoryVGroupWidth, StoryVGroupHeight, StorySlices) // For now...

	return story
}

// Update updates the rooms.
func (s *Story) Update(req *ActivityRequests) {
	if s.textTimer <= 0 {
		s.textTimer = 30
		t := MakeFloatingText("ok", color.NRGBA{255, 255, 255, 255}, 30)
		t.SetOrigin(0, 60)
		t.YOffset = 11

		s.text = append(s.text, t)
	} else {
		s.textTimer--
	}
	for i := 0; i < len(s.text); i++ {
		s.text[i].Update()
		if !s.text[i].Alive() {
			s.text = append(s.text[:i], s.text[i+1:]...)
			i--
		}
	}

	// Update the floors in case they have sweet animations.
	for _, stack := range s.stacks {
		stack.Update()
		//stack.SetRotation(stack.Rotation() + 0.01) // Spin the floors. FIXME: Camera no longer works due to fake perspective trick, so we spin here.
	}

	// Bail if the story is not yet open.
	if !s.open {
		return
	}

	// Update the rooms.
	var updatedRooms []*Room
	for _, room := range s.rooms {
		if room != nil {
			for _, updatedRoom := range updatedRooms {
				if updatedRoom == room {
					continue
				}
			}
			room.Update(req) // Just forward up room updates to tower
			updatedRooms = append(updatedRooms, room)
		}
	}

	// Update the dudes (maybe should be just handled in rooms?)
	var dudeUpdates ActivityRequests
	for _, dude := range s.dudes {
		dude.Update(s, &dudeUpdates)
	}
	for _, u := range dudeUpdates {
		success := true
		switch u := u.(type) {
		case MoveActivity:
			roomIndex := s.RoomIndexFromAngle(s.AngleFromCenter(u.x, u.y))
			room := s.rooms[roomIndex]
			if room != u.initiator.Room() {
				// FIXME: Remove the Actor concept and assume Dude for all
				// Add dude to the given room and remove from existing room.
				if d, ok := u.initiator.(*Dude); ok {
					if d.room != nil {
						d.room.RemoveDude(d)
					}
					if room != nil {
						room.AddDude(d)
					}
				}
				req.Add(RoomEnterActivity{initiator: u.initiator, room: room})
			}
			// Check if the initiator is in the center of the room and update as appropriate.
			if room != nil {
				if s.IsInCenterOfRoom(s.AngleFromCenter(u.x, u.y), roomIndex) {
					if !room.IsActorInCenter(u.initiator) {
						room.AddActorToCenter(u.initiator)
						req.Add(RoomCenterActivity{initiator: u.initiator, room: room})
					}
				} else if room.IsActorInCenter(u.initiator) {
					room.RemoveActorFromCenter(u.initiator)
					req.Add(RoomEndActivity{initiator: u.initiator, room: room})
				}
			}
		}
		if success {
			u.Apply()
		}
		if cb := u.Cb(); cb != nil {
			cb(success)
		}
	}
}

// Draw draws the rooms.
func (s *Story) Draw(o *render.Options) {
	s.vgroup.Clear()

	opts := &render.Options{
		Camera:        o.Camera,
		Screen:        o.Screen,
		Pitch:         o.Pitch,
		VGroup:        s.vgroup,
		TowerRotation: o.TowerRotation,
	}

	// We can't use the camera's own functionality, so we do it ourselves here.
	opts.DrawImageOptions.GeoM.Translate(-StoryVGroupWidth/2, -StoryVGroupHeight/2)
	opts.DrawImageOptions.GeoM.Rotate(o.TowerRotation)
	opts.DrawImageOptions.GeoM.Translate(StoryVGroupWidth/2, StoryVGroupHeight/2)

	// If the story is not yet open, just draw the tower exterior stacks.
	if !s.open {
		s.walls.Draw(opts)
	} else {
		// Conditionally render the walls based upon rotation.
		for _, stack := range s.walls {
			r := stack.Rotation() + o.TowerRotation
			r += math.Pi / 2

			// Ensure r is constrained from 0 to 2*math.Pi
			for r < 0 {
				r += math.Pi * 2
			}
			for r >= math.Pi*2 {
				r -= math.Pi * 2
			}

			min := math.Pi / 4
			max := math.Pi * 4 / 4

			opts.DrawImageOptions.ColorScale.Reset()

			if r >= min && r < max {
				continue
			} else if r >= min-math.Pi/4 && r < max+math.Pi/4 {
				opts.DrawImageOptions.ColorScale.ScaleAlpha(0.25)
			}

			stack.Draw(opts)
		}
		opts.DrawImageOptions.ColorScale.Reset()

		s.stacks.Draw(opts)

		for _, room := range s.rooms {
			if room != nil {
				room.Draw(opts)
			}
		}

		for _, dude := range s.dudes {
			dude.Draw(opts)
		}
	}

	s.vgroup.Draw(o)

	for _, text := range s.text {
		opts2 := render.Options{
			Camera:        o.Camera,
			Screen:        o.Screen,
			Pitch:         o.Pitch,
			VGroup:        s.vgroup,
			TowerRotation: o.TowerRotation,
		}

		opts2.DrawImageOptions.GeoM.Concat(o.DrawImageOptions.GeoM)

		text.Draw(&opts2)
	}
}

// Complete returns if the story is considered complete based upon full room saturation.
func (s *Story) Complete() bool {
	for _, room := range s.rooms {
		if room == nil {
			return false
		}
	}
	return true
}

func (s *Story) AddDude(d *Dude) {
	d.story = s
	s.dudes = append(s.dudes, d)
}

func (s *Story) RemoveDude(d *Dude) {
	for i, v := range s.dudes {
		if v == d {
			d.story = nil
			s.dudes = append(s.dudes[:i], s.dudes[i+1:]...)
			return
		}
	}
}

// Open marks the story as being open. This activates full updates and rendering.
func (s *Story) Open() {
	s.open = true
}

// PlaceRoom places a room in the story, populating the rooms slice's pointer references accordingly.
func (s *Story) PlaceRoom(r *Room, index int) error {
	if index < 0 || index >= len(s.rooms) {
		return ErrRoomIndexInvalid
	}
	if index+int(r.size) > len(s.rooms) {
		return ErrRoomTooLarge
	}
	for i := 0; i < int(r.size); i++ {
		if s.rooms[index+i] != nil {
			return ErrRoomNoSpace
		}
	}

	x := float64(StoryVGroupWidth) / 2
	y := float64(StoryVGroupHeight) / 2

	// Assign position and origin offsets as needed due to different pie sizes having different origins. FIXME: This should be controlled by the room itself, not here. Additionally, using the "plural" stack assignment won't apply to any stacks that are added to the room (such as monsters, etc.), so this should be readjusted so that it applies to DrawImageOptions in Draw() or similar.
	if r.size == Large {
		r.stacks.SetPositions(x, y-LargeOriginY)
		r.stacks.SetOrigins(0, LargeOriginY)
		r.walls.SetPositions(x, y-LargeOriginY)
		r.walls.SetOrigins(0, LargeOriginY)
	} else if r.size == Huge {
		r.stacks.SetPositions(x, y-HugeOriginY)
		r.stacks.SetOrigins(0, HugeOriginY)
		r.walls.SetPositions(x, y-HugeOriginY)
		r.walls.SetOrigins(0, HugeOriginY)
	} else {
		r.stacks.SetPositions(x, y)
		r.walls.SetPositions(x, y)
	}

	r.stacks.SetRotations(float64(index) * -(math.Pi / 4)) // We go counter-clockwise...
	r.walls.SetRotations(float64(index) * -(math.Pi / 4))  // We go counter-clockwise...
	for i := 0; i < int(r.size); i++ {
		s.rooms[index+i] = r
	}
	r.story = s
	r.index = index
	return nil
}

// RemoveRoom removes the room at the given index. This function always gets the "head" of the room and removes its size, regardless of if the target index is not the head.
func (s *Story) RemoveRoom(index int) error {
	if index < 0 || index >= len(s.rooms) {
		return ErrRoomIndexInvalid
	}
	if s.rooms[index] == nil {
		return ErrRoomNotPresent
	}
	room := s.rooms[index]
	// Get the "head" of the room and use that instead. Obviously this is redundant for single unit rooms.
	for i := 0; i < len(s.rooms); i++ {
		if s.rooms[i] == room {
			index = i
			break
		}
	}
	// Clear out the room references.
	for i := 0; i < int(room.size); i++ {
		s.rooms[index+i] = nil
	}
	room.story = nil
	room.index = -1
	return nil
}

// RoomIndexFromAngle returns the room index based upon the radians provided.
func (s *Story) RoomIndexFromAngle(rads float64) int {
	rads -= math.Pi / 2 // Adjust a lil
	// We go counter-clockwise, so we need to reverse the angle.
	rads = -rads
	if rads < 0 {
		rads += math.Pi * 2
	}
	return int(math.Floor(rads/(math.Pi/4))) % len(s.rooms)
}

func (s *Story) IsInCenterOfRoom(rads float64, roomIndex int) bool {
	room := s.rooms[roomIndex]
	if room == nil {
		return false
	}
	rads -= math.Pi / 2 // Adjust a lil
	rads = -rads
	if rads < 0 {
		rads += math.Pi * 2
	}

	start := float64(room.index) * (math.Pi / 4)
	end := start + float64(room.size)*(math.Pi/4)

	center := (start + end) / 2
	centerRad := math.Pi / 16
	centerStart := center - centerRad
	centerEnd := center + centerRad
	return rads >= centerStart && rads <= centerEnd
}

func (s *Story) AngleFromCenter(x, y float64) float64 {
	cx := float64(TowerCenterX)
	cy := float64(TowerCenterY)
	return math.Atan2(y-cy, x-cx)
}

func (s *Story) PositionFromCenter(rads float64, amount float64) (float64, float64) {
	cx := float64(TowerCenterX)
	cy := float64(TowerCenterY)
	return cx + math.Cos(rads)*amount, cy + math.Sin(rads)*amount
}

func (s *Story) DistanceFromCenter(x, y float64) float64 {
	cx := float64(TowerCenterX)
	cy := float64(TowerCenterY)
	return math.Sqrt(math.Pow(x-cx, 2) + math.Pow(y-cy, 2))
}

const (
	ErrRoomIndexInvalid = Error("invalid room index")
	ErrRoomTooLarge     = Error("room too large")
	ErrRoomNoSpace      = Error("no space for room")
	ErrRoomNotPresent   = Error("room not present")
)
