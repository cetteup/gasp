package vehicle

const (
	Armor         uint8 = 0
	Jet           uint8 = 1
	AirDefense    uint8 = 2
	Helicopter    uint8 = 3
	Transport     uint8 = 4
	Artillery     uint8 = 5
	GroundDefense uint8 = 6
)

var (
	IDs = []uint8{
		Armor,
		Jet,
		AirDefense,
		Helicopter,
		Transport,
		Artillery,
		GroundDefense,
	}
)

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
