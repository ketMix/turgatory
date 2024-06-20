package game

import (
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Story is a single story in the tower. It contains rooms.
type Story struct {
	rooms  []*Room // Rooms represent the counter-clockwise "pie" of rooms. This field is sized according to the capacity of the story (which is assumed to always be 8, but not necessarily).
	stacks render.Stacks
	vgroup *render.VGroup
}

// StoryHeight is the height of a story in da tower.
const StoryHeight = 20        // StoryHeight is used to space stories apart from each other vertically.
const StorySlices = 20        // The amount of slices used for the frame buffers, should be equal to maximum staxie slice count used in a story.
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
		stack, err := render.NewStack("walls/pie", "", "")
		if err != nil {
			continue
		}
		stack.SetRotation(float64(i) * (math.Pi / 2))

		// This feels hacky atm, but position from the center of our vgroup.
		x := float64(StoryVGroupWidth) / 2
		y := float64(StoryVGroupHeight) / 2
		stack.SetPosition(x, y)

		story.stacks.Add(stack)
	}

	// Test
	{
		dude, err := render.NewStack("dudes/liltest", "", "")
		if err != nil {
			panic(err)
		}
		dude.SetPosition(StoryVGroupWidth/2-32, StoryVGroupHeight/2-32)
		story.stacks = append(story.stacks, dude)
	}

	story.vgroup = render.NewVGroup(StoryVGroupWidth, StoryVGroupHeight, StorySlices) // For now...

	return story
}

// Update updates the rooms.
func (s *Story) Update() {
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
	for i := 0; i < int(r.size); i++ {
		s.rooms[index+i] = r
	}
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
	return nil
}

const (
	ErrRoomIndexInvalid = Error("invalid room index")
	ErrRoomTooLarge     = Error("room too large")
	ErrRoomNoSpace      = Error("no space for room")
	ErrRoomNotPresent   = Error("room not present")
)
