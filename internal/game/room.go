package game

// RoomSize represents the different sizes of rooms, with 1 equating to 1/8th of a circle.
type RoomSize int

const (
	Small  RoomSize = 1
	Medium RoomSize = 2
	Large  RoomSize = 3
	Huge   RoomSize = 4
)

// RoomKind is an enumeration of the different kinds of rooms in za toweru.
type RoomKind int

const (
	Empty RoomKind = iota
	// Armory provide... armor up? damage up? Maybe should be different types.
	Armory
	// Healing shrine heals the adventurers over time.
	HealingShrine
	// Combat is where it goes down. $$$ is acquired.
	Combat
)

// Room is a room within a story of za toweru.
type Room struct {
	story *Story // Might as well keep for reference.
	kind  RoomKind
	size  RoomSize
	power int // ???
}
