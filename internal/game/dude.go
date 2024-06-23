package game

import (
	"fmt"
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
	name         string
	xp           int
	gold         float32
	profession   ProfessionKind
	stats        Stats
	equipped     map[EquipmentType]*Equipment
	inventory    []*Equipment
	story        *Story // current story da dude be in
	room         *Room  // current room the dude is in
	stack        *render.Stack
	timer        int
	activity     DudeActivity
	activityDone bool
	variation    float64
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
	dude.profession = pk

	// Initialize stats and equipment
	profession := NewProfession(pk, level)
	dude.stats = profession.StartingStats()
	dude.inventory = profession.StartingEquipment()
	dude.AutoEquip() // equip starting equipment

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
			r := story.AngleFromCenter(cx, cy) + d.variation/5000
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
		if distance >= RoomPath+d.variation {
			d.activity = Moving
		} else {
			r := story.AngleFromCenter(cx, cy) + d.variation/5000
			nx, ny := story.PositionFromCenter(r, distance+d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny})
		}
	case Moving:
		cx, cy := d.Position()
		r := story.AngleFromCenter(cx, cy) + d.variation/5000
		nx, ny := story.PositionFromCenter(r-d.Speed(), RoomPath+d.variation)

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
		// Roll for loot on exit
		if eq := e.room.RollLoot(d.stats.luck); eq != nil {
			fmt.Println(d.name, "found", eq.Name())
			d.inventory = append(d.inventory, eq)
			d.AutoEquip() // Equip if possible
		}
		d.room = nil
	case EventEquip:
		fmt.Println(d.name, "equipped", e.equipment.Name())
	case EventUnequip:
		fmt.Println(d.name, "unequipped", e.equipment.Name())
	}

	// Trigger equipped equipment
	for _, eq := range d.equipped {
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
	return d.stats.level
}

// Scale speed with agility
// Thinkin we probably shouldn't calculate this like this...
func (d *Dude) Speed() float64 {
	// This values probably belong somewhere else
	speedScale := 0.1
	baseSpeed := 0.005
	return baseSpeed * (1 + float64(d.stats.agility)*speedScale)
}

func (d *Dude) Stats() *Stats {
	return &d.stats
}

func (d *Dude) Profession() ProfessionKind {
	return d.profession
}

func (d *Dude) Gold() float32 {
	return d.gold
}

func (d *Dude) UpdateGold(gold float32) {
	d.gold += gold
}

// Auto equips from inventory to empty slots
// If equipment can't be equipped by the class, it will not be equipped
// If slot is taken by another item, it will not be equipped and remain in inventory
func (d *Dude) AutoEquip() {
	for _, eq := range d.inventory {
		if eq.CanEquip(d.profession) {
			if _, ok := d.equipped[eq.Type()]; !ok {
				d.equip(eq)
			}
		}
	}
}

// Equips item to dude
func (d *Dude) equip(eq *Equipment) {
	if _, ok := d.equipped[eq.Type()]; ok {
		d.unequip(eq.Type())
		d.Trigger(EventUnequip{dude: d, equipment: eq}) // Event isolated to dude?

		d.equipped[eq.Type()] = eq
		d.Trigger(EventEquip{dude: d, equipment: eq}) // Event isolated to dude?

		// If equipment is in inventory, remove it
		for i, e := range d.inventory {
			if e == eq {
				d.inventory = append(d.inventory[:i], d.inventory[i+1:]...)
				break
			}
		}
	}
}

func (d *Dude) unequip(t EquipmentType) {
	if _, ok := d.equipped[t]; ok {
		// Add to inventory
		d.inventory = append(d.inventory, d.equipped[t])
		d.equipped[t] = nil
		d.Trigger(EventUnequip{dude: d, equipment: d.equipped[t]}) // Event isolated to dude?
	}
}

func (d *Dude) Inventory() []*Equipment {
	return d.inventory
}

func (d *Dude) Equipped() map[EquipmentType]*Equipment {
	return d.equipped
}

// Returns the stats of the dude with the equipment stats added
func (d *Dude) GetCalculatedStats() *Stats {
	stats := NewStats(nil)
	stats = stats.Add(&d.stats)
	for _, eq := range d.equipped {
		stats = stats.Add(eq.Stats())
	}
	return stats
}

func (d *Dude) AddXP(xp int) {
	d.xp += xp
	// If level reached
	nextLevelXP := 100 * d.Level()
	if d.xp >= nextLevelXP {
		d.xp -= nextLevelXP
		d.stats.LevelUp()
	}
}

func (d *Dude) XP() int {
	return d.xp
}

func (d *Dude) Heal(amount int) {
	if d.stats.currentHp+amount > d.stats.totalHp {
		d.stats.currentHp = d.stats.totalHp
	} else {
		d.stats.currentHp += amount
	}
}
