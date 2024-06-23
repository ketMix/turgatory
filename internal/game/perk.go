package game

import (
	"fmt"
	"math/rand"
)

// Stubbing to maybe batch call perks based on trigger
//
// type PerkTrigger int
// const (
// 	PerkTriggerOnAttack PerkTrigger = iota
// )

// PerkQuality is the quality of a given perk.
type PerkQuality int

const (
	PerkQualityTrash PerkQuality = iota
	PerkQualityPathetic
	PerkQualityLesser
	PerkQualityCommon
	PerkQualityGreater
	PerkQualityGodly
)

func (pq *PerkQuality) String() string {
	switch *pq {
	case PerkQualityTrash:
		return "Trash"
	case PerkQualityPathetic:
		return "Pathetic"
	case PerkQualityLesser:
		return "Lesser"
	case PerkQualityCommon:
		return "Common"
	case PerkQualityGreater:
		return "Greater"
	case PerkQualityGodly:
		return "Godly"
	default:
		return "Unknown"
	}
}

func (pq PerkQuality) Color() string {
	switch pq {
	case PerkQualityTrash:
		return "gray"
	case PerkQualityPathetic:
		return "white"
	case PerkQualityLesser:
		return "green"
	case PerkQualityCommon:
		return "blue"
	case PerkQualityGreater:
		return "purple"
	case PerkQualityGodly:
		return "orange"
	default:
		return "white"
	}
}

// etc.
func constructName(ability string, quality PerkQuality, modifier *string) string {

	// Construct the name
	// Lesser Stat Boost of Strength On Attack

	name := quality.String()
	name += " "
	name += ability
	name += " "
	if modifier != nil {
		name += "of " + *modifier
	}

	// Probably don't need this if we're assigning events directly to perks
	// triggerText := ""

	// // Join them with or
	// triggersText := ""
	// for i, trigger := range triggers {
	// 	if i == 0 {
	// 		triggersText += " when " + trigger
	// 	} else {
	// 		triggersText += " or " + trigger
	// 	}
	// }
	// name += triggerText

	return name
}

// IPerk is our interface for perks that can be applied to equipment, dudes, etc.
type IPerk interface {
	Check(Event) bool
	Name() string   // Name of the perk
	String() string // Full name of the perk
	Description() string
	Quality() PerkQuality
	LevelUp(PerkQuality)
	LevelDown()
}

type Perk struct {
	IPerk
	quality PerkQuality
}

func (p Perk) String() string {
	return "What is zis, zis is nothing!"
}
func (p Perk) Description() string {
	return ""
}
func (p Perk) Quality() PerkQuality {
	return p.quality
}

func (p Perk) Check(e Event) bool {
	return false
}
func (p Perk) LevelUp(maxQuality PerkQuality) {
	if p.quality >= maxQuality {
		return
	}
	p.quality++
}

func (p Perk) LevelDown() {
	if p.quality <= PerkQualityTrash {
		return
	}
	p.quality--
}

// PerkNone represents an empty perk. Not sure if this will be used.
type PerkNone struct {
	Perk
}

func (p PerkNone) Quality() PerkQuality {
	return PerkQualityTrash
}

// PerkFindGold finds gold based upon the quality of the perk.
// +0.25 per quality level per room.
type PerkFindGold struct {
	Perk
}

func (p PerkFindGold) chance() float64 {
	return 0.25 * float64(p.quality)
}

func (p PerkFindGold) Name() string {
	return "Find Gold"
}

func (p PerkFindGold) String() string {
	return constructName(p.Name(), p.quality, nil)
}

func (p PerkFindGold) Description() string {
	amount := float32(p.quality) * 0.25
	return fmt.Sprintf("Has a Chance to find finds %f gold", amount)
}

func (p PerkFindGold) Check(e Event) bool {
	switch e := e.(type) {
	case EventEnterRoom:
		amount := float32(p.quality) * 0.25
		e.dude.UpdateGold(amount)
		return true
	}
	return false
}

// PerkStatBoost is a perk that modifies a stat based on teh quality of the perk.
// +1 target stat per quality level.
type PerkStatBoost struct {
	Perk
	stat Stat
}

func (p PerkStatBoost) Name() string {
	return "Stat Boost"
}

func (p PerkStatBoost) String() string {
	statStr := string(p.stat)
	return constructName(p.Name(), p.quality, &statStr)
}

func (p PerkStatBoost) Description() string {
	if p.stat == "" {
		return "No bonus! How sad."
	}
	return fmt.Sprintf("Boosts %s stat by %d", p.stat, p.quality)
}

func (p PerkStatBoost) Check(e Event) bool {
	switch e := e.(type) {
	case EventEquip:
		e.dude.Stats().ModifyStat(p.stat, int(p.quality))
		return true
	case EventUnequip:
		e.dude.Stats().ModifyStat(p.stat, -int(p.quality))
		return true
	}
	return false
}

// PerkHeal heals dude based on wisdom when entering a room.
type PerkHealOnRoomEnter struct {
	Perk
}

func (p PerkHealOnRoomEnter) amount(wisdom int) int {
	return int(p.quality) * wisdom
}

func (p PerkHealOnRoomEnter) Name() string {
	return "Heal On Room Enter"
}

func (p PerkHealOnRoomEnter) String() string {
	return constructName(p.Name(), p.quality, nil)
}

func (p PerkHealOnRoomEnter) Description() string {
	return fmt.Sprintf("Heals %d * wisdom on room enter", int(p.quality))
}

func (p PerkHealOnRoomEnter) Check(e Event) bool {
	switch e := e.(type) {
	case EventEnterRoom:
		e.dude.Heal(int(p.amount(e.dude.stats.wisdom)))
	}
	return false
}

// PerkHeal heals all dudes based on quality when equip is sold
type PerkHealOnSell struct {
	Perk
}

func (p PerkHealOnSell) amount() int {
	return int(p.quality) * 10
}
func (p PerkHealOnSell) Name() string {
	return "Heal On Sell"
}

func (p PerkHealOnSell) String() string {
	return constructName(p.Name(), p.quality, nil)
}

func (p PerkHealOnSell) Description() string {
	return fmt.Sprintf("Heals %d when sold", p.amount())
}

func (p PerkHealOnSell) Check(e Event) bool {
	switch e := e.(type) {
	case EventSell:
		// Heal all dudes
		fmt.Println("Unimplemented", e.String())
		return true
	}
	return false
}

func GetRandomPerk(quality PerkQuality) IPerk {
	// Randomly select a perk
	perk := Perk{
		quality: quality,
	}

	// Set of all perks
	perkList := []IPerk{
		PerkFindGold{perk},
		PerkStatBoost{perk, StatStrength},
		PerkStatBoost{perk, StatWisdom},
		PerkStatBoost{perk, StatDefense},
		PerkStatBoost{perk, StatAgility},
		PerkStatBoost{perk, StatCowardice},
		PerkStatBoost{perk, StatLuck},
		PerkStatBoost{perk, StatHP},
		PerkHealOnRoomEnter{perk},
		PerkHealOnSell{perk},
	}

	// Randomly select a perk
	index := rand.Intn(len(perkList))
	return perkList[index]
}
