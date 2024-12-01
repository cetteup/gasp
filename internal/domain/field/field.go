// Package field would be called "map" if that wasn't a reserved keyword
package field

// Record Only used fields are modeled
type Record struct {
	Player PlayerRef
	Field  FieldRef
	Time   uint32
	Wins   uint16
	Losses uint16
}

type FieldRef struct {
	ID uint16
}

type PlayerRef struct {
	ID uint32
}
