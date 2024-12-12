package award

import (
	"github.com/cetteup/gasp/internal/domain/round"
)

type Type int

const (
	TypeRibbon = 0
	TypeBadge  = 1
	TypeMedal  = 2
)

// Award Only used fields are modeled
type Award struct {
	ID   uint32
	Type Type
}

type Record struct {
	Player PlayerRef
	Award  Award
	Round  round.Round
	Level  uint64
}

type PlayerRef struct {
	ID uint32
}

func IsKitBadge(awardID uint32) bool {
	switch awardID {
	case
		1031119, // Assault
		1031120, // Anti-tank
		1031109, // Sniper
		1031115, // Spec-Ops
		1031121, // Support
		1031105, // Engineer
		1031113: // Medic
		return true
	default:
		return false
	}
}
