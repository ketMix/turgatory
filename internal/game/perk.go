package game

import (
	"fmt"
	"image/color"
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

func (pq PerkQuality) TextColor() color.Color {
	switch pq {
	case PerkQualityTrash: // gray
		return color.RGBA{150, 150, 150, 255}
	case PerkQualityPathetic: // white
		return color.RGBA{200, 200, 200, 255}
	case PerkQualityLesser: // green
		return color.RGBA{75, 250, 75, 255}
	case PerkQualityCommon: // blue
		return color.RGBA{75, 75, 250, 255}
	case PerkQualityGreater: // purple
		return color.RGBA{250, 75, 250, 255}
	case PerkQualityGodly: // orange
		return color.RGBA{250, 75, 75, 255}
	default:
		return color.RGBA{200, 200, 200, 255}
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

func (p *Perk) String() string {
	return "What is zis, zis is nothing!"
}

func (p *Perk) Description() string {
	return ""
}

func (p *Perk) Quality() PerkQuality {
	return p.quality
}

func (p *Perk) Check(e Event) bool {
	return false
}
func (p *Perk) LevelUp(maxQuality PerkQuality) {
	if p.quality >= maxQuality {
		return
	}
	p.quality++
}

func (p *Perk) LevelDown() {
	if p.quality <= PerkQualityTrash {
		return
	}
	p.quality--
}

func (p *Perk) Name() string {
	return constructName(p.String(), p.quality, nil)
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
	*Perk
}

func (p PerkFindGold) chance() float64 {
	return 0.25 * float64(p.quality)
}

func (p PerkFindGold) amount() int {
	return int(p.quality) * 5
}

func (p PerkFindGold) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkFindGold) String() string {
	return "Find Gold"
}

func (p PerkFindGold) Description() string {
	return fmt.Sprintf("Has a %.2f Chance to find finds %d gold", p.chance()*100, p.amount())
}

func (p PerkFindGold) Check(e Event) bool {
	switch e := e.(type) {
	case EventEnterRoom:
		e.dude.Trigger(EventGoldGain{dude: e.dude, amount: p.amount()})
		return true
	}
	return false
}

// PerkStatBoost is a perk that modifies a stat based on teh quality of the perk.
// +1 target stat per quality level.
type PerkStatBoost struct {
	*Perk
	stat Stat
}

func (p PerkStatBoost) amount() int {
	if p.stat == StatMaxHP {
		return int(p.quality) * 10
	}
	return (int(p.quality) + 1) * 3
}

func (p PerkStatBoost) Name() string {
	statStr := string(p.stat)
	return constructName(p.String(), p.quality, &statStr)
}

func (p PerkStatBoost) String() string {
	return "Stat Boost"
}

func (p PerkStatBoost) Description() string {
	if p.stat == "" {
		return "No bonus! How sad."
	}
	return fmt.Sprintf("Boosts %s stat by %d", p.stat, p.amount())
}

func (p PerkStatBoost) Check(e Event) bool {
	switch e := e.(type) {
	case EventEquip:
		e.dude.Stats().ModifyStat(p.stat, p.amount())
		return false
	case EventUnequip:
		e.dude.Stats().ModifyStat(p.stat, p.amount())
		return false
	}
	return false
}

// PerkHeal heals dude based on wisdom when entering a room.
type PerkHealOnRoomEnter struct {
	*Perk
}

func (p PerkHealOnRoomEnter) amount(wisdom int) int {
	return (wisdom / 4) * int(p.quality+1)
}

func (p PerkHealOnRoomEnter) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnRoomEnter) String() string {
	return "Heal On Room Enter"
}

func (p PerkHealOnRoomEnter) Description() string {
	return fmt.Sprintf("Heals %d * wisdom/4 on room enter", int(p.quality))
}

func (p PerkHealOnRoomEnter) Check(e Event) bool {
	switch e := e.(type) {
	case EventEnterRoom:
		if e.dude == nil {
			return false
		}
		stats := e.dude.GetCalculatedStats()
		amount := e.dude.Heal(p.amount(stats.wisdom))
		if amount == 0 {
			return false
		}
		return true
	}
	return false
}

// PerkHealOnSell heals all dudes based on quality when equip is sold
type PerkHealOnSell struct {
	*Perk
}

func (p PerkHealOnSell) amount() int {
	return int(p.quality+1) * 10
}

func (p PerkHealOnSell) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnSell) String() string {
	return "Heal On Sell"
}

func (p PerkHealOnSell) Description() string {
	return fmt.Sprintf("Heals all dudes for %d when sold", p.amount())
}

func (p PerkHealOnSell) Check(e Event) bool {
	switch e := e.(type) {
	case EventSell:
		// Heal all dudes
		for _, dude := range e.dudes {
			dude.Heal(p.amount())
		}
		return true
	}
	return false
}

// PerkHealOnGoldGain heals dude when they get gold
type PerkHealOnGoldGain struct {
	*Perk
}

func (p PerkHealOnGoldGain) amount() int {
	return int(p.quality+1) * 1
}

func (p PerkHealOnGoldGain) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnGoldGain) String() string {
	return "Smell of Gold"
}

func (p PerkHealOnGoldGain) Description() string {
	return fmt.Sprintf("Heals dude for %d when they gain gold", p.amount())
}

func (p PerkHealOnGoldGain) Check(e Event) bool {
	switch e := e.(type) {
	case EventGoldGain:
		e.dude.Heal(p.amount())
		return true
	}
	return false
}

// PerkHealOnGoldGain heals dude when they get gold
type PerkHealOnGoldLoss struct {
	*Perk
}

func (p PerkHealOnGoldLoss) amount() int {
	return int(p.quality+1) * 1
}

func (p PerkHealOnGoldLoss) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnGoldLoss) String() string {
	return "Food Tax"
}

func (p PerkHealOnGoldLoss) Description() string {
	return fmt.Sprintf("Heals dude for %d when they lose gold", p.amount())
}

func (p PerkHealOnGoldLoss) Check(e Event) bool {
	switch e := e.(type) {
	case EventGoldGain:
		e.dude.Heal(p.amount())
		return true
	}
	return false
}

// PerkStickyFingers reduces gold loss
type PerkStickyFingers struct {
	*Perk
}

func (p PerkStickyFingers) amount() float64 {
	return float64(p.quality+1) * 0.1
}

func (p PerkStickyFingers) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkStickyFingers) String() string {
	return "Sticky Fingers"
}

func (p PerkStickyFingers) Description() string {
	return fmt.Sprintf("Reduces gold loss by %.2f percent", p.amount()*100)
}

func (p PerkStickyFingers) Check(e Event) bool {
	switch e := e.(type) {
	case EventGoldLoss:
		e.dude.UpdateGold(int(float64(-e.amount) * p.amount()))
		return true
	}
	return false
}

// PerkGoldBoost increases gold gain
type PerkGoldBoost struct {
	*Perk
}

func (p PerkGoldBoost) amount() float64 {
	return float64(p.quality+1) * 0.1
}

func (p PerkGoldBoost) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkGoldBoost) String() string {
	return "Miser's Touch"
}

func (p PerkGoldBoost) Description() string {
	return fmt.Sprintf("Increases gold gain by %.2f percent", p.amount()*100)
}

func (p PerkGoldBoost) Check(e Event) bool {
	switch e := e.(type) {
	case EventGoldGain:
		e.dude.UpdateGold(int(float64(e.amount) * p.amount()))
		return true
	}
	return false
}

type PerkHealOnDodge struct {
	*Perk
}

func (p PerkHealOnDodge) amount() int {
	return int(p.quality+1) * 2
}

func (p PerkHealOnDodge) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnDodge) String() string {
	return "Narrow Recovery"
}

func (p PerkHealOnDodge) Description() string {
	return fmt.Sprintf("Heals dude for %d when they dodge", p.amount())
}

func (p PerkHealOnDodge) Check(e Event) bool {
	switch e := e.(type) {
	case EventDudeDodge:
		e.dude.Heal(p.amount())
		return true
	}
	return false
}

type PerkHealOnCrit struct {
	*Perk
}

func (p PerkHealOnCrit) amount() int {
	return int(p.quality+1) * 4
}

func (p PerkHealOnCrit) Name() string {
	return constructName(p.String(), p.quality, nil)
}

func (p PerkHealOnCrit) String() string {
	return "Crit Heal"
}

func (p PerkHealOnCrit) Description() string {
	return fmt.Sprintf("Heals dude for %d when they crit", p.amount())
}

func (p PerkHealOnCrit) Check(e Event) bool {
	switch e := e.(type) {
	case EventDudeCrit:
		e.dude.Heal(p.amount())
		return true
	}
	return false
}

func GetRandomPerk(quality PerkQuality) IPerk {
	// Randomly select a perk
	perk := &Perk{
		quality: quality,
	}

	// Set of all perks
	perkList := []IPerk{
		PerkFindGold{perk},
		PerkStatBoost{perk, StatStrength},
		PerkStatBoost{perk, StatWisdom},
		PerkStatBoost{perk, StatDefense},
		PerkStatBoost{perk, StatAgility},
		PerkStatBoost{perk, StatLuck},
		PerkStatBoost{perk, StatMaxHP},
		PerkHealOnRoomEnter{perk},
		PerkHealOnSell{perk},
		PerkHealOnGoldGain{perk},
		PerkHealOnDodge{perk},
		PerkHealOnCrit{perk},
		PerkGoldBoost{perk},

		// Curse dependent perks
		// PerkStickyFingers{perk},
		// PerkHealOnGoldLoss{perk},
	}

	// Randomly select a perk
	index := rand.Intn(len(perkList))
	return perkList[index]
}
