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
		return "Boss Slimer"
	case EnemyBossSkelly:
		return "Boss Undeath"
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
	case EnemySlime:
		return "bossslime"
	case EnemySkelly:
		return "bossskull"
	case EnemyBossEbi:
		return "bossebi"
	default:
		return ""
	}
}

func (e EnemyKind) Stats() *Stats {
	switch e {
	case EnemyRat:
		return &Stats{strength: 3, defense: 3, totalHp: 30}
	case EnemySlime:
		return &Stats{strength: 6, defense: 6, totalHp: 50}
	case EnemySkelly:
		return &Stats{strength: 9, defense: 9, totalHp: 100}
	case EnemyEbi:
		return &Stats{strength: 12, defense: 12, totalHp: 200}
	case EnemyBossRat:
		return &Stats{strength: 15, defense: 25, totalHp: 1000}
	case EnemyBossSlime:
		return &Stats{strength: 30, defense: 50, totalHp: 1500}
	case EnemyBossSkelly:
		return &Stats{strength: 30, defense: 35, totalHp: 2000} // level two by default so double these
	case EnemyBossEbi:
		return &Stats{strength: 40, defense: 40, totalHp: 2000} // level three by default so triple these
	default:
		return &Stats{strength: 1, defense: 0, totalHp: 1}
	}
}

const ENEMY_SCALE = 1.0

type Enemy struct {
	name  EnemyKind
	stack *render.Stack
	stats *Stats
}

func NewEnemy(name EnemyKind, level int, stack *render.Stack) *Enemy {
	level = max(1, level/4)

	stats := NewStats(name.Stats(), true)
	for i := 0; i < level; i++ {
		stats.LevelUp(true)
	}

	// Modify stats by stat scale
	stats.strength = int(float64(stats.strength) * ENEMY_SCALE)
	stats.defense = int(float64(stats.defense) * ENEMY_SCALE)
	stats.luck = int(float64(stats.luck) * ENEMY_SCALE)
	stats.totalHp = int(float64(stats.totalHp) * ENEMY_SCALE)
	stats.currentHp = stats.totalHp

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
	return min(1, e.stats.level) * 10
}

// Random gold multiplier between 0.5 and 1.25
func (e *Enemy) Gold() int {
	randMultiplier := 0.5 + rand.Float64()
	return int(float64(e.stats.totalHp) * randMultiplier)
}

func (e *Enemy) IsDead() bool {
	return e.stats.currentHp <= 0
}

// Hit target with highest confidence
func (e *Enemy) GetTarget(dudes []*Dude) *Dude {
	if len(dudes) == 0 {
		return nil
	}

	highestConfidence := 0
	var target *Dude
	for _, d := range dudes {
		stats := d.GetCalculatedStats()
		if stats.confidence >= highestConfidence && !d.IsDead() {
			highestConfidence = d.stats.confidence
			target = d
		}
	}
	return target
}
