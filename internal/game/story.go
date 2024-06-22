package game

import (
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Story is a single story in the tower. It contains rooms.
type Story struct {
	rooms  []*Room // Rooms represent the counter-clockwise "pie" of rooms. This field is sized according to the capacity of the story (which is assumed to always be 8, but not necessarily).
	dudes  []*Dude
	stacks render.Stacks
	vgroup *render.VGroup
	level  int
}

// StoryHeight is the height of a story in da tower.
const StoryHeight = 28        // StoryHeight is used to space stories apart from each other vertically.
const StorySlices = 28        // The amount of slices used for the frame buffers, should be equal to maximum staxie slice count used in a story.
const StoryVGroupWidth = 256  // Framebuffer's maximum width for rendering.
const StoryVGroupHeight = 256 // Framebuffer's maximum height for rendering.

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

	room := NewRoom(Small, Combat)
	PanicIfErr(story.PlaceRoom(room, 4))

	room2 := NewRoom(Medium, Armory)
	PanicIfErr(story.PlaceRoom(room2, 0))

	/*{
		center := Must(render.NewStack("rooms/center", "", ""))
		center.SetPosition(float64(StoryVGroupWidth)/2-16, float64(StoryVGroupHeight)/2-16)
		center.SetOriginToCenter()
		story.stacks.Add(center)
	}*/

	// Test dude
	dude := NewDude()
	dude.stack.SetPosition(story.PositionFromCenter(math.Pi/2, TowerEntrance))
	story.AddDude(dude)

	story.vgroup = render.NewVGroup(StoryVGroupWidth, StoryVGroupHeight, StorySlices) // For now...

	return story
}

// Update updates the rooms.
func (s *Story) Update(req *ActivityRequests) {
	// Update the floors in case they have sweet animations.
	for _, stack := range s.stacks {
		stack.Update()
		//stack.SetRotation(stack.Rotation() + 0.01) // Spin the floors. FIXME: Camera no longer works due to fake perspective trick, so we spin here.
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
			room.Update()
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
			if room := s.rooms[roomIndex]; room != u.initiator.Room() {
				req.Add(RoomEnterActivity{initiator: u.initiator, room: room})
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
		Screen: o.Screen,
		Pitch:  o.Pitch,
		VGroup: s.vgroup,
	}

	// We can't use the camera's own functionality, so we do it ourselves here.
	opts.DrawImageOptions.GeoM.Translate(-StoryVGroupWidth/2, -StoryVGroupHeight/2)
	opts.DrawImageOptions.GeoM.Rotate(o.TowerRotation)
	opts.DrawImageOptions.GeoM.Translate(StoryVGroupWidth/2, StoryVGroupHeight/2)

	s.stacks.Draw(opts)

	for _, room := range s.rooms {
		if room != nil {
			room.Draw(opts)
		}
	}

	for _, dude := range s.dudes {
		dude.Draw(opts)
	}

	s.vgroup.Draw(o)
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
	} else if r.size == Huge {
		r.stacks.SetPositions(x, y-HugeOriginY)
		r.stacks.SetOrigins(0, HugeOriginY)
	} else {
		r.stacks.SetPositions(x, y)
	}

	r.stacks.SetRotations(float64(index) * -(math.Pi / 4)) // We go counter-clockwise...
	for i := 0; i < int(r.size); i++ {
		s.rooms[index+i] = r
	}
	r.story = s
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

func (s *Story) AngleFromCenter(x, y float64) float64 {
	cx := float64(StoryVGroupWidth) / 2
	cy := float64(StoryVGroupHeight) / 2
	return math.Atan2(y-cy, x-cx)
}

func (s *Story) PositionFromCenter(rads float64, amount float64) (float64, float64) {
	cx := float64(StoryVGroupWidth) / 2
	cy := float64(StoryVGroupHeight) / 2
	return cx + math.Cos(rads)*amount, cy + math.Sin(rads)*amount
}

func (s *Story) DistanceFromCenter(x, y float64) float64 {
	cx := float64(StoryVGroupWidth) / 2
	cy := float64(StoryVGroupHeight) / 2
	return math.Sqrt(math.Pow(x-cx, 2) + math.Pow(y-cy, 2))
}

const (
	ErrRoomIndexInvalid = Error("invalid room index")
	ErrRoomTooLarge     = Error("room too large")
	ErrRoomNoSpace      = Error("no space for room")
	ErrRoomNotPresent   = Error("room not present")
)
