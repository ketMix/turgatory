package game

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/kettek/ebijam24/internal/render"
)

// RoomSize represents the different sizes of rooms, with 1 equating to 1/8th of a circle.
type RoomSize int

func (r RoomSize) String() string {
	switch r {
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Large:
		return "large"
	case Huge:
		return "huge"
	default:
		return "Unknown"
	}
}

const (
	Small  RoomSize = 1
	Medium RoomSize = 2
	Large  RoomSize = 3
	Huge   RoomSize = 4
)

// These origins are used to re-position a room "pie" image so that its center is in the appropriate place.
const (
	LargeOriginY = 44
	HugeOriginY  = 64
)

const CombatTickrate = 30 // Tick combat every 30 ticks

// RoomStairsEntrance is the distance from the center that a room's stairs is expected to be at.
const RoomStairsEntrance = 12
const RoomPath = 53
const TowerStairs = 60
const TowerEntrance = 80

// RoomKind is an enumeration of the different kinds of rooms in za toweru.
type RoomKind int

func (r *RoomKind) String() string {
	switch *r {
	case Stairs:
		return "stairs"
	case Armory:
		return "armory"
	case HealingShrine:
		return "healing"
	case Treasure:
		return "treasure"
	case Curse:
		return "curse"
	case Combat:
		return "combat"
	case Well:
		return "well"
	case Library:
		return "library"
	case Empty:
		return "template"
	case Trap:
		return "trap"
	case Boss:
		return "boss"
	default:
		return "Unknown"
	}
}

func (r *RoomKind) GetRoomEnemy(roomSize RoomSize, storyLevel int) EnemyKind {
	switch *r {
	case Combat:
		switch roomSize {
		case Small:
			return EnemyRat
		case Medium:
			return EnemySlime
		case Large:
			return EnemySkelly
		}
	case Boss:
		// Every 3 levels, boss is upgraded
		if storyLevel <= 3 {
			return EnemyBossRat
		} else if storyLevel <= 6 {
			return EnemyBossRat
		} else if storyLevel <= 9 {
			return EnemyBossRat
		} else if storyLevel <= 12 {
			// Last level?
			return EnemyBossRat
		}
	}
	return EnemyUnknown
}

// Equipment you can find in room
// hmm, kinda defeats the current structure of equipment in yaml if we have to add new equipment here
func (r *RoomKind) Equipment() []*string {
	switch *r {
	case Armory:
		types := []EquipmentType{EquipmentTypeWeapon, EquipmentTypeArmor}
		return GetEquipmentNamesWithTypes(types)
	case Treasure:
		types := []EquipmentType{EquipmentTypeAccessory}
		return GetEquipmentNamesWithTypes(types)
	}
	return nil
}

const (
	Empty RoomKind = iota
	//
	Stairs
	// Armory provide... armor up? damage up? Maybe should be different types.
	Armory
	// Healing shrine heals the adventurers over time.
	HealingShrine
	// Combat is where it goes down. $$$ is acquired.
	Combat
	// Well restores magic items?
	Well
	// Treasure room - $$$
	Treasure
	// Library - enchant
	Library
	// Curse - % to curse dude
	Curse
	// Trap - % to damage dude based on stats
	Trap

	// Boss room
	Boss

	// Marker for the end... allows for iteration
	RoomKindEnd
)

type RoomTemplate struct {
	kind RoomKind
	size RoomSize
}

// Room is a room within a story of za toweru.
type Room struct {
	story          *Story
	index          int // Reference to the index within the story.
	kind           RoomKind
	size           RoomSize
	power          int // ???
	combatTicks    int
	boss           *Enemy
	killedBoss     bool
	stacks         render.Stacks
	walls          render.Stacks
	dudesInCenter  []*Dude
	dudesInWaiting []*Dude
	dudes          []*Dude
	highlight      bool
}

func NewRoom(size RoomSize, kind RoomKind) *Room {
	r := &Room{
		size:       size,
		kind:       kind,
		killedBoss: false,
	}

	stack, err := render.NewStack(fmt.Sprintf("rooms/%s", size.String()), kind.String(), "")
	if err != nil {
		panic(err)
	}
	r.stacks.Add(stack)

	// Add our walls.
	for j := 0; j < 3; j++ {
		for i := 0; i < 8; i++ {
			stack, err := render.NewStack(fmt.Sprintf("walls/%s", size.String()), kind.String(), "base")
			if err != nil {
				stack, err = render.NewStack(fmt.Sprintf("walls/%s", size.String()), "template", "base")
				if err != nil {
					continue
				}
			}
			stack.VgroupOffset = j * StoryWallHeight
			r.walls.Add(stack)
		}
	}

	return r
}

// Update updates the stuff in the room.
func (r *Room) Update(req *ActivityRequests) {
	r.stacks.Update()
	if r.boss != nil {
		r.boss.RoomUpdate(r)
	}
	if r.kind == Combat || r.kind == Trap {
		r.combatTicks++
		if r.combatTicks >= CombatTickrate {
			r.combatTicks = 0
			for _, d := range r.dudes {
				req.Add(RoomCombatActivity{room: r, dude: d})
			}
		}
	}
	if r.kind == Boss {
		if r.boss != nil {
			if r.boss.IsDead() {
				goldPerDude := r.boss.Gold() / float64(len(r.dudes))
				xp := r.boss.XP() / len(r.dudes)
				r.boss = nil
				r.killedBoss = true
				for _, d := range r.dudes {
					d.AddXP(xp)
					d.UpdateGold(goldPerDude)
					req.Add(RoomEndBossActivity{room: r, dude: d})
				}
			} else {
				// Boss combat
				r.combatTicks++
				if r.combatTicks >= CombatTickrate {
					r.combatTicks = 0
					bossTarget := r.boss.GetTarget(r.dudes)
					if bossTarget != nil {
						bossTarget.ApplyDamage(r.boss.Hit())
						bossTarget.stats.ModifyStat(StatCowardice, r.boss.stats.strength)
					}
					for _, d := range r.dudes {
						if !d.IsDead() && !r.boss.IsDead() {
							dmg, _ := d.GetDamage()
							if dmg > 0 {
								r.boss.Damage(dmg)
							}
						}
					}
				}
			}

		} else {
			// If all dudes are waiting, trigger boss fight
			aliveDudes := 0
			for _, d := range r.story.dudes {
				if !d.IsDead() {
					aliveDudes++
				}
			}
			if aliveDudes > 0 && r.AreAllDudesWaiting(aliveDudes) {
				r.RemoveDudesFromWaiting()
				for _, d := range r.dudes {
					req.Add(RoomStartBossActivity{room: r, dude: d})
				}
				bossEnemy := r.kind.GetRoomEnemy(Huge, r.story.level)
				bossStack, err := render.NewStack("enemies/"+r.size.String(), strings.ToLower(bossEnemy.BossStack()), "")
				if err != nil {
					fmt.Println("Error creating boss stack for", bossEnemy.String(), err)
				}
				r.boss = NewEnemy(bossEnemy, r.story.level+1, bossStack)
			}
		}
	}

}

func (r *Room) Reset() {
	r.dudes = nil
	r.dudesInCenter = nil
	r.dudesInWaiting = nil
	r.boss = nil
	r.killedBoss = false
	r.combatTicks = 0
}

// Determins pan and vol of room track based on camera position
// TODO:
// - Move this out a bit so we can consolidate duplicate rooms and not set pan/vol twice for same track (take highest)
// - Determine by not only rotation but camera height, so scrolling up tower changes vol
func (r *Room) getPanVol(rads float64, multiplier float64) (float64, float64) {
	cR := rads
	rR := r.stacks[0].Rotation()

	// Determine pan and vol based on camera and room rotation
	pan := math.Cos(cR-rR) * 0.5
	vol := math.Sin(cR-rR) * 0.5

	vol *= multiplier

	// Return pan and vol
	return pan, vol
}

// Draw our room bits and bobs.
func (r *Room) Draw(o *render.Options) {
	if r.highlight {
		o.DrawImageOptions.ColorScale.Scale(2, 2, 2, 1)
	}
	r.stacks.Draw(o)
	r.walls.Draw(o)
	o.DrawImageOptions.ColorScale.Reset()

	if r.boss != nil {
		r.boss.Draw(*o)
	}
}

func (r *Room) AddDude(d *Dude) {
	r.dudes = append(r.dudes, d)
}

func (r *Room) RemoveDude(d *Dude) {
	for i, dude := range r.dudes {
		if dude == d {
			r.dudes = append(r.dudes[:i], r.dudes[i+1:]...)
			return
		}
	}
}

func (r *Room) AddDudeToWaiting(a *Dude) {
	r.dudesInWaiting = append(r.dudesInWaiting, a)
}
func (r *Room) IsDudeWaiting(a *Dude) bool {
	for _, dude := range r.dudesInWaiting {
		if dude == a {
			return true
		}
	}
	return false
}
func (r *Room) AreAllDudesWaiting(n int) bool {
	return len(r.dudesInWaiting) == n
}
func (r *Room) RemoveDudesFromWaiting() {
	r.dudesInWaiting = nil
}

func (r *Room) AddDudeToCenter(a *Dude) {
	r.dudesInCenter = append(r.dudesInCenter, a)
}
func (r *Room) RemoveDudeFromCenter(a *Dude) {
	for i, dude := range r.dudesInCenter {
		if dude == a {
			r.dudesInCenter = append(r.dudesInCenter[:i], r.dudesInCenter[i+1:]...)
			return
		}
	}
}
func (r *Room) IsDudeInCenter(a *Dude) bool {
	for _, dude := range r.dudesInCenter {
		if dude == a {
			return true
		}
	}
	return false
}

// Roll for new equipment from list
// Modifies the equipment by luck stat and room height
// Luck and room level determines chance of finding equipment, harder to find at higher levels
// Luck determines the quality
func (r *Room) RollLoot(luck int) *Equipment {
	if len(r.kind.Equipment()) == 0 {
		return nil
	}

	// Determine if we get equipment at all
	if rand.Intn(100) > ((luck+1)*10 - r.story.level) {
		return nil
	}

	// Determine the initial quality of the equipment based on luck
	fromLuck := float64(luck) / 5.0
	fromRoomLevel := float64(r.story.level) / 2.0
	initialQuality := EquipmentQuality((math.Floor(fromLuck + fromRoomLevel)))
	if initialQuality > EquipmentQualityLegendary {
		initialQuality = EquipmentQualityLegendary
	}

	// Determine if perk exists based on luck
	// Determine perk quality based on luck and room level
	hasPerk := rand.Intn(100) < luck
	var perk IPerk = nil
	if hasPerk {
		fromLuck = float64(luck) / 5.0
		fromRoomLevel = float64(r.story.level) / 2.0
		perkQuality := PerkQuality((math.Floor(fromLuck + fromRoomLevel)))
		if perkQuality > PerkQualityGodly {
			perkQuality = PerkQualityGodly
		}

		perk = GetRandomPerk(perkQuality)
	}

	// Create equipment
	list := r.kind.Equipment()
	equipmentName := list[rand.Intn(len(list))]
	equipment := NewEquipment(*equipmentName, 1, initialQuality, perk)
	if equipment == nil {
		return nil
	}

	// Level up the equipment based on floor level
	for i := 0; i < r.story.level; i++ {
		equipment.LevelUp(EquipmentQualityLegendary)
	}

	return equipment
}

func (r *Room) GetRoomEffect(e Event) Activity {
	if r == nil {
		return nil
	}
	switch e := e.(type) {
	case EventCombatRoom:
		switch r.kind {
		case Trap:
			// Damage dude based on stats
			e.dude.TrapDamage(r.story.level + 1)
			if e.dude.IsDead() {
				return DudeDeadActivity{dude: e.dude}
			}
		}
	case EventEnterRoom:
		switch r.kind {
		case Combat:
			// Add enemy based on room size
			enemyName := r.kind.GetRoomEnemy(r.size, r.story.level)
			enemyStack, err := render.NewStack("enemies/"+r.size.String(), "", "")
			// Randomize which enemy flavor it is
			enemyStack.SetStack(enemyStack.Stacks()[rand.Intn(len(enemyStack.Stacks()))])
			if err != nil {
				fmt.Println("Error creating enemy stack", err)
			} else {
				enemy := NewEnemy(enemyName, r.story.level, enemyStack)
				e.dude.enemy = enemy
			}
		}
	case EventLeaveRoom:
		// If enemy is attached to dude, remove it
		if e.dude.enemy != nil {
			e.dude.enemy = nil
		}
		// Add XP
		e.dude.AddXP(5 * (r.story.level + 1))

		// Roll for any loot if the room has any
		if r.kind.Equipment() != nil {
			// Roll for loot on exit
			if eq := r.RollLoot(e.dude.stats.luck); eq != nil {
				//fmt.Println(e.dude.name, "found", eq.Name())

				// Add to inventory and equip if slot is empty
				e.dude.AddToInventory(eq)
			}
		}
		// Add other leave events here
	case EventCenterRoom:
		// Add center room events here
		switch r.kind {
		case Armory:
			// Level up equipment
			maxQuality := EquipmentQuality(r.story.level/2 + 1)
			if maxQuality > EquipmentQualityLegendary {
				maxQuality = EquipmentQualityLegendary
			}
			if r.size == Large {
				e.dude.LevelUpEquipment(5, maxQuality)
			} else {
				e.dude.LevelUpEquipment(1, maxQuality)
			}
		case HealingShrine:
			// Heal
			stats := e.dude.GetCalculatedStats()
			switch r.size {
			case Small:
				// Heal 25% of max hp
				e.dude.Heal(stats.totalHp / 4)
			case Medium:
				// Heal 75% of max hp
				e.dude.Heal(stats.totalHp * 3 / 4)
			case Large:
				// Heal 100% of max hp
				e.dude.Heal(stats.totalHp)
			}
		case Curse:
			// Curse
			e.dude.Cursify(r.story.level + 1)
		case Well:
			// Restore all equipment uses
			e.dude.RestoreUses()
		case Combat:
			// He be in combat on entering room
		case Treasure:
			// Add gold
			goldAmount := (r.story.level + 1) * rand.Intn(10*int(r.size))
			e.dude.Trigger(EventGoldGain{dude: e.dude, amount: float64(goldAmount)})
		case Library:
			// Level up a random equipment perk or add one
			maxQuality := PerkQuality(r.story.level/2 + 1)
			if maxQuality > PerkQualityGodly {
				maxQuality = PerkQualityGodly
			}
			e.dude.Perkify(maxQuality)
		}
	}
	return nil
}

// For populating the required rooms to place
// Number of bad rooms based on requested size count.
// Stories 3, 6, 9, and 12 (?) are boss rooms
// TODO: make sure you can actually fit all required rooms
func GetRequiredRooms(storyLevel int, sizeCount int) []*RoomDef {
	level := storyLevel + 1

	if level%3 == 0 {
		roomDef := GetRoomDef(Boss, Huge)
		return []*RoomDef{roomDef}
	}

	potentialRooms := []RoomTemplate{}
	if level < 3 {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: Combat, size: Small},
			RoomTemplate{kind: Trap, size: Small},
		)
	} else if level < 6 {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: Combat, size: Small},
			RoomTemplate{kind: Trap, size: Small},
			RoomTemplate{kind: Combat, size: Medium},
			RoomTemplate{kind: Trap, size: Medium},
		)
	} else if level < 9 {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: Curse, size: Medium},
			RoomTemplate{kind: Combat, size: Small},
			RoomTemplate{kind: Trap, size: Small},
			RoomTemplate{kind: Combat, size: Medium},
			RoomTemplate{kind: Trap, size: Medium},
			RoomTemplate{kind: Combat, size: Large},
			RoomTemplate{kind: Trap, size: Large},
		)
	} else {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: Curse, size: Medium},
			RoomTemplate{kind: Combat, size: Small},
			RoomTemplate{kind: Trap, size: Small},
			RoomTemplate{kind: Combat, size: Medium},
			RoomTemplate{kind: Trap, size: Medium},
			RoomTemplate{kind: Combat, size: Large},
			RoomTemplate{kind: Trap, size: Large},
			RoomTemplate{kind: Combat, size: Huge},
			RoomTemplate{kind: Trap, size: Huge},
		)
	}
	rooms := make([]*RoomDef, 0)
	for i := 0; i < sizeCount; {
		room := potentialRooms[rand.Intn(len(potentialRooms))]
		if i+int(room.size) > sizeCount {
			continue
		}
		roomDef := GetRoomDef(room.kind, room.size)
		rooms = append(rooms, roomDef)
		i += int(room.size)
	}
	return rooms
}

// Returns an amount of rooms until the given size/space count is reached.
func GetOptionalRooms(storyLevel int, sizeCount int) []*RoomDef {
	level := storyLevel + 1

	// if we are at boss level, no optional rooms
	if level%3 == 0 {
		return nil
	}
	potentialRooms := []RoomTemplate{}
	potentialRooms = append(
		potentialRooms,
		RoomTemplate{kind: Armory, size: Medium},
		RoomTemplate{kind: HealingShrine, size: Small},
		RoomTemplate{kind: Well, size: Small},
		RoomTemplate{kind: Treasure, size: Small},
		RoomTemplate{kind: Library, size: Medium},
	)
	if level > 3 {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: HealingShrine, size: Medium},
			RoomTemplate{kind: Treasure, size: Medium},
		)
	}
	if level > 6 {
		potentialRooms = append(
			potentialRooms,
			RoomTemplate{kind: HealingShrine, size: Large},
			RoomTemplate{kind: Treasure, size: Large},
		)
	}

	rooms := make([]*RoomDef, 0)
	for i := 0; i < sizeCount; {
		room := potentialRooms[rand.Intn(len(potentialRooms))]
		if i+int(room.size) > sizeCount {
			continue
		}
		//
		roomDef := GetRoomDef(room.kind, room.size)
		rooms = append(rooms, roomDef)
		i += int(room.size)
	}
	return rooms
}
