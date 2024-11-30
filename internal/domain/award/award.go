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
