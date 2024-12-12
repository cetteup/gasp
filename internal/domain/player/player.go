package player

type Player struct {
	ID                uint32
	Name              string
	Joined            uint32
	LastOnline        uint32
	Time              uint32
	Rounds            uint16
	Rank              RankRef
	Score             int64
	CommandScore      int64
	CombatScore       int64
	TeamScore         int64
	Kills             uint64
	Deaths            uint64
	Captures          uint64
	Neutralizes       uint64
	CaptureAssists    uint64
	NeutralizeAssists uint64
	Defends           uint64
	Heals             uint32
	Revives           uint32
	Resupplies        uint32
	Repairs           uint32
	DamageAssists     uint32
	TargetAssists     uint32
	DriverSpecials    uint32
	DriverAssists     uint32
	TeamKills         uint32
	TeamDamage        uint32
	TeamVehicleDamage uint32
	Suicides          uint16
	KillStreak        uint16
	DeathStreak       uint16
	CommandTime       uint32
	SquadLeaderTime   uint32
	SquadMemberTime   uint32
	LoneWolfTime      uint32
	TimeParachute     int32
	Wins              uint16
	Losses            uint16
	BestScore         uint16
	RankChanged       bool
	RankDecreased     bool
	Mode0             uint16
	Mode1             uint16
	Mode2             uint16
	TimesKicked       uint16
	TimesBanned       uint16
	PermanentlyBanned bool
}

type RankRef struct {
	ID uint8
}
