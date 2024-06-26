package game

type ActivityRequests []Activity

func (a *ActivityRequests) Add(act Activity) {
	*a = append(*a, act)
}

type Activity interface {
	Apply()
	Cb() func(success bool)
}

type MoveActivity struct {
	dude *Dude
	face float64
	x    float64
	y    float64
	cb   func(success bool)
}

func (m MoveActivity) Apply() {
	m.dude.SetPosition(m.x, m.y)
	m.dude.SetRotation(m.face)
}

func (m MoveActivity) Cb() func(success bool) {
	return m.cb
}

type StoryEnterNextActivity struct {
	dude  *Dude
	story *Story
	cb    func(success bool)
}

func (s StoryEnterNextActivity) Apply() {
}

func (s StoryEnterNextActivity) Cb() func(success bool) {
	return s.cb
}

type TowerLeaveActivity struct {
	dude *Dude
	cb   func(success bool)
}

func (t TowerLeaveActivity) Apply() {
}

func (t TowerLeaveActivity) Cb() func(success bool) {
	return t.cb
}

type TowerCompleteActivity struct {
	dude *Dude
	cb   func(success bool)
}

func (t TowerCompleteActivity) Apply() {
}

func (t TowerCompleteActivity) Cb() func(success bool) {
	return t.cb
}

type RoomEnterActivity struct {
	dude *Dude
	room *Room
	cb   func(success bool)
}

func (r RoomEnterActivity) Apply() {
	r.dude.SetRoom(r.room)
}

func (r RoomEnterActivity) Cb() func(success bool) {
	return r.cb
}

type RoomLeaveActivity struct {
	dude *Dude
	room *Room
	cb   func(success bool)
}

func (r RoomLeaveActivity) Apply() {
	r.dude.SetRoom(nil)
}

func (r RoomLeaveActivity) Cb() func(success bool) {
	return r.cb
}

type RoomCenterActivity struct {
	dude *Dude
	room *Room
	cb   func(success bool)
}

func (r RoomCenterActivity) Apply() {
}

func (r RoomCenterActivity) Cb() func(success bool) {
	return r.cb
}

type RoomEndActivity struct {
	dude     *Dude
	room     *Room
	lastRoom bool
	cb       func(success bool)
}

func (r RoomEndActivity) Apply() {
}

func (r RoomEndActivity) Cb() func(success bool) {
	return r.cb
}

type RoomCombatActivity struct {
	dude *Dude
	room *Room
}

func (r RoomCombatActivity) Apply() {
}

func (r RoomCombatActivity) Cb() func(success bool) {
	return nil
}
