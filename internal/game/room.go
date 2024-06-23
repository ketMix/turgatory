package game

import (
	"fmt"
	"math"
	"math/rand"

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

// RoomStairsEntrance is the distance from the center that a room's stairs is expected to be at.
const RoomStairsEntrance = 12
const RoomPath = 53
const TowerStairs = 60
const TowerEntrance = 80

// RoomKind is an enumeration of the different kinds of rooms in za toweru.
type RoomKind int

func (r *RoomKind) String() string {
	switch *r {
	case Armory:
		return "armory"
	case HealingShrine:
		return "healing"
	case Treasure:
		return "treasure"
	case Combat:
		//return "combat"
		return "template"
	case Well:
		return "well"
	case Library:
		return "library"
	default:
		return "Unknown"
	}
}

// Equipment you can find in room
// hmm, kinda defeats the current structure of equipment in yaml if we have to add new equipment here
func (r *RoomKind) Equipment() []string {
	switch *r {
	case Armory:
		return []string{"Sword", "Shield", "Bow", "Book", "Robe", "Plate", "Leather"}
	case HealingShrine:
		return []string{"Ring", "Necklace"} // temporary
	}
	return []string{}
}

const (
	Empty RoomKind = iota
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
)

// Room is a room within a story of za toweru.
type Room struct {
	story *Story
	index int // Reference to the index within the story.
	kind  RoomKind
	size  RoomSize
	power int // ???

	stacks         render.Stacks
	walls          render.Stacks
	actorsInCenter []Actor
}

func NewRoom(size RoomSize, kind RoomKind) *Room {
	r := &Room{
		size: size,
		kind: kind,
	}

	stack, err := render.NewStack(fmt.Sprintf("rooms/%s", size.String()), kind.String(), "")
	if err != nil {
		panic(err)
	}
	r.stacks.Add(stack)

	// Add our walls.
	for j := 0; j < 3; j++ {
		for i := 0; i < 8; i++ {
			stack, err := render.NewStack(fmt.Sprintf("walls/%s", size.String()), "template", "base")
			if err != nil {
				continue
			}
			stack.VgroupOffset = j * StoryWallHeight
			r.walls.Add(stack)
		}
	}

	return r
}

// Update updates the stuff in the room.
func (r *Room) Update() {
	r.stacks.Update()
}

// Draw our room bits and bobs.
func (r *Room) Draw(o *render.Options) {
	r.stacks.Draw(o)
	r.walls.Draw(o)
}

func (r *Room) AddActorToCenter(a Actor) {
	r.actorsInCenter = append(r.actorsInCenter, a)
}

func (r *Room) RemoveActorFromCenter(a Actor) {
	for i, actor := range r.actorsInCenter {
		if actor == a {
			r.actorsInCenter = append(r.actorsInCenter[:i], r.actorsInCenter[i+1:]...)
			return
		}
	}
}

func (r *Room) IsActorInCenter(a Actor) bool {
	for _, actor := range r.actorsInCenter {
		if actor == a {
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
	equipment := NewEquipment(equipmentName, 1, initialQuality, perk)
	if equipment == nil {
		return nil
	}

	// Level up the equipment based on floor level
	for i := 0; i < r.story.level; i++ {
		equipment.LevelUp()
	}

	return equipment
}

func (r *Room) GetRoomEffect(e Event) {
	switch e := e.(type) {
	case EventLeaveRoom:
		// Add XP
		e.dude.AddXP(1 * r.story.level)

		// Roll for any loot if the room has any
		if r.kind.Equipment() != nil {
			// Roll for loot on exit
			if eq := r.RollLoot(e.dude.stats.luck); eq != nil {
				fmt.Println(e.dude.name, "found", eq.Name())

				// Add to inventory and equip if slot is empty
				e.dude.AddToInventory(eq)
			}
		}
		// Add other leave events here
	case EventCenterRoom:
		fmt.Println("Center room", r.kind.String())
		// Add center room events here
		switch r.kind {
		case Armory:
			// Level up equipment
			if r.size == Large {
				e.dude.LevelUpEquipment(5)
			} else {
				e.dude.LevelUpEquipment(1)
			}
		case HealingShrine:
			// Heal
			e.dude.Heal((r.story.level + 1) * 5)
		case Well:
			// Restore equipment uses
			e.dude.RestoreUses(r.story.level + 1)
		case Combat:
			// Damage the dude
			e.dude.Damage(10 * (r.story.level + 1))
		case Treasure:
			// Add gold
			goldAmount := (r.story.level + 1) * rand.Intn(10*int(r.size))
			e.dude.UpdateGold(float32(goldAmount))
		case Library:
			// Level up a random equipment perk or add one
			maxQuality := PerkQuality(r.story.level + 1)
			if maxQuality > PerkQualityGodly {
				maxQuality = PerkQualityGodly
			}
			e.dude.Perkify(maxQuality)
		}
	}
}
