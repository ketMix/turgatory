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
	StairsToUp
	StairsFromDown
	GoingUp   // Entering the room from a staircase, this basically does the fancy slice offset/limiting.
	Centering // Move the dude to the center of the room.
	Moving    // Move the dude counter-clockwise.
	Leaving   // Move the dude to the stairs.
	GoingDown // Leaving the room to the stairs, opposite of GoingUp.
	EnterPortal
	Ded
	BornAgain
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
	shadow       *render.Stack
	timer        int
	activity     DudeActivity
	activityDone bool
	variation    float64
	enemy        *Enemy  // currently fighting enemy
	trueRotation float64 // This is the absolute rotation of the dude, ignoring facing.
}

func NewDude(pk ProfessionKind, level int) *Dude {
	dude := &Dude{}

	stack, err := render.NewStack("dudes/liltest", "", "")
	if err != nil {
		panic(err)
	}

	// Randomize which dude it be.
	stack.SetStack(stack.Stacks()[rand.Intn(len(stack.Stacks()))])
	stack.SetAnimation("base")

	// Get shadow.
	shadowStack, err := render.NewStack("dudes/shadow", "", "")
	if err != nil {
		panic(err)
	}
	dude.shadow = shadowStack

	// Assign a random dude skin
	stackNames := stack.Stacks()
	stack.SetStack(stackNames[rand.Intn(len(stackNames))])
	stack.SetOriginToCenter()

	dude.name = assets.GetRandomName()
	dude.xp = 0
	dude.gold = 0
	dude.profession = pk

	// Initialize stats and equipment
	profession := NewProfession(pk, level)
	dude.stats = profession.StartingStats()
	dude.inventory = make([]*Equipment, 0)

	dude.equipped = make(map[EquipmentType]*Equipment)
	for _, eq := range profession.StartingEquipment() {
		dude.Equip(eq)
	}

	dude.variation = -6 + rand.Float64()*12

	dude.stack = stack
	dude.stack.VgroupOffset = 1

	return dude
}

func (d *Dude) SetActivity(a DudeActivity) {
	d.activity = a
	d.timer = 0
	if a == Ded {
		d.stack.SetAnimation("ded")
	} else if a == BornAgain {
		d.stack.SetAnimation("base")
	}
}

func (d *Dude) Update(story *Story, req *ActivityRequests) {
	// NOTE: We should replace Centering/Moving direct position/rotation setting with a "pathing node" that the dude seeks to follow. This would allow more smoothly doing turns and such, as we could have a turn limit the dude would follow automatically...
	switch d.activity {
	case Idle:
		// Do nothing.
	case Ded:
		// Also do nothing!
	case BornAgain:
		d.SetActivity(Moving) // Is this safe to just set moving?
	case FirstEntering:
		cx, cy := d.Position()
		distance := story.DistanceFromCenter(cx, cy)
		if distance < 50+d.variation {
			d.SetActivity(Centering)
			d.stack.HeightOffset = 0
		} else {
			r := story.AngleFromCenter(cx, cy) + d.variation/5000
			nx, ny := story.PositionFromCenter(r, distance-d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)
			d.trueRotation = face

			req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
				d.stack.HeightOffset -= 0.15
				if d.stack.HeightOffset <= 0 {
					d.stack.HeightOffset = 0
				}
				d.SyncEquipment()
			}})
		}
	case StairsToUp:
		d.timer++

		if d.timer < 40 {
			cx, cy := d.Position()
			r := story.AngleFromCenter(cx, cy) + d.variation/5000
			nx, ny := story.PositionFromCenter(r-0.005, RoomPath+d.variation)

			d.stack.VgroupOffset = d.timer / 2

			face := math.Atan2(ny-cy, nx-cx)
			d.trueRotation = face

			req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
				if success {
					d.SyncEquipment()
				}
			}})
		} else {
			req.Add(StoryEnterNextActivity{dude: d, cb: func(success bool) {
				if success {
					d.SyncEquipment()
				}
			}})
			d.stack.VgroupOffset = 0
			d.SetActivity(StairsFromDown)
		}
	case StairsFromDown:
		d.timer++
		if d.stack.SliceOffset == 0 {
			d.stack.SliceOffset = d.stack.SliceCount()
			d.stack.MaxSliceIndex = 1
		}
		cx, cy := d.Position()
		r := story.AngleFromCenter(cx, cy) + d.variation/5000
		nx, ny := story.PositionFromCenter(r-0.01, RoomPath+d.variation)

		face := math.Atan2(ny-cy, nx-cx)
		d.trueRotation = face

		req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
			d.SyncEquipment()
		}})
		if d.timer >= 2 {
			d.stack.SliceOffset--
			d.stack.MaxSliceIndex++
			d.timer = 0
		}
		if d.stack.SliceOffset <= 0 {
			d.stack.SliceOffset = 0
			d.stack.MaxSliceIndex = 0
			d.SetActivity(Moving)
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
			d.trueRotation = face

			req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
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
			d.SetActivity(Centering)
		}
	case Centering:
		cx, cy := d.Position()
		distance := story.DistanceFromCenter(cx, cy)
		if distance >= RoomPath+d.variation {
			d.SetActivity(Moving)
		} else {
			r := story.AngleFromCenter(cx, cy) + d.variation/5000
			nx, ny := story.PositionFromCenter(r, distance+d.Speed()*100)

			face := math.Atan2(ny-cy, nx-cx)
			d.trueRotation = face

			req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
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
		d.trueRotation = face

		// Face inwards if we have an enemy!
		if d.enemy != nil {
			fx, fy := story.PositionFromCenter(r-d.Speed(), d.variation)
			face = math.Atan2(fy-cy, fx-cx)
		}

		req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
			if success {
				d.SyncEquipment()
			}
		}})
	case Leaving:
		// TODO
	case GoingDown:
		// TODO
	case EnterPortal:
		d.timer++
		// Wait a little bit before entering!
		if d.timer >= 30 {
			cx, cy := d.Position()
			distance := story.DistanceFromCenter(cx, cy)

			d.stack.Transparency = float32(d.timer-30) / 20
			d.shadow.Transparency = float32(d.timer-30) / 20

			if distance < PortalDistance-4+d.variation {
				d.stack.Transparency = 1
				d.shadow.Transparency = 1
				d.SetActivity(Idle)
				req.Add(TowerLeaveActivity{dude: d})
			} else {
				r := story.AngleFromCenter(cx, cy)
				nx, ny := story.PositionFromCenter(r, distance-0.005*100)

				face := math.Atan2(ny-cy, nx-cx)
				d.trueRotation = face

				req.Add(MoveActivity{dude: d, face: face, x: nx, y: ny, cb: func(success bool) {
					d.SyncEquipment()
				}})
			}
		}

	}

	d.stack.Update()

	// Update equipment
	for _, eq := range d.equipped {
		if eq != nil {
			eq.Update()
		}
	}

	// Update enemy if there is one
	if d.enemy != nil {
		d.enemy.Update(d)
	}
}

func (d *Dude) SyncEquipment() {
	// Piggy-back syncing shadow here
	d.shadow.SetOrigin(d.stack.Origin())
	d.shadow.SetPosition(d.stack.Position())
	d.shadow.SetRotation(d.stack.Rotation())
	for _, eq := range d.equipped {
		if eq != nil && eq.stack != nil {
			// Set equipment position to dude position
			eq.stack.SliceOffset = d.stack.SliceOffset
			eq.stack.MaxSliceIndex = d.stack.MaxSliceIndex
			eq.stack.HeightOffset = d.stack.HeightOffset
			eq.stack.VgroupOffset = d.stack.VgroupOffset
			eq.stack.Transparency = d.stack.Transparency
			eq.stack.SetOrigin(d.stack.Origin())
			eq.stack.SetPosition(d.stack.Position())
			eq.stack.SetRotation(d.stack.Rotation())
		}
	}
}

func (d *Dude) Draw(o *render.Options) {
	d.stack.Draw(o)
	d.shadow.Draw(o)

	if d.IsDead() {
		return
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

func (d *Dude) Trigger(e Event) Activity {
	// Trigger equipped equipment
	// It may modify event amounts
	for _, eq := range d.equipped {
		if eq != nil {
			eq.Activate(e)
		}
	}

	switch e := e.(type) {
	case EventDudeHit:
		if d.IsDead() {
			return DudeDeadActivity{dude: d}
		}
	case EventCombatRoom:
		// Attack enemy if there is one
		if d.enemy != nil {
			damage, isCrit := d.GetDamage()
			if damage == 0 {
				d.Trigger(EventDudeMiss{dude: d, enemy: d.enemy})
				AddMessage(
					MessageCombat,
					fmt.Sprintf("%s missed their attack against %s!", d.name, d.enemy.name),
				)
			} else if isCrit {
				AddMessage(
					MessageCombat,
					fmt.Sprintf("%s crit %s for %d damage!", d.name, d.enemy.name, damage),
				)
				d.Trigger(EventDudeCrit{dude: d, enemy: d.enemy, amount: damage})
			}
			isDead := d.enemy.Damage(d.stats.strength)

			if isDead {
				xp := d.enemy.XP()
				gold := d.enemy.Gold()
				d.UpdateGold(gold)
				d.AddXP(xp)
				AddMessage(
					MessageCombat,
					fmt.Sprintf("%s defeated %s and gained %d xp and %.0fgp", d.name, d.enemy.name, xp, gold),
				)
				d.enemy = nil
			} else {
				takenDamage, isDodge := d.ApplyDamage(d.enemy.Hit())
				if !isDodge {
					if act := d.Trigger(EventDudeHit{dude: d, amount: takenDamage}); act != nil {
						return act
					}
				} else {
					d.Trigger(EventDudeDodge{dude: d, enemy: d.enemy})
					AddMessage(
						MessageCombat,
						fmt.Sprintf("%s dodged an attack from %s", d.name, d.enemy.name),
					)
				}
				if d.IsDead() {
					d.enemy = nil
				}
			}
		}
		// Else it may be a trap room
		if act := d.room.GetRoomEffect(e); act != nil {
			return act
		}
	case EventEnterRoom:
		if act := d.room.GetRoomEffect(e); act != nil {
			return act
		}
	case EventCenterRoom:
		if act := d.room.GetRoomEffect(e); act != nil {
			return act
		}
	case EventLeaveRoom:
		if act := d.room.GetRoomEffect(e); act != nil {
			return act
		}
	case EventEquip:
		//fmt.Println(d.name, "equipped", e.equipment.Name())
		if d.stack != nil {
			t := MakeFloatingTextFromDude(d, fmt.Sprintf("equip %s", e.equipment.Name()), color.NRGBA{100, 200, 200, 255}, 120, 0.4)
			AddMessage(
				MessageInfo,
				fmt.Sprintf("%s equipped %s", d.name, e.equipment.Name()),
			)
			d.story.AddText(t)
		}
	case EventUnequip:
		//fmt.Println(d.name, "unequipped", e.equipment.Name())
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("remove %s", e.equipment.Name()), color.NRGBA{200, 100, 100, 255}, 120, 0.4)
		d.story.AddText(t)
	case EventGoldGain:
		//fmt.Println(d.name, "gained", e.amount, "gold")
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("+%.0fgp", e.amount), color.NRGBA{255, 255, 0, 255}, 40, 0.6)
		d.story.AddText(t)
	case EventGoldLoss:
		//fmt.Println(d.name, "lost", e.amount, "gold")
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("-%.0fgp", e.amount), color.NRGBA{255, 255, 0, 255}, 40, 0.4)
		d.story.AddText(t)
	}
	return nil
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
	baseSpeed := 0.01
	// Slow dude down when in combat.
	if d.enemy != nil {
		return baseSpeed * (1 + speedScale)
	}
	stats := d.GetCalculatedStats()
	return baseSpeed * (1 + float64(stats.agility/10)*speedScale)
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
		d.AddXP(1)
		t := MakeFloatingTextFromDude(d, "*CRIT*", color.NRGBA{255, 128, 255, 128}, 60, 1.0)
		d.story.AddText(t)
		multiplier = 2.0
		wasCrit = true
	}

	// Cowardice can cause miss
	// capped at 50% miss chance
	// 1 cowardice = 1% miss chance
	// 100 cowardice = 50% miss chance
	missRoll := math.Min(float64(stats.cowardice)/100.0, 0.5)

	if rand.Float64() < missRoll {
		t := MakeFloatingTextFromDude(d, "*miss*", color.NRGBA{128, 128, 128, 128}, 30, 0.5)
		d.story.AddText(t)
		multiplier = 0.0
	}

	return int(float64(stats.strength) * multiplier), wasCrit
}

func (d *Dude) ApplyDamage(amount int) (int, bool) {
	// Luck and agility can cause dodge
	stats := d.GetCalculatedStats()
	baseChance := 0.05                                   // 5% base dodge chance
	luckContribution := float64(stats.luck) * 0.005      // 0.5% per luck point
	agilityContribution := float64(stats.agility) * 0.01 // 1% per agility point

	dodgeChance := baseChance + luckContribution + agilityContribution

	// Cap the maximum dodge chance at 50%
	chance := math.Min(dodgeChance, 0.5)
	if rand.Float64() < chance {
		d.AddXP(1)
		t := MakeFloatingTextFromDude(d, "*dodge*", color.NRGBA{255, 255, 0, 128}, 30, 0.5)
		d.story.AddText(t)
		return 0, true
	}

	// Apply defense stat
	// Higher defense means less damage taken with diminishing returns
	// 1 defense = 1% damage reduction
	// 100 defense = 50% damage reduction
	amount = int(float64(amount) * (1.0 - float64(stats.defense)/200.0))
	if amount < 0 {
		amount = 1
	}

	d.stats.currentHp -= amount

	if d.stats.currentHp < 0 {
		d.stats.currentHp = 0
	}

	if d.stats.currentHp == 0 {
		d.SetActivity(Ded)
		t := MakeFloatingTextFromDude(d, "RIP", color.NRGBA{64, 64, 64, 255}, 80, 1)
		d.story.AddText(t)
	} else {
		//fmt.Println(d.name, "took", amount, "damage and has", d.stats.currentHp, "HP left")
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("%d", -amount), color.NRGBA{255, 0, 0, 255}, 40, 0.5)
		d.story.AddText(t)
	}
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

	//fmt.Println(d.name, "added", eq.Name(), "to inventory")
	t := MakeFloatingTextFromDude(d, fmt.Sprintf("+%s", eq.Name()), color.NRGBA{200, 200, 50, 128}, 100, 1.0)
	d.story.AddText(t)
}

func (d *Dude) Inventory() []*Equipment {
	return d.inventory
}

func (d *Dude) Equipped() map[EquipmentType]*Equipment {
	return d.equipped
}

// Returns the stats of the dude with the equipment stats added
func (d *Dude) GetCalculatedStats() *Stats {
	stats := NewStats(nil, false)
	stats = stats.Add(&d.stats)
	for _, eq := range d.equipped {
		stats = stats.Add(eq.Stats())
	}
	stats.currentHp = d.stats.currentHp
	return stats
}

func (d *Dude) AddXP(xp int) {
	d.xp += xp
	// If level reached
	nextLevelXP := d.NextLevelXP()
	if d.xp >= nextLevelXP {
		d.xp -= nextLevelXP
		d.stats.LevelUp()
		t := MakeFloatingTextFromDude(d, "LEVEL UP", color.NRGBA{100, 255, 255, 255}, 80, 1)
		d.story.AddText(t)
	} else {
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("+%dxp", xp), color.NRGBA{100, 200, 200, 200}, 50, 1)
		d.story.AddText(t)
	}
}

func (d *Dude) XP() int {
	return d.xp
}
func (d *Dude) NextLevelXP() int {
	return 50 * d.Level()
}

func (d *Dude) Heal(amount int) {
	initialHP := d.stats.currentHp
	stats := d.GetCalculatedStats()
	d.stats.currentHp += amount
	if d.stats.currentHp > stats.totalHp {
		d.stats.currentHp = stats.totalHp
	}
	amount = d.stats.currentHp - initialHP
	//fmt.Println(d.name, "healed", amount, "HP", " and went from ", initialHP, " to ", d.stats.currentHp)
	if d.story != nil {
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("+%d", amount), color.NRGBA{0, 255, 0, 255}, 40, 0.5)
		d.story.AddText(t)
	}
}

func (d *Dude) FullHeal() {
	stats := d.GetCalculatedStats()

	// No rez
	if d.stats.currentHp >= 0 {
		d.Heal(stats.totalHp)
	}
}

func (d *Dude) RestoreUses(amount int) {
	for _, eq := range d.equipped {
		if eq != nil {
			eq.RestoreUses(amount)
		}
	}
	//fmt.Println(d.name, "restored equipment uses by", amount)
	t := MakeFloatingTextFromDude(d, fmt.Sprintf("+eq restore %d", amount), color.NRGBA{0, 128, 255, 200}, 40, 0.5)
	d.story.AddText(t)
}

func (d *Dude) RandomEquippedItem() *Equipment {
	equippedTypes := []EquipmentType{}
	for t, eq := range d.equipped {
		if eq != nil {
			equippedTypes = append(equippedTypes, t)
		}
	}

	if len(equippedTypes) == 0 {
		return nil
	}

	et := equippedTypes[rand.Intn(len(equippedTypes))]
	return d.equipped[et]
}

func (d *Dude) LevelUpEquipment(amount int) {
	// Random equipped item
	eq := d.RandomEquippedItem()
	if eq == nil {
		return
	}

	for i := 0; i < amount; i++ {
		eq.LevelUp()
	}
	//fmt.Println(d.name, "leveled up equipment by", amount)
	t := MakeFloatingTextFromDude(d, fmt.Sprintf("+eq up %d", amount), color.NRGBA{128, 128, 255, 255}, 50, 0.5)
	d.story.AddText(t)
}

func (d *Dude) Perkify(maxQuality PerkQuality) {
	eq := d.RandomEquippedItem()
	if eq == nil {
		return
	}

	// Assign random perk
	if eq.perk == nil {
		prevName := eq.Name()
		eq.perk = GetRandomPerk(PerkQualityTrash)
		//fmt.Println(d.name, "upgraded his equipment", prevName, "with", eq.perk.Name())
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("+%s perk %s", prevName, eq.perk.Name()), color.NRGBA{128, 255, 128, 255}, 100, 0.5)
		d.story.AddText(t)
	} else {
		// Level up perk
		previousQuality := eq.perk.Quality()
		previousName := eq.Name()
		eq.perk.LevelUp(maxQuality)
		if eq.perk.Quality() != previousQuality {
			//fmt.Println(eq.perk.Quality(), previousQuality)
			//fmt.Println(d.name, "upgraded his equipment", previousName, "to", eq.Name())
			t := MakeFloatingTextFromDude(d, fmt.Sprintf("+eq %s upgrade to %s", previousName, eq.Name()), color.NRGBA{128, 255, 128, 255}, 100, 0.5)
			d.story.AddText(t)
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
	stats := d.GetCalculatedStats()

	// Ensure stats are not negative
	wis := max(stats.wisdom, 1)
	luck := max(stats.luck, 1)

	wisdomRoll := rand.Intn(wis) + 1 // ensure non-zero roll
	luckRoll := rand.Intn(luck) + 1
	highestRoll := max(wisdomRoll, luckRoll)

	// Higher the roll, lower the chance of being cursed
	threshold := 1.0 - math.Log10(float64(highestRoll+1))

	curseRoll := rand.Float64()
	if curseRoll > threshold {
		// Spared
		return
	}

	// Check for gold loss
	if curseRoll <= threshold*0.75 { // high chance for gold loss
		goldLoss := float32(roomLevel * 10)
		d.UpdateGold(-goldLoss)
		//fmt.Println(d.name, "lost", goldLoss, "gold")
		t := MakeFloatingTextFromDude(d, fmt.Sprintf("-%.0fgp", goldLoss), color.NRGBA{255, 255, 0, 200}, 40, 0.5)
		d.story.AddText(t)
	}

	// Check for equipment delevel
	if curseRoll <= threshold*0.5 { // reduced chance for equipment delevel
		equipmentType := RandomEquipmentType()
		if eq := d.equipped[equipmentType]; eq != nil {
			eq.LevelDown()
			//fmt.Println(d.name, "lost a level on", eq.Name())
			t := MakeFloatingTextFromDude(d, fmt.Sprintf("-eq level %s", eq.Name()), color.NRGBA{200, 200, 32, 200}, 50, 0.5)
			d.story.AddText(t)
		}
	}

	// Check for perk delevel
	if curseRoll <= threshold*0.25 { // even lower chance for perk delevel
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
				//fmt.Println(d.name, "lost a perk level on", eq.Name())
				t := MakeFloatingTextFromDude(d, fmt.Sprintf("-eq perk %s", eq.Name()), color.NRGBA{200, 200, 32, 200}, 50, 0.5)
				d.story.AddText(t)
			}
		}
	}

	// Check for dude delevel
	if curseRoll <= threshold*0.1 { // lowest chance for dude delevel
		d.stats.LevelDown()
		//fmt.Println(d.name, "lost a level")
		t := MakeFloatingTextFromDude(d, "-level", color.NRGBA{100, 0, 0, 200}, 50, 0.5)
		d.story.AddText(t)
	}
}

func (d *Dude) TrapDamage(roomLevel int) {
	// Chance based on agility
	agilityRoll := rand.Intn(d.stats.agility + 1)

	// Higher agility, lower chance of being hit
	threshold := 1.0 - math.Log10(float64(agilityRoll+1))
	trapRoll := rand.Float64()

	if trapRoll > threshold {
		return
	}

	// Damage based on room level
	damage := roomLevel * 2

	amount, miss := d.ApplyDamage(damage)
	if !miss {
		AddMessage(
			MessageCombat,
			fmt.Sprintf("%s took %d damage from a trap", d.name, amount),
		)
	}
}

func (d *Dude) IsDead() bool {
	return d.stats.currentHp <= 0
}
