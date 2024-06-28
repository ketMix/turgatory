package game

import (
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Story is a single story in the tower. It contains rooms.
type Story struct {
	rooms       []*Room // Rooms represent the counter-clockwise "pie" of rooms. This field is sized according to the capacity of the story (which is assumed to always be 8, but not necessarily).
	dudes       []*Dude
	portalStack *render.Stack
	stacks      render.Stacks
	doorStack   *render.Stack
	walls       render.Stacks
	vgroup      *render.VGroup
	level       int
	open        bool
	texts       []FloatingText
}

// StoryHeight is the height of a story in da tower.
const StoryHeight = 28                                      // StoryHeight is used to space stories apart from each other vertically.
const StorySlices = 28                                      // The amount of slices used for the frame buffers, should be equal to maximum staxie slice count used in a story.
const StoryWallHeight = 9                                   // The height of the wall stack -- this is repeated 3 times to get the full height (roughly)
const StoryVGroupWidth = 256                                // Framebuffer's maximum width for rendering.
const StoryVGroupHeight = 256                               // Framebuffer's maximum height for rendering.
const PortalDistance = 44                                   // Distance from the center of the tower to the portal.
const PortalRotationMedium = 7.0*-(math.Pi/8) - math.Pi/3.5 // Rotation of the portal.
const PortalRotationSmall = 7.0*-(math.Pi/8) - math.Pi/2.5  // Rotation of the portal.
const TowerCenterX = StoryVGroupWidth/2 - 5                 // Center of the tower. Have to offset lightly for some dumb reason...
const TowerCenterY = StoryVGroupHeight/2 - 5                // Center of the tower.

// NewStory creates a grand new spankin' story.
func NewStory() *Story {
	return NewStoryWithSize(8)
}
func NewStoryWithSize(size int) *Story {
	story := &Story{}
	story.rooms = make([]*Room, size)

	// Fill with template rooms
	for i := 0; i < size-1; i++ {
		story.PlaceRoom(NewRoom(Small, Empty, false), i)
	}
	// Place entrance/exit room
	room := NewRoom(Small, Stairs, true)
	PanicIfErr(story.PlaceRoom(room, 7))

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

	// Add our closed door to our final room position.
	{
		stack := Must(render.NewStack("walls/door", "", ""))
		x := float64(StoryVGroupWidth) / 2
		y := float64(StoryVGroupHeight) / 2
		stack.SetPosition(x, y)
		stack.SetRotation(float64(7) * -(math.Pi / 4))
		story.doorStack = stack
	}

	// Add our center pillar
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
func (s *Story) Update(req *ActivityRequests, g *Game) {
	for i := 0; i < len(s.texts); i++ {
		s.texts[i].Update()
		if !s.texts[i].Alive() {
			s.texts = append(s.texts[:i], s.texts[i+1:]...)
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
			room.Update(req, g) // Just forward up room updates to tower
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
			if room != u.dude.room {
				// Add dude to the given room and remove from existing room.
				if u.dude.room != nil {
					u.dude.room.RemoveDude(u.dude)
				}
				if room != nil {
					room.AddDude(u.dude)
				}
				req.Add(RoomEnterActivity{dude: u.dude, room: room})
			}
			// Check if the dude is in the center of the room and update as appropriate.
			if room != nil {
				// Special case for boss room
				if room.kind == Boss && !room.killedBoss {
					// If dude is in first fourth of boss room, add it to the waiting list
					if s.IsInCenterOfRoom(s.AngleFromCenter(u.x, u.y)-math.Pi/4, roomIndex) {
						if !room.IsDudeWaiting(u.dude) {
							room.AddDudeToWaiting(u.dude)
							req.Add(RoomWaitActivity{dude: u.dude, room: room})
						}
					}
				} else {
					if s.IsInCenterOfRoom(s.AngleFromCenter(u.x, u.y), roomIndex) {
						if !room.IsDudeInCenter(u.dude) {
							room.AddDudeToCenter(u.dude)
							req.Add(RoomCenterActivity{dude: u.dude, room: room})
						}
					} else if room.IsDudeInCenter(u.dude) {
						room.RemoveDudeFromCenter(u.dude)
						req.Add(RoomEndActivity{dude: u.dude, room: room})
					}
				}
			}
		case StoryEnterNextActivity:
			// Pass it up but with our story added.
			req.Add(StoryEnterNextActivity{dude: u.dude, story: s})
		case TowerLeaveActivity:
			// Pass it up.
			req.Add(u)
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

		// Draw the door stack.
		if s.doorStack != nil {
			s.doorStack.Draw(opts)
		}

		if s.portalStack != nil {
			s.portalStack.Draw(opts)
		}

	}

	s.vgroup.Draw(o)

	textOpts := render.Options{
		Camera:        o.Camera,
		Screen:        o.Overlay, // Use overlay for render target
		Pitch:         o.Pitch,
		VGroup:        s.vgroup,
		TowerRotation: o.TowerRotation,
	}
	textOpts.DrawImageOptions.GeoM.Concat(o.DrawImageOptions.GeoM)
	for _, text := range s.texts {
		text.Draw(&textOpts)
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
			// Might as well remove dude from rooms.
			if d.room != nil {
				d.room.RemoveDudeFromCenter(d)
				d.room.RemoveDude(d)
			}
			s.dudes = append(s.dudes[:i], s.dudes[i+1:]...)
			return
		}
	}
}

// Open marks the story as being open. This activates full updates and rendering.
func (s *Story) Open() {
	s.open = true
}

func (s *Story) RemoveDoor() {
	s.doorStack = nil
}

func (s *Story) AddPortal() {
	stack := Must(render.NewStack("walls/portal", "", ""))

	r := PortalRotationMedium
	switch s.rooms[6].size {
	case Small:
		r = PortalRotationSmall
	}

	x, y := s.PositionFromCenter(r, PortalDistance)

	stack.SetRotation(math.Pi + math.Pi/2 - math.Pi/8)
	stack.SetPosition(x, y)
	stack.NoLighting = true
	stack.VgroupOffset = 2
	s.portalStack = stack
}

func (s *Story) RemovePortal() {
	s.portalStack = nil
}

func (s *Story) Reset() {
	for _, room := range s.rooms {
		if room != nil {
			room.Reset()
		}
	}
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
		existingRoom := s.rooms[index+i]
		if existingRoom != nil && existingRoom.kind != Empty {
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
		s.PlaceRoom(NewRoom(Small, Empty, false), index+i)
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

func (s *Story) GetRoomCenterRad(roomIndex int) float64 {
	room := s.rooms[roomIndex]
	if room == nil {
		return 0
	}
	start := float64(room.index) * (math.Pi / 4)
	end := start + float64(room.size)*(math.Pi/4)
	return (start + end) / 2
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

func (s *Story) AddText(t FloatingText) {
	s.texts = append(s.texts, t)
}

const (
	ErrRoomIndexInvalid = Error("invalid room index")
	ErrRoomTooLarge     = Error("room too large")
	ErrRoomNoSpace      = Error("no space for room")
	ErrRoomNotPresent   = Error("room not present")
)
