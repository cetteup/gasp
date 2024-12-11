package kill

type RelationType int

const (
	RelationTypeVictim RelationType = iota
	RelationTypeAttacker
)

type HistoryRecord struct {
	Player PlayerRef
	Other  PlayerStub
	Kills  uint16
	// Relation of the other player to the player.
	// Read as [other] is top [relation type] of [player].
	RelationType RelationType
}

type PlayerStub struct {
	ID     uint32
	Name   string
	RankID uint32
}

type PlayerRef struct {
	ID uint32
}
