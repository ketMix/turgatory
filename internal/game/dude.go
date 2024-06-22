package game

import (
	"math"
	"math/rand"

	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type DudeActivity int

const (
	Idle          DudeActivity = iota
	FirstEntering              // First entering the tower.
	GoingUp                    // Entering the room from a staircase, this basically does the fancy slice offset/limiting.
	Centering                  // Move the dude to the center of the room.
	Moving                     // Move the dude counter-clockwise.
	Leaving                    // Move the dude to the stairs.
	GoingDown                  // Leaving the room to the stairs, opposite of GoingUp.
)

type Dude struct {
	name           string
	xp             int
	gold           float32
	professionKind ProfessionKind
	stats          Stats
	equipment      []*Equipment
	story          *Story // current story da dude be in
	room           *Room  // current room the dude is in
	stack          *render.Stack
	timer          int
	activity       DudeActivity
	activityDone   bool
	variation      float64
}

func NewDude(pk ProfessionKind, level int) *Dude {
	dude := &Dude{}

	stack, err := render.NewStack("dudes/liltest", "", "")
	if err != nil {
		panic(err)
	}
	stack.SetOriginToCenter()

	dude.name = assets.GetRandomName()
	dude.xp = 0
	dude.gold = 0
	dude.professionKind = pk
	profession := NewProfession(pk, level)
	dude.stats = profession.StartingStats()
	dude.equipment = profession.StartingEquipment()
	dude.variation = -6 + rand.Float64()*12

	dude.stack = stack

	return dude
}

func (d *Dude) Update(story *Story, req *ActivityRequests) {
	// NOTE: We should replace Centering/Moving direct position/rotation setting with a "pathing node" that the dude seeks to follow. This would allow more smoothly doing turns and such, as we could have a turn limit the dude would follow automatically...
	switch d.activity {
	case Idle:
		// Do nothing.
	case FirstEntering:
		cx, cy := d.Position()
		distance := story.DistanceFromCenter(cx, cy)
		if distance < 50+d.variation {
			d.activity = Centering
			d.stack.HeightOffset = 0
		} else {
			r := story.AngleFromCenter(cx, cy)
			nx, ny := story.PositionFromCenter(r, distance-d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny, cb: func(success bool) {
				d.stack.HeightOffset -= 0.15
				if d.stack.HeightOffset <= 0 {
					d.stack.HeightOffset = 0
				}
			}})
		}
	case GoingUp:
		d.timer++
		if d.stack.SliceOffset == 0 {
			d.stack.SliceOffset = d.stack.SliceCount()
			d.stack.MaxSliceIndex = 1
			cx, cy := d.Position()
			distance := story.DistanceFromCenter(cx, cy)
			r := story.AngleFromCenter(cx, cy)
			nx, ny := story.PositionFromCenter(r, distance+d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny})
		}
		if d.timer >= 15 {
			d.stack.SliceOffset--
			d.stack.MaxSliceIndex++
			d.timer = 0
		}
		if d.stack.SliceOffset <= 0 {
			d.stack.SliceOffset = 0
			d.stack.MaxSliceIndex = 0
			d.activity = Centering
		}
	case Centering:
		cx, cy := d.Position()
		distance := story.DistanceFromCenter(cx, cy)
		if distance >= 48+d.variation {
			d.activity = Moving
		} else {
			r := story.AngleFromCenter(cx, cy)
			nx, ny := story.PositionFromCenter(r, distance+d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny})
		}
	case Moving:
		cx, cy := d.Position()
		r := story.AngleFromCenter(cx, cy)
		nx, ny := story.PositionFromCenter(r-d.Speed(), 48+d.variation)

		face := math.Atan2(ny-cy, nx-cx)

		req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny})
	case Leaving:
		// TODO
	case GoingDown:
		// TODO
	}

	d.stack.Update()
}

func (d *Dude) Draw(o *render.Options) {
	d.stack.Draw(o)
}

func (d *Dude) Trigger(e Event) {
	switch e := e.(type) {
	case EventEnterRoom:
		d.room = e.room
	case EventLeaveRoom:
		d.room = nil
	}
	for _, eq := range d.equipment {
		eq.Activate(e)
	}
}

func (d *Dude) Position() (float64, float64) {
	return d.stack.Position()
}

func (d *Dude) SetPosition(x, y float64) {
	d.stack.SetPosition(x, y)
}

func (d *Dude) Rotation() float64 {
	return d.stack.Rotation()
}

func (d *Dude) SetRotation(r float64) {
	d.stack.SetRotation(r)
}

func (d *Dude) Room() *Room {
	return d.room
}

func (d *Dude) SetRoom(r *Room) {
	d.room = r
}

func (d *Dude) Name() string {
	return d.name
}

func (d *Dude) Level() int {
	return d.stats.Level()
}

// Scale speed with agility
// Thinkin we probably shouldn't calculate this like this...
func (d *Dude) Speed() float64 {
	// This values probably belong somewhere else
	speedScale := 0.1
	baseSpeed := 0.005
	return baseSpeed * (1 + float64(d.stats.Agility())*speedScale)
}

func (d *Dude) Stats() *Stats {
	return &d.stats
}

func (d *Dude) Profession() ProfessionKind {
	return d.professionKind
}

func (d *Dude) UpdateGold(gold float32) {
	d.gold += gold
}
