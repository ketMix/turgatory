package game

type ActivityRequests []Activity

func (a *ActivityRequests) Add(act Activity) {
	*a = append(*a, act)
}

type Activity interface {
	Apply()
	Cb() func(success bool)
}

type PositionRotator interface {
	Position() (float64, float64)
	SetPosition(x, y float64)
	Rotation() float64
	SetRotation(r float64)
}

type Actor interface {
	PositionRotator
	Name() string
	Room() *Room
	SetRoom(r *Room)
}

type MoveActivity struct {
	initiator Actor
	face      float64
	x         float64
	y         float64
	cb        func(success bool)
}

func (m MoveActivity) Apply() {
	m.initiator.SetPosition(m.x, m.y)
	m.initiator.SetRotation(m.face)
}

func (m MoveActivity) Cb() func(success bool) {
	return m.cb
}

type RoomEnterActivity struct {
	initiator Actor
	room      *Room
	cb        func(success bool)
}

func (r RoomEnterActivity) Apply() {
	r.initiator.SetRoom(r.room)
}

func (r RoomEnterActivity) Cb() func(success bool) {
	return r.cb
}

type RoomLeaveActivity struct {
	initiator Actor
	room      *Room
	cb        func(success bool)
}

func (r RoomLeaveActivity) Apply() {
	r.initiator.SetRoom(nil)
}

func (r RoomLeaveActivity) Cb() func(success bool) {
	return r.cb
}
