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
	EnemySkelly
	EnemyEbi
	EnemyBossRat
	EnemyBossSlime
	EnemyBossSkelly
	EnemyBossEbi
	EnemyUnknown
)

func (e EnemyKind) String() string {
	switch e {
	case EnemyRat:
		return "Rat"
	case EnemySlime:
		return "Slime"
	case EnemySkelly:
		return "Skelly"
	case EnemyEbi:
		return "Ebi"
	case EnemyBossRat:
		return "Boss Rat"
	case EnemyBossSlime:
		return "Boss Slime"
	case EnemyBossEbi:
		return "Boss Ebi"
	default:
		return "Unknown"
	}
}

func (e EnemyKind) BossStack() string {
	switch e {
	case EnemyBossRat:
		return "bossrat"
	default:
		return ""
	}
}

func (e EnemyKind) Stats() *Stats {
	switch e {
	case EnemyRat:
		return &Stats{strength: 3, defense: 3, totalHp: 15, luck: 1}
	case EnemySlime:
		return &Stats{strength: 5, defense: 5, totalHp: 30, luck: 3}
	case EnemySkelly:
		return &Stats{strength: 10, defense: 10, totalHp: 60, luck: 5}
	case EnemyEbi:
		return &Stats{strength: 15, defense: 15, totalHp: 90, luck: 7}
	case EnemyBossRat:
		return &Stats{strength: 8, defense: 10, totalHp: 150, luck: 1}
	case EnemyBossSlime:
		return &Stats{strength: 12, defense: 15, totalHp: 300, luck: 3}
	case EnemyBossSkelly:
		return &Stats{strength: 20, defense: 25, totalHp: 600, luck: 5}
	case EnemyBossEbi:
		return &Stats{strength: 30, defense: 35, totalHp: 1000, luck: 7}
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
	stats := NewStats(name.Stats(), true)
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
	if d == nil {
		e.stack.Update()
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

func (e *Enemy) RoomUpdate(r *Room) {
	if e.stack == nil {
		return
	}
	if r == nil || len(r.dudes) == 0 {
		e.stack.Update()
		return
	}
	d := e.GetTarget(r.dudes)
	if d == nil {
		e.stack.Update()
		return
	}

	// Face the enemy towards the dude
	cx, cy := d.stack.Position()
	distance := d.story.DistanceFromCenter(cx, cy)
	rot := d.story.AngleFromCenter(cx, cy)
	nx, ny := d.story.PositionFromCenter(rot, distance-15)

	e.stack.SetRotation(rot + math.Pi)
	e.stack.SetPosition(nx, ny)
	e.stack.Update()
}

func (e *Enemy) Draw(o render.Options) {
	e.stack.Draw(&o)
}

// Damage deals damage to the enemy, returns true if the enemy is dead.
func (e *Enemy) Damage(amount int) bool {

	// Apply defense reduction
	reducedDamage := e.stats.ApplyDefense(amount)

	e.stats.currentHp -= reducedDamage
	return e.stats.currentHp <= 0
}

func (e *Enemy) Hit() int {
	return e.stats.strength
}

func (e *Enemy) Name() string {
	return e.name.String()
}

func (e *Enemy) XP() int {
	return e.stats.totalHp
}

func (e *Enemy) Gold() int {
	randMultiplier := 0.5 + rand.Float64()
	return int(float64(e.stats.totalHp) * randMultiplier)
}

func (e *Enemy) IsDead() bool {
	return e.stats.currentHp <= 0
}

// Hit target with lowest cowardice
func (e *Enemy) GetTarget(dudes []*Dude) *Dude {
	if len(dudes) == 0 {
		return nil
	}

	lowestCowardice := math.MaxInt32
	var target *Dude
	for _, d := range dudes {
		stats := d.GetCalculatedStats()
		if stats.cowardice < lowestCowardice && !d.IsDead() {
			lowestCowardice = d.stats.cowardice
			target = d
		}
	}
	return target
}
