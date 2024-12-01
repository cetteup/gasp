package vehicle

type Record struct {
	Player    PlayerRef
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

type PlayerRef struct {
	ID uint32
}
