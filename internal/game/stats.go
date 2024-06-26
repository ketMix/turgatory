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
	StatStrength  Stat = "Str"
	StatWisdom    Stat = "Wis"
	StatDefense   Stat = "Def"
	StatAgility   Stat = "Agi"
	StatCowardice Stat = "Cow"
	StatLuck      Stat = "Luc"
	StatMaxHP     Stat = "MaxHP"
	StatCurrentHP Stat = "CurrentHP"
)

type Stats struct {
	level     int // how much they've grown
	currentHp int // how dead are they
	totalHp   int // how dead could they not be
	strength  int // how hard they hit in combat
	wisdom    int // how well they heal, influences stat variation
	defense   int // reduces enemy attack
	agility   int // how fast they zip
	cowardice int // combat priority (who gets hit first)
	luck      int // how lucky they are

	levelUpChange *Stats
}

const WISDOM_PER_VARIANCE = 5

// ApplyLevelUp applies the level up changes to the stats
// with some variance depend on wisdom
func (s *Stats) LevelUp() {
	// what did you think was going to happen
	s.level += 1

	// variance is a random number between 1 and (wisdom/5) + 1
	// need 5 wisdom to get a variance of 2
	// lowest variance is 1
	// multiplier can get pretty high with his wisdom and level
	level := s.level
	if level <= 0 {
		level = 1
	}
	variance := func() float64 {
		return 1 + rand.Float64()*(float64(s.wisdom)/WISDOM_PER_VARIANCE)
	}
	getValue := func(base int) int {
		return int(math.Round(float64(base) * variance()))
	}

	// apply the variance to the stats
	s.ModifyStat(StatMaxHP, getValue(s.levelUpChange.totalHp))
	s.ModifyStat(StatStrength, getValue(s.levelUpChange.strength))
	s.ModifyStat(StatWisdom, getValue(s.levelUpChange.wisdom))
	s.ModifyStat(StatDefense, getValue(s.levelUpChange.defense))
	s.ModifyStat(StatAgility, getValue(s.levelUpChange.agility))
	s.ModifyStat(StatCowardice, getValue(s.levelUpChange.cowardice))
	s.ModifyStat(StatLuck, getValue(s.levelUpChange.luck))

	// a blesing from jesus himself
	s.currentHp = s.totalHp
}

// ApplyLevelDown applies the level down changes to the stats
// with some variance depend on wisdom
func (s *Stats) LevelDown() {
	// the devil refuses to heal you
	s.currentHp += 0

	// instead of using variance to determine the amount it increases,
	// we use it to determine if the stat is skipped from being reduced
	// variance is a number between 1 and (wisdom/5) + 1
	// lowest variance is 0.1
	level := s.level
	if level <= 0 {
		level = 1
	}
	threshold := 0.1 * (float64(s.wisdom)/WISDOM_PER_VARIANCE + float64(level))

	// conditionally apply the variance to the stats
	if rand.Float64() > threshold {
		s.ModifyStat(StatMaxHP, -s.levelUpChange.totalHp)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatStrength, -s.levelUpChange.strength)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatWisdom, -s.levelUpChange.wisdom)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatDefense, -s.levelUpChange.defense)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatAgility, -s.levelUpChange.agility)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatCowardice, -s.levelUpChange.cowardice)
	}
	if rand.Float64() > threshold {
		s.ModifyStat(StatLuck, -s.levelUpChange.luck)
	}

	// cant forget this part
	s.level -= 1
}

func (s *Stats) ModifyStat(stat Stat, amount int) {
	switch stat {
	case StatStrength:
		s.strength += amount
	case StatWisdom:
		s.wisdom += amount
	case StatDefense:
		s.defense += amount
	case StatAgility:
		s.agility += amount
	// Cowardice can't go below 0
	case StatCowardice:
		s.cowardice += amount
		if s.cowardice < 0 {
			s.cowardice = 0
		}
	case StatLuck:
		s.luck += amount
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

	case StatCurrentHP:
		s.currentHp += amount
		if s.currentHp < 0 {
			s.currentHp = 0
		}
		if s.currentHp > s.totalHp {
			s.currentHp = s.totalHp
		}
	default:
		fmt.Printf("Unknown stat %s\n", stat)
	}
}

func NewStats(levelUpChange *Stats) *Stats {
	// start the stats at a negative level
	// then level up a few times in order to set the starting stats
	startingLevels := 3
	// Added this because GetCalculatedStats passes nil to this func... --kts
	level := -3
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
		cowardice:     0,
		luck:          0,
	}
	if levelUpChange == nil {
		return stats
	}

	// Level up the stats a few times
	for i := 0; i < startingLevels; i++ {
		stats.LevelUp()
	}
	return stats
}

func (s *Stats) Add(a *Stats) *Stats {
	stats := &Stats{
		currentHp: s.currentHp,
		totalHp:   s.totalHp,
		strength:  s.strength,
		wisdom:    s.wisdom,
		defense:   s.defense,
		agility:   s.agility,
		cowardice: s.cowardice,
		luck:      s.luck,
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
	stats.ModifyStat(StatCowardice, a.cowardice)
	stats.ModifyStat(StatLuck, a.luck)

	return stats
}
