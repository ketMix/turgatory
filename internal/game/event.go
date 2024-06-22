package game

// Event represents when something in the game has happened, such as entering a room, a dude dying, etc.
type Event interface {
	String() string
}

// EventEnterRoom is triggered when a dude enters a room
type EventEnterRoom struct {
	room *Room
	dude *Dude
}

func (e EventEnterRoom) String() string {
	return "Enter Room"
}

// EventLeaveRoom is triggered when a dude leaves a room
type EventLeaveRoom struct {
	room *Room
	dude *Dude
}

func (e EventLeaveRoom) String() string {
	return "Leave Room"
}

// EventEquip is triggered when a dude equips an item
type EventEquip struct {
	dude      *Dude
	equipment *Equipment
}

func (e EventEquip) String() string {
	return "Equip"
}

// EventUnequip is triggered when a dude unequips an item
type EventUnequip struct {
	dude      *Dude
	equipment *Equipment
}

func (e EventUnequip) String() string {
	return "Unequip"
}
