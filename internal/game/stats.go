package game

import (
	"fmt"
	"math"
	"math/rand"
)

// A dude's inherent capability to deal with their environment
// (how they survive dem rooms)
type Stat string

const (
	StatStrength   Stat = "Str"
	StatWisdom     Stat = "Wis"
	StatDefense    Stat = "Def"
	StatAgility    Stat = "Agi"
	StatConfidence Stat = "Con"
	StatLuck       Stat = "Luc"
	StatMaxHP      Stat = "MaxHP"
	StatCurrentHP  Stat = "CurrentHP"
)

type Stats struct {
	level      int // how much they've grown
	currentHp  int // how dead are they
	totalHp    int // how dead could they not be
	strength   int // how hard they hit in combat
	wisdom     int // how well they heal, influences stat variation
	defense    int // reduces enemy attack
	agility    int // how fast they zip
	confidence int // combat priority (who gets hit first)
	luck       int // how lucky they are

	levelUpChange *Stats
}

func StatsSort(a, b Stats) bool {
	return a.level < b.level
}

func (s *Stats) Print() {
	fmt.Printf("Level: %d\n", s.level)
	fmt.Printf("HP: %d/%d\n", s.currentHp, s.totalHp)
	fmt.Printf("Strength: %d\n", s.strength)
	fmt.Printf("Wisdom: %d\n", s.wisdom)
	fmt.Printf("Defense: %d\n", s.defense)
	fmt.Printf("Agility: %d\n", s.agility)
	fmt.Printf("Confidence: %d\n", s.confidence)
	fmt.Printf("Luck: %d\n", s.luck)
}

// ApplyLevelUp applies the level up changes to the stats
// with some variance depend on wisdom
func (s *Stats) LevelUp(isEnemy bool) {
	// what did you think was going to happen
	s.level += 1

	// variance is a random number between 1 and (wisdom/5) + 1
	// lowest variance is 0.75, highest is 1.25
	variance := func() float64 {
		return rand.Float64()*0.5 + 0.75
	}

	// If is enemy, dont' apply variance
	if isEnemy {
		variance = func() float64 {
			return 1
		}
	}

	getValue := func(base int) int {
		return int(math.Floor(float64(base) * variance()))
	}

	// apply the variance to the stats
	s.ModifyStat(StatMaxHP, getValue(s.levelUpChange.totalHp))
	s.ModifyStat(StatStrength, getValue(s.levelUpChange.strength))
	s.ModifyStat(StatWisdom, getValue(s.levelUpChange.wisdom))
	s.ModifyStat(StatDefense, getValue(s.levelUpChange.defense))
	s.ModifyStat(StatAgility, getValue(s.levelUpChange.agility))
	s.ModifyStat(StatConfidence, getValue(s.levelUpChange.confidence))
	s.ModifyStat(StatLuck, getValue(s.levelUpChange.luck))

	// a blesing from jesus himself
	s.currentHp = s.totalHp
}

// ApplyLevelDown applies the level down changes to the stats
func (s *Stats) LevelDown() {
	// the devil refuses to heal you
	s.currentHp += 0

	s.ModifyStat(StatMaxHP, -s.levelUpChange.totalHp)
	s.ModifyStat(StatStrength, -s.levelUpChange.strength)
	s.ModifyStat(StatWisdom, -s.levelUpChange.wisdom)
	s.ModifyStat(StatDefense, -s.levelUpChange.defense)
	s.ModifyStat(StatAgility, -s.levelUpChange.agility)
	s.ModifyStat(StatConfidence, -s.levelUpChange.confidence)
	s.ModifyStat(StatLuck, -s.levelUpChange.luck)

	// cant forget this part
	s.level -= 1
}

func (s *Stats) ModifyStat(stat Stat, amount int) {
	switch stat {
	case StatStrength:
		s.strength += amount
		if s.strength < 0 {
			s.strength = 0
		}
	case StatWisdom:
		s.wisdom += amount
		if s.wisdom < 0 {
			s.wisdom = 0
		}
	case StatDefense:
		s.defense += amount
		if s.defense < 0 {
			s.defense = 0
		}
	case StatAgility:
		s.agility += amount
		if s.agility < 0 {
			s.agility = 0
		}
	case StatConfidence:
		s.confidence += amount
		if s.confidence < 0 {
			s.confidence = 0
		}
	case StatLuck:
		s.luck += amount
		if s.luck < 0 {
			s.luck = 0
		}
	// HP is a special case
	case StatMaxHP:
		prevHp := s.totalHp
		s.totalHp += amount
		if s.totalHp < 1 {
			s.totalHp = 1
		}

		// If the total HP is increased, increase the current HP by the same amount
		if s.totalHp > prevHp {
			s.currentHp += amount
		}
	default:
		fmt.Printf("Unknown stat %s\n", stat)
	}
}

func (s *Stats) ApplyDefense(damage int) int {
	// Apply defense stat using a logarithmic function
	// for diminishing returns
	defenseReduction := 1 - 1/(1+math.Log1p(float64(s.defense)/100))

	// Apply defense reduction
	reducedDamage := float64(damage) * (1 - defenseReduction)

	// Round the result and convert back to int
	return max(1, int(math.Round(reducedDamage)))
}

func NewStats(levelUpChange *Stats, isEnemy bool) *Stats {
	// start the stats at a negative level
	// then level up a few times in order to set the starting stats
	startingLevels := 3

	// Don't start enemy with extra levels
	if isEnemy {
		startingLevels = 0
	}

	// Added this because GetCalculatedStats passes nil to this func... --kts
	level := -startingLevels
	if levelUpChange != nil {
		level = levelUpChange.level - startingLevels
	}
	stats := &Stats{
		level:         level,
		levelUpChange: levelUpChange,
		totalHp:       0,
		strength:      0,
		wisdom:        0,
		defense:       0,
		agility:       0,
		confidence:    0,
		luck:          0,
	}
	if levelUpChange == nil {
		return stats
	}

	// Level up the stats a few times
	for i := 0; i < startingLevels; i++ {
		stats.LevelUp(isEnemy)
	}
	return stats
}

func (s *Stats) Add(a *Stats) *Stats {
	stats := &Stats{
		currentHp:  s.currentHp,
		totalHp:    s.totalHp,
		strength:   s.strength,
		wisdom:     s.wisdom,
		defense:    s.defense,
		agility:    s.agility,
		confidence: s.confidence,
		luck:       s.luck,
	}

	if a == nil {
		return stats
	}

	// Modify the stats using wrapper
	stats.ModifyStat(StatMaxHP, a.totalHp)
	stats.ModifyStat(StatStrength, a.strength)
	stats.ModifyStat(StatWisdom, a.wisdom)
	stats.ModifyStat(StatDefense, a.defense)
	stats.ModifyStat(StatAgility, a.agility)
	stats.ModifyStat(StatConfidence, a.confidence)
	stats.ModifyStat(StatLuck, a.luck)

	return stats
}
