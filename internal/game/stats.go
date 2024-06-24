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
	StatHP        Stat = "HP"
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

	// apply the variance to the stats
	s.totalHp += int(math.Round(float64(s.levelUpChange.totalHp) * variance()))
	s.strength += int(math.Round(float64(s.levelUpChange.strength) * variance()))
	s.wisdom += int(math.Round(float64(s.levelUpChange.wisdom) * variance()))
	s.defense += int(math.Round(float64(s.levelUpChange.defense) * variance()))
	s.agility += int(math.Round(float64(s.levelUpChange.agility) * variance()))
	s.cowardice += int(math.Round(float64(s.levelUpChange.cowardice) * variance()))

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
		s.totalHp -= s.levelUpChange.totalHp
	}
	if rand.Float64() > threshold {
		s.strength -= s.levelUpChange.strength
	}
	if rand.Float64() > threshold {
		s.wisdom -= s.levelUpChange.wisdom
	}
	if rand.Float64() > threshold {
		s.defense -= s.levelUpChange.defense
	}
	if rand.Float64() > threshold {
		s.agility -= s.levelUpChange.agility
	}
	if rand.Float64() > threshold {
		s.cowardice -= s.levelUpChange.cowardice
	}
	if rand.Float64() > threshold {
		s.luck -= s.levelUpChange.luck
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
	case StatCowardice:
		s.cowardice += amount
	case StatLuck:
		s.luck += amount
		// HP is a special case
	case StatHP:
		s.totalHp += amount
		s.currentHp += amount
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
	return &Stats{
		currentHp: s.currentHp + a.currentHp,
		totalHp:   s.totalHp + a.totalHp,
		strength:  s.strength + a.strength,
		wisdom:    s.wisdom + a.wisdom,
		defense:   s.defense + a.defense,
		agility:   s.agility + a.agility,
		cowardice: s.cowardice + a.cowardice,
		luck:      s.luck + a.luck,
	}
}
