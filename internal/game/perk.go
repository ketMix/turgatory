package game

import "fmt"

// Stubbing to maybe batch call perks based on trigger
//
// type PerkTrigger int
// const (
// 	PerkTriggerOnAttack PerkTrigger = iota
// )

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

type PerkAbility int

const (
	AbilityStatBoost PerkAbility = iota
	AbilityFindGold
)

func (pa *PerkAbility) String() string {
	switch *pa {
	case AbilityStatBoost:
		return "Stat Boost"
	case AbilityFindGold:
		return "Find Gold"
	default:
		return "Unknown"
	}
}

type ActivateFunc func(*Event, *Dude, *Equipment) bool

type Perk struct {
	Activate    ActivateFunc
	Name        func() string
	Description func() string
}

var NO_EFFECT_METHOD = func(e *Event, d *Dude, eq *Equipment) bool {
	return false
}

// Returns perk object with the given quality, ability, triggers, and modifier
func NewPerk(quality PerkQuality, ability PerkAbility, triggers []EventName, modifier *string) *Perk {
	constructedName := constructName(quality, ability, triggers, modifier)
	effect, description := getAbility(ability, quality, modifier)

	// Set the activate function
	// Checks the event against the triggers
	activate := func(e *Event, d *Dude, eq *Equipment) bool {
		// If the event is in the triggers, then apply the effect and return success
		for _, trigger := range triggers {
			if trigger == e.Name() {
				return effect(e, d, eq)
			}
		}
		return false
	}

	return &Perk{
		Activate: activate,
		Description: func() string {
			return description
		},
		Name: func() string {
			return constructedName
		},
	}
}

/**
 * Get Ability
 * Returns the ability function and the description of the ability
 *
 * Just a fun switch statement through abilitiy names
 */
func getAbility(ability PerkAbility, quality PerkQuality, modifier *string) (ActivateFunc, string) {
	switch ability {
	case AbilityStatBoost:
		return statBoost(quality, modifier)
	default:
		return NO_EFFECT_METHOD, "What is zis, zis is nothing!"
	}
}

func constructName(quality PerkQuality, ability PerkAbility, triggers []EventName, modifier *string) string {
	triggerText := ""

	// Join them with and
	triggersText := ""
	for i, trigger := range triggers {
		if i == 0 {
			triggersText += " when " + trigger.String()
		} else {
			triggersText += " or " + trigger.String()
		}
	}

	// Construct the name
	// Lesser Stat Boost of Strength On Attack

	name := quality.String()
	name += " "
	name += ability.String()
	name += " "
	if modifier != nil {
		name += "of " + *modifier
	}

	name += triggerText

	return name
}

/**
 * Stat Boost
 * Stat boost is a perk that modifies a stat based on the quality of the perk.
 * +1 per quality level
 */
func statBoost(quality PerkQuality, modifier *string) (ActivateFunc, string) {
	if modifier == nil {
		return NO_EFFECT_METHOD, "No bonus! How sad."
	}

	description := fmt.Sprintf("Boosts %s stat by %d", Stat(*modifier), int(quality))

	return func(e *Event, d *Dude, eq *Equipment) bool {
		d.Stats().ModifyStat(Stat(*modifier), int(quality))
		return true
	}, description
}

/**
 * Find Gold
 * The amount of gold is based on the quality of the perk.
 * +0.25 per quality level
 */
func findGold(quality PerkQuality, modifier *string) (ActivateFunc, string) {
	amount := float32(quality) * 0.25
	description := fmt.Sprintf("Finds %f gold", amount)

	return func(e *Event, d *Dude, eq *Equipment) bool {
		d.UpdateGold(amount)
		return true
	}, description
}
