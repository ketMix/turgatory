package game

import (
	"math"
	"math/rand"

	"github.com/kettek/ebijam24/internal/render"
)

type DudeActivity int

const (
	Idle      DudeActivity = iota
	GoingUp                // Entering the room from a staircase, this basically does the fancy slice offset/limiting.
	Centering              // Move the dude to the center of the room.
	Moving                 // Move the dude counter-clockwise.
	Leaving                // Move the dude to the stairs.
	GoingDown              // Leaving the room to the stairs, opposite of GoingUp.
)

type Dude struct {
	stack        *render.Stack
	speed        float64
	activity     DudeActivity
	activityDone bool
	variation    float64
}

func NewDude() *Dude {
	dude := &Dude{}

	stack, err := render.NewStack("dudes/liltest", "", "")
	if err != nil {
		panic(err)
	}
	stack.SetOriginToCenter()

	dude.speed = 0.002
	dude.activity = Centering
	dude.variation = -3 + rand.Float64()*6

	dude.stack = stack

	return dude
}

func (d *Dude) Update(story *Story) {
	switch d.activity {
	case Idle:
		// Do nothing.
	case GoingUp:
		// TODO
	case Centering:
		cx, cy := d.Position()
		distance := story.DistanceFromCenter(cx, cy)
		if distance >= 48+d.variation {
			d.activity = Moving
		} else {
			r := story.AngleFromCenter(cx, cy)
			nx, ny := story.PositionFromCenter(r, distance+d.speed*100)

			face := math.Atan2(ny-cy, nx-cx)

			d.SetRotation(face)
			d.SetPosition(nx, ny)
		}
	case Moving:
		cx, cy := d.Position()
		r := story.AngleFromCenter(cx, cy)
		nx, ny := story.PositionFromCenter(r-d.speed, 48+d.variation)

		face := math.Atan2(ny-cy, nx-cx)

		d.SetRotation(face)
		d.SetPosition(nx, ny)
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
