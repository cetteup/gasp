package kit

const (
	AntiTank uint8 = 0
	Assault  uint8 = 1
	Engineer uint8 = 2
	Medic    uint8 = 3
	SpecOps  uint8 = 4
	Support  uint8 = 5
	Sniper   uint8 = 6
)

var (
	IDs = []uint8{
		AntiTank,
		Assault,
		Engineer,
		Medic,
		SpecOps,
		Support,
		Sniper,
	}
)

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
