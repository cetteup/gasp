// Package field would be called "map" if that wasn't a reserved keyword
package field

const (
	KubraDam             uint16 = 0
	MashtuurCity         uint16 = 1
	OperationCleanSweep  uint16 = 2
	ZatarWetlands        uint16 = 3
	StrikeAtKarkand      uint16 = 4
	SharqiPeninsula      uint16 = 5
	GulfOfOman           uint16 = 6
	OperationSmokescreen uint16 = 10
	TarabaQuarry         uint16 = 11
	RoadToJalalabad      uint16 = 12
	DaqingOilfields      uint16 = 100
	DalianPlant          uint16 = 101
	DragonValley         uint16 = 102
	FushePass            uint16 = 103
	HinganHills          uint16 = 104
	SonghuaStalemate     uint16 = 105
	GreatWall            uint16 = 110
	MidnightSun          uint16 = 200
	OperationRoadRage    uint16 = 201
	OperationHarvest     uint16 = 202
	DevilsPerch          uint16 = 300
	IronGator            uint16 = 301
	NightFlight          uint16 = 302
	Warlord              uint16 = 303
	Leviathan            uint16 = 304
	MassDestruction      uint16 = 305
	Surge                uint16 = 306
	GhostTown            uint16 = 307
	WakeIsland2007       uint16 = 601
	HighwayTampa         uint16 = 602
	// OperationBluePearl currently is referenced as 603 in the database, rather than 120
	OperationBluePearl uint16 = 603
)

var (
	IDs = []uint16{
		KubraDam,
		MashtuurCity,
		OperationCleanSweep,
		ZatarWetlands,
		StrikeAtKarkand,
		SharqiPeninsula,
		GulfOfOman,
		OperationSmokescreen,
		TarabaQuarry,
		RoadToJalalabad,
		DaqingOilfields,
		DalianPlant,
		DragonValley,
		FushePass,
		HinganHills,
		SonghuaStalemate,
		GreatWall,
		MidnightSun,
		OperationRoadRage,
		OperationHarvest,
		DevilsPerch,
		IronGator,
		NightFlight,
		Warlord,
		Leviathan,
		MassDestruction,
		Surge,
		GhostTown,
		WakeIsland2007,
		HighwayTampa,
		OperationBluePearl,
	}
)

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
