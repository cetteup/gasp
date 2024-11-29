package unlock

type Unlock struct {
	ID          uint16
	KitID       uint8
	Name        string
	Description string
}

type Record struct {
	Player    PlayerStub // Only ID, name and rank are loaded from db (values are zero if unlocked is false)
	Unlock    Unlock
	Unlocked  bool
	Timestamp uint32 // zero if unlocked is false
}

// PlayerStub Minimal representation only containing fields relevant for unlocks
type PlayerStub struct {
	ID     uint32
	Name   string
	RankID uint8
}
