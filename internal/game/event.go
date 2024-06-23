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

// EventCenterRoom is triggered when a dude is roughly in the center of a room
type EventCenterRoom struct {
	room *Room
	dude *Dude
}

func (e EventCenterRoom) String() string {
	return "Center of Room"
}

// EventEndRoom is triggered when a dude is near the end part of a room.
type EventEndRoom struct {
	room *Room
	dude *Dude
}

func (e EventEndRoom) String() string {
	return "End of Room"
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

// EventSell occurs when an equipment is sold
type EventSell struct {
	tower     *Tower
	equipment *Equipment
}

func (e EventSell) String() string {
	return "Sell"
}

// EventGoldGain occurs when a dude gains gold
type EventGoldGain struct {
	dude   *Dude
	amount float32
}

func (e EventGoldGain) String() string {
	return "Gold Gain"
}

// EventGoldGain occurs when a dude gains gold
type EventGoldLoss struct {
	dude   *Dude
	amount float32
}

func (e EventGoldLoss) String() string {
	return "Gold Loss"
}
