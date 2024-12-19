package leaderboard

type Filter struct {
	PID   *uint32
	First uint32
	Last  uint32
}

func NewPIDFilter(pid uint32) Filter {
	return Filter{
		PID: &pid,
	}
}

func NewPositionFilter(first, last uint32) Filter {
	return Filter{
		First: first,
		Last:  last,
	}
}

type Entry[T any] struct {
	Position uint32
	Data     T
}

type PlayerStub struct {
	ID           uint32
	Name         string
	Joined       uint32
	Country      string
	Time         uint32
	Rank         RankRef
	Score        int64
	CommandScore int64
	CombatScore  int64
	TeamScore    int64
	Kills        uint64
	CommandTime  uint32
}

type RankRef struct {
	ID uint8
}

type KitRecord struct {
	Player PlayerStub
	Kit    KitRef
	Time   uint32
	Score  int
	Kills  uint32
	Deaths uint32
}

type KitRef struct {
	ID uint8
}

type VehicleRecord struct {
	Player    PlayerStub
	Vehicle   VehicleRef
	Time      uint32
	Score     int
	Kills     uint32
	Deaths    uint32
	RoadKills uint32
}

type VehicleRef struct {
	ID uint8
}

type WeaponRecord struct {
	Player        PlayerStub
	Weapon        WeaponRef
	Time          uint32
	Score         int
	Kills         uint32
	Deaths        uint32
	ShotsFired    uint32
	ShotsHit      uint32
	TimesDeployed uint16
}

type WeaponRef struct {
	ID uint8
}

type RisingStar struct {
	Player      PlayerStub
	WeeklyScore uint32
}
