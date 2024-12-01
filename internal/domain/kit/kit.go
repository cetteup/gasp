package kit

type Record struct {
	Player PlayerRef
	Kit    KitRef
	Time   uint32
	Score  int
	Kills  uint32
	Deaths uint32
}

type KitRef struct {
	ID uint8
}

type PlayerRef struct {
	ID uint32
}
