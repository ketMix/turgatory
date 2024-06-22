package game

import "fmt"

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

// Perk is our interface for perks that can be applied to equipment, dudes, etc.
type Perk interface {
	Check(Event) bool
	String() string
	Description() string
}

// PerkNone represents an empty perk. Not sure if this will be used.
type PerkNone struct{}

func (p PerkNone) String() string {
	return "What is zis, zis is nothing!"
}

func (p PerkNone) Description() string {
	return ""
}

func (p PerkNone) Check(e Event) bool {
	return false
}

// PerkFindGold finds gold based upon the quality of the perk.
// +0.25 per quality level per room.
type PerkFindGold struct {
	quality  PerkQuality
	modifier *string // What be this?
}

func (p PerkFindGold) String() string {
	return "Find Gold"
}

func (p PerkFindGold) Description() string {
	amount := float32(p.quality) * 0.25
	return fmt.Sprintf("Finds %f gold", amount)
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
	quality PerkQuality
	stat    Stat
}

func (p PerkStatBoost) String() string {
	return "Stat Boost"
}

func (p PerkStatBoost) Name() string {
	triggers := []string{EventEquip{}.String(), EventUnequip{}.String()}
	return constructName(p.String(), p.quality, triggers, string(p.stat))
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

// etc.
func constructName(ability string, quality PerkQuality, triggers []string, modifier string) string {
	triggerText := ""

	// Join them with and
	triggersText := ""
	for i, trigger := range triggers {
		if i == 0 {
			triggersText += " when " + trigger
		} else {
			triggersText += " or " + trigger
		}
	}

	// Construct the name
	// Lesser Stat Boost of Strength On Attack

	name := string(quality)
	name += " "
	name += ability
	name += " "
	if modifier != "" {
		name += "of " + modifier
	}

	name += triggerText

	return name
}
