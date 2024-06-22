package game

type EventName int

const (
	EventUnknown   EventName = iota
	EventEnterRoom           // EventEnterRoom is triggered when a dude enters a room
	EventLeaveRoom           // EventLeaveRoom is triggered when a dude leaves a room
	EventEquip               // EventEquip is triggered when a dude equips an item
	EventUnequip             // EventUnequip is triggered when a dude unequips an item
)

func (e EventName) String() string {
	switch e {
	case EventEnterRoom:
		return "Enter Room"
	case EventLeaveRoom:
		return "Leave Room"
	case EventEquip:
		return "Equip"
	case EventUnequip:
		return "Unequip"
	case EventUnknown:
	default:
	}
	return "Unknown"
}

type IEvent interface {
	Name() EventName
	Data() interface{}
}

type Event struct {
	name EventName
	data interface{}
}

func (e *Event) Name() EventName {
	return e.name
}

func (e *Event) Data() interface{} {
	return e.data
}

func NewEnterRoomEvent(room *Room) *Event {
	return &Event{
		name: EventEnterRoom,
		data: room,
	}
}

func NewLeaveRoomEvent(room *Room) *Event {
	return &Event{
		name: EventLeaveRoom,
		data: room,
	}
}

func NewEquipEvent(equipment *Equipment) *Event {
	return &Event{
		name: EventEquip,
		data: equipment,
	}
}

func NewUnequipEvent(equipment *Equipment) *Event {
	return &Event{
		name: EventUnequip,
		data: equipment,
	}
}
