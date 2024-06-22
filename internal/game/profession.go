package game

// ProfessionKind is an enumeration of the different kinds of Professions a dude can have
type ProfessionKind int

func (p *ProfessionKind) String() string {
	switch *p {
	case Vagabond:
		return "Vagabond"
	case Knight:
		return "Knight"
	case Cleric:
		return "Cleric"
	case Ranger:
		return "Ranger"
	default:
		return "Unknown"
	}
}

const (
	// Medium defense, medium attack, medium hp
	Vagabond ProfessionKind = iota

	// High defense, low attack, high hp
	Knight

	// Low defense, low attack, *can heal*
	Cleric

	// Medium defense, high attack, low hp *ranged*
	Ranger
)

// A profession defines a dude's abilities.
// It also defines the dude's appearance.
type Profession struct {
	kind              ProfessionKind
	description       string
	startingStats     Stats
	startingEquipment []*Equipment
}

func NewProfession(kind ProfessionKind, level int) *Profession {
	switch kind {
	case Knight:
		return &Profession{
			kind:          Knight,
			description:   "A knight in shining armor",
			startingStats: *getStartingStats(Knight, 1),
		}
	case Cleric:
		return &Profession{
			kind:          Cleric,
			description:   "A cleric who can heal",
			startingStats: *getStartingStats(Cleric, 1),
		}
	case Vagabond:
		return &Profession{
			kind:          Vagabond,
			description:   "A vagabond with no home",
			startingStats: *getStartingStats(Vagabond, 1),
		}
	case Ranger:
		return &Profession{
			kind:          Ranger,
			description:   "A ranger who can shoot from afar",
			startingStats: *getStartingStats(Ranger, 1),
		}
	}
	return nil
}

func (p *Profession) String() string {
	return p.kind.String()
}
func (p *Profession) Description() string {
	return p.description
}
func (p *Profession) StartingStats() Stats {
	return p.startingStats
}
func (p *Profession) StartingEquipment() []*Equipment {
	return p.startingEquipment
}

// Professions are created using their level change modifiers to stats and a given level
// Then they level up and apply the changes
func getStartingStats(kind ProfessionKind, level int) *Stats {
	switch kind {
	case Knight:
		return NewStats(&Stats{
			level:     level,
			totalHp:   15,
			strength:  2,
			wisdom:    1,
			defense:   3,
			agility:   1,
			cowardice: -10, // balls get bigger
		})
	case Cleric:
		return NewStats(&Stats{
			level:     level,
			totalHp:   7,
			strength:  1,
			wisdom:    3,
			defense:   2,
			agility:   2,
			cowardice: 10, // balls get smaller
		})
	case Vagabond:
		return NewStats(&Stats{
			level:     level,
			totalHp:   7,
			strength:  1,
			wisdom:    3,
			defense:   2,
			agility:   2,
			cowardice: -5, // no change
		})

	case Ranger:
		return NewStats(&Stats{
			level:     1,
			totalHp:   7,
			strength:  1,
			wisdom:    3,
			defense:   2,
			agility:   2,
			cowardice: -5, // no change
		})
	default:
		// you useless jobless bum
		return NewStats(nil)

	}
}
