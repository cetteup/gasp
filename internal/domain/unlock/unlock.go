package unlock

type Unlock struct {
	ID          uint16
	Name        string
	Description string
	Kit         KitRef
}

type KitRef struct {
	ID uint8
}

type Record struct {
	Player    PlayerRef
	Unlock    Unlock
	Unlocked  bool
	Timestamp uint32 // zero if unlocked is false
}

type PlayerRef struct {
	ID uint32
}
