package game

import (
	"math"
	"math/rand"

	"github.com/kettek/ebijam24/internal/render"
)

type EnemyKind int

const (
	EnemyRat EnemyKind = iota
	EnemySlime
	EnemyUnknown
)

func (e EnemyKind) String() string {
	switch e {
	case EnemyRat:
		return "Rat"
	case EnemySlime:
		return "Slime"
	default:
		return "Unknown"
	}
}

func (e EnemyKind) Stats() *Stats {
	switch e {
	case EnemyRat:
		return &Stats{strength: 1, defense: 1, totalHp: 5, luck: 1}
	case EnemySlime:
		return &Stats{strength: 3, defense: 3, totalHp: 15, luck: 2}
	default:
		return &Stats{strength: 1, defense: 0, totalHp: 1, luck: 1}
	}
}

type Enemy struct {
	name  EnemyKind
	stack *render.Stack
	stats *Stats
}

func NewEnemy(name EnemyKind, level int, stack *render.Stack) *Enemy {
	if level < 1 {
		level = 1
	}
	stats := NewStats(name.Stats())
	for i := 0; i < level; i++ {
		stats.LevelUp()
	}
	return &Enemy{
		name:  name,
		stack: stack,
		stats: stats,
	}
}

func (e *Enemy) Update(d *Dude) {
	if e.stack == nil {
		return
	}
	// Face the enemy towards the dude
	e.stack.SetRotation(d.stack.Rotation() + math.Pi)

	// Position enemy slightly closer to center than the dude
	// slightly off
	cx, cy := d.stack.Position()
	distance := d.story.DistanceFromCenter(cx, cy)
	r := d.story.AngleFromCenter(cx, cy)
	nx, ny := d.story.PositionFromCenter(r, distance-15)

	e.stack.SetPosition(nx, ny)
	e.stack.Update()
}

func (e *Enemy) Draw(o render.Options) {
	e.stack.Draw(&o)
}

// Damage deals damage to the enemy, returns true if the enemy is dead.
func (e *Enemy) Damage(amount int) bool {
	e.stats.currentHp -= amount - e.stats.defense
	return e.stats.currentHp <= 0
}

func (e *Enemy) Hit() int {
	return e.stats.strength
}

func (e *Enemy) Name() string {
	return e.name.String()
}

func (e *Enemy) XP() int {
	return e.stats.level * 10
}

func (e *Enemy) Gold() float32 {
	randMultiplier := 0.5 + rand.Float32()
	return float32(e.stats.totalHp*e.stats.level) * randMultiplier
}
