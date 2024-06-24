package game

import (
	"fmt"
	"image/color"
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
	enemy        *Enemy // currently fighting enemy
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
	fmt.Println(dude.name, "the", pk.String(), "has been created with stats", dude.stats)
	dude.inventory = profession.StartingEquipment()
	dude.equipped = make(map[EquipmentType]*Equipment)
	for _, eq := range dude.inventory {
		dude.Equip(eq)
	}

	dude.variation = -6 + rand.Float64()*12

	dude.stack = stack
	dude.stack.VgroupOffset = 1

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
				d.SyncEquipment()
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

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny, cb: func(success bool) {
				d.SyncEquipment()
			}})
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

			req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny, cb: func(success bool) {
				if success {
					d.SyncEquipment()
				}
			}})
		}
	case Moving:
		cx, cy := d.Position()
		r := story.AngleFromCenter(cx, cy) + d.variation/5000
		nx, ny := story.PositionFromCenter(r-d.Speed(), RoomPath+d.variation)

		face := math.Atan2(ny-cy, nx-cx)

		req.Add(MoveActivity{initiator: d, face: face, x: nx, y: ny, cb: func(success bool) {
			if success {
				d.SyncEquipment()
			}
		}})
	case Leaving:
		// TODO
	case GoingDown:
		// TODO
	}

	d.stack.Update()

	// Update equipment
	for _, eq := range d.equipped {
		if eq != nil {
			eq.Update()
		}
	}
}

func (d *Dude) SyncEquipment() {
	for _, eq := range d.equipped {
		if eq != nil && eq.stack != nil {
			// Set equipment position to dude position
			eq.stack.SliceOffset = d.stack.SliceOffset
			eq.stack.MaxSliceIndex = d.stack.MaxSliceIndex
			eq.stack.HeightOffset = d.stack.HeightOffset
			eq.stack.VgroupOffset = d.stack.VgroupOffset
			eq.stack.SetOrigin(d.stack.Origin())
			eq.stack.SetPosition(d.stack.Position())
			eq.stack.SetRotation(d.stack.Rotation())
		}
	}
}

func (d *Dude) Draw(o *render.Options) {
	d.stack.Draw(o)

	// THIS IS BORKED
	if len(o.VGroup.Images) > 0 {
		ro := render.TextOptions{
			Font:   assets.BodyFont,
			Color:  color.NRGBA{255, 0, 255, 255},
			Screen: o.VGroup.Overlay,
		}

		x, y := d.stack.Position()
		y /= o.Camera.Zoom

		ro.GeoM.Translate(x, y)

		render.DrawText(&ro, d.name)
	}
	// Draw equipment
	for _, eq := range d.equipped {
		if eq != nil {
			eq.Draw(o)
		}
	}
	// Reset colors, as equipment may have munged it.
	o.DrawImageOptions.ColorScale.Reset()

	// Draw enemy if there is one
	if d.enemy != nil {
		d.enemy.Draw(*o)
	}
}

func (d *Dude) DrawProfile(o *render.Options) {
	stack := render.CopyStack(d.stack)
	stack.SetPosition(0, 0)
	stack.SetOrigin(0, 0)
	stack.SetRotation(-math.Pi / 2)
	stack.Draw(o)

	// Draw armor (like helmet or soemthing) ?
	armor := d.equipped[EquipmentTypeArmor]
	if armor != nil && armor.stack != nil {
		stack = render.CopyStack(armor.stack)
		stack.SetPosition(0, 0)
		stack.SetOrigin(0, 0)
		stack.SetRotation(-math.Pi / 2)
		stack.Draw(o)
	}
}

func (d *Dude) Trigger(e Event) {
	// Trigger equipped equipment
	// It may modify event amounts
	for _, eq := range d.equipped {
		if eq != nil {
			eq.Activate(e)
		}
	}

	switch e := e.(type) {
	case EventGlobalTick:
		// Attack enemy if there is one
		if d.enemy != nil {
			damage, isCrit := d.GetDamage()
			if damage == 0 {
				d.Trigger(EventDudeMiss{dude: d, enemy: d.enemy})
			} else if isCrit {
				d.Trigger(EventDudeCrit{dude: d, enemy: d.enemy, amount: damage})
			}
			isDead := d.enemy.Damage(d.stats.strength)

			if isDead {
				xp := d.enemy.XP()
				gold := d.enemy.Gold()
				d.UpdateGold(gold)
				d.AddXP(xp)
				d.enemy = nil
			} else {
				takenDamage, isDodge := d.ApplyDamage(d.enemy.Hit())
				if !isDodge {
					d.Trigger(EventDudeHit{dude: d, amount: takenDamage})
				}
			}
		}
	case EventEnterRoom:
		d.room = e.room
		d.room.GetRoomEffect(e)
	case EventCenterRoom:
		d.room.GetRoomEffect(e)
	case EventLeaveRoom:
		d.room.GetRoomEffect(e)
		d.room = nil
	case EventEquip:
		fmt.Println(d.name, "equipped", e.equipment.Name())
	case EventUnequip:
		fmt.Println(d.name, "unequipped", e.equipment.Name())
	case EventGoldGain:
		fmt.Println(d.name, "gained", e.amount, "gold")
	case EventGoldLoss:
		fmt.Println(d.name, "lost", e.amount, "gold")
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

// TODO: Refine this
func (d *Dude) GetDamage() (int, bool) {
	wasCrit := false
	stats := d.GetCalculatedStats()

	// Luck can cause crit
	luckRoll := float64(stats.luck+1) * 0.1
	randRoll := rand.Float64()

	multiplier := 1.0
	if randRoll < luckRoll {
		fmt.Println(d.name, "crit!")
		multiplier = 2.0
		wasCrit = true
	}

	// Cowardice can cause miss
	missRoll := float64(stats.cowardice+1) * 0.01
	if rand.Float64() < missRoll {
		fmt.Println(d.name, "missed!")
		multiplier = 0.0
	}

	return int(float64(stats.strength) * multiplier), wasCrit
}

func (d *Dude) ApplyDamage(amount int) (int, bool) {
	// Luck and agility can cause dodge
	luckRoll := float64(d.stats.luck+1) * 0.1
	agilityRoll := float64(d.stats.agility+1) * 0.1
	if rand.Float64() < luckRoll || rand.Float64() < agilityRoll {
		fmt.Println(d.name, "dodged!")
		return 0, true
	}

	// Apply defense stat
	amount -= d.stats.defense
	if amount < 0 {
		amount = 0
	}

	d.stats.currentHp -= amount

	if d.stats.currentHp < 0 {
		d.stats.currentHp = 0
	}

	fmt.Println(d.name, "took", amount, "damage and has", d.stats.currentHp, "HP left")
	return amount, false

	// If dead, uh, do something right? maybe an event or something idk
	// maybe just doomed to roam the story forever
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

func (d *Dude) UpdateGold(amount float32) {
	d.gold += amount
	if d.gold < 0 {
		d.gold = 0
	}

	if amount > 0 {
		d.Trigger(EventGoldGain{dude: d, amount: amount})
	} else {
		d.Trigger(EventGoldLoss{dude: d, amount: amount})
	}
}

// Equips item to dude
func (d *Dude) Equip(eq *Equipment) {
	if _, ok := d.equipped[eq.Type()]; ok {
		d.Unequip(eq.Type())
		d.Trigger(EventUnequip{dude: d, equipment: eq}) // Event isolated to dude?
	}

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

func (d *Dude) Unequip(t EquipmentType) {
	if _, ok := d.equipped[t]; ok {
		// Add to inventory
		d.inventory = append(d.inventory, d.equipped[t])
		d.Trigger(EventUnequip{dude: d, equipment: d.equipped[t]}) // Event isolated to dude?
		d.equipped[t] = nil
	}
}

func (d *Dude) AddToInventory(eq *Equipment) {
	d.inventory = append(d.inventory, eq)
	if d.equipped[eq.Type()] == nil {
		d.Equip(eq)
	}

	fmt.Println(d.name, "added", eq.Name(), "to inventory")
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
	initialHP := d.stats.currentHp
	if d.stats.currentHp+amount > d.stats.totalHp {
		d.stats.currentHp = d.stats.totalHp
	} else {
		d.stats.currentHp += amount
	}
	fmt.Println(d.name, "healed", amount, "HP", " and went from ", initialHP, " to ", d.stats.currentHp)
}

func (d *Dude) RestoreUses(amount int) {
	for _, eq := range d.equipped {
		if eq != nil {
			eq.RestoreUses(amount)
		}
	}
	fmt.Println(d.name, "restored equipment uses by", amount)
}

func (d *Dude) LevelUpEquipment(amount int) {
	for _, eq := range d.equipped {
		if eq != nil {
			for i := 0; i < amount; i++ {
				eq.LevelUp()
			}
		}
	}
	fmt.Println(d.name, "leveled up equipment by", amount)
}

func (d *Dude) Perkify(maxQuality PerkQuality) {
	// Random equipped item
	equipmentType := RandomEquipmentType()

	if eq := d.equipped[equipmentType]; eq != nil {
		// Assign random perk
		if eq.perk == nil {
			prevName := eq.Name()
			eq.perk = GetRandomPerk(PerkQualityTrash)
			fmt.Println(d.name, "upgraded his equipment", prevName, "with", eq.perk.Name())
		} else {
			// Level up perk
			previousQuality := eq.perk.Quality()
			previousName := eq.Name()
			eq.perk.LevelUp(maxQuality)
			if eq.perk.Quality() != previousQuality {
				fmt.Println(eq.perk.Quality(), previousQuality)
				fmt.Println(d.name, "upgraded his equipment", previousName, "to", eq.Name())
			}
		}
	}
}

// Cursify the dude
// Rolls twice, once with wisdom, the other with luck, and takes the highest
// Has a chance to
// - Delevel equipment (high chance)
// - Delevel perk (medium chance)
// - Delevel dude (low chance)
func (d *Dude) Cursify(roomLevel int) {
	wisdomRoll := rand.Intn(d.stats.wisdom)
	luckRoll := rand.Intn(d.stats.luck)
	threshold := 1 - 1/math.Max(float64(wisdomRoll), float64(luckRoll))

	curseRoll := rand.Float64()
	if curseRoll > threshold {
		// Spared
		return
	}

	if curseRoll > threshold/2 {
		// Roll for equipment type
		equipmentType := RandomEquipmentType()
		if eq := d.equipped[equipmentType]; eq != nil {
			eq.LevelDown()
		}
	}

	// Roll for perk
	if curseRoll > threshold/4 {
		// Random equipped item (that has a perk)
		equipmentWithPerks := []EquipmentType{}

		for t, eq := range d.equipped {
			if eq != nil && eq.perk != nil {
				equipmentWithPerks = append(equipmentWithPerks, t)
			}
		}

		if len(equipmentWithPerks) > 0 {
			randomEquipType := equipmentWithPerks[rand.Intn(len(equipmentWithPerks))]
			if eq := d.equipped[randomEquipType]; eq != nil {
				eq.perk.LevelDown()
			}
		}
	}

	// Roll for dude
	if curseRoll <= threshold/4 {
		d.stats.LevelDown()
	}
}
