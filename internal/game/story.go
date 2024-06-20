package game

import (
	"math"

	"github.com/kettek/ebijam24/internal/render"
)

// Story is a single story in the tower. It contains rooms.
type Story struct {
	rooms       []*Room // Rooms represent the counter-clockwise "pie" of rooms. This field is sized according to the capacity of the story (which is assumed to always be 8, but not necessarily).
	floorStacks []*render.Stack
}

// StoryHeight is the height of a story in da tower.
const StoryHeight = 16

// NewStory creates a grand new story.
func NewStory(size int) *Story {
	story := &Story{}
	story.rooms = make([]*Room, size)

	for i := 0; i < 4; i++ {
		stack, err := render.NewStack("walls/pie", "", "")
		if err != nil {
			continue
		}
		stack.SetRotation(float64(i) * (math.Pi / 2))
		stack.SetRotationDistance(0)
		story.floorStacks = append(story.floorStacks, stack)
	}

	return story
}

// Update updates the rooms.
func (s *Story) Update() {
	// Update the floors in case they have sweet animations.
	for _, stack := range s.floorStacks {
		stack.Update()
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

	// TODO: We need to sort our draw operations based upon the rotation of the camera vs. the rotation+position of the stack...
	for _, stack := range s.floorStacks {
		stack.Draw(o)
	}
	// NOTE: Maybe we should actually apply room GeoM rotations here? This would make it so everything in a room is transformed appropriately. Technically, we could store unit/thing/positions in cartesian coordinates, then transform them into polar coordinates(???)???
	for _, room := range s.rooms {
		if room != nil {
			room.Draw(o)
		}
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
