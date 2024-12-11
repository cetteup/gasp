package dto

const (
	FieldKubraDam             uint16 = 0
	FieldMashtuurCity         uint16 = 1
	FieldOperationCleanSweep  uint16 = 2
	FieldZatarWetlands        uint16 = 3
	FieldStrikeAtKarkand      uint16 = 4
	FieldSharqiPeninsula      uint16 = 5
	FieldGulfOfOman           uint16 = 6
	FieldOperationSmokescreen uint16 = 10
	FieldTarabaQuarry         uint16 = 11
	FieldRoadToJalalabad      uint16 = 12
	FieldDaqingOilfields      uint16 = 100
	FieldDalianPlant          uint16 = 101
	FieldDragonValley         uint16 = 102
	FieldFushePass            uint16 = 103
	FieldHinganHills          uint16 = 104
	FieldSonghuaStalemate     uint16 = 105
	FieldGreatWall            uint16 = 110
	FieldOperationBluePearl   uint16 = 120
	FieldMidnightSun          uint16 = 200
	FieldOperationRoadRage    uint16 = 201
	FieldOperationHarvest     uint16 = 202
	FieldDevilsPerch          uint16 = 300
	FieldIronGator            uint16 = 301
	FieldNightFlight          uint16 = 302
	FieldWarlord              uint16 = 303
	FieldLeviathan            uint16 = 304
	FieldMassDestruction      uint16 = 305
	FieldSurge                uint16 = 306
	FieldGhostTown            uint16 = 307
	FieldWakeIsland2007       uint16 = 601
	FieldHighwayTampa         uint16 = 602
)

var (
	// FieldIDs Ordered according to the original GameSpy ASP response
	FieldIDs = []uint16{
		FieldKubraDam,
		FieldMashtuurCity,
		FieldOperationCleanSweep,
		FieldZatarWetlands,
		FieldStrikeAtKarkand,
		FieldSharqiPeninsula,
		FieldGulfOfOman,
		FieldDaqingOilfields,
		FieldDalianPlant,
		FieldDragonValley,
		FieldFushePass,
		FieldHinganHills,
		FieldSonghuaStalemate,
		FieldWakeIsland2007,
		FieldDevilsPerch,
		FieldIronGator,
		FieldNightFlight,
		FieldWarlord,
		FieldLeviathan,
		FieldMassDestruction,
		FieldSurge,
		FieldGhostTown,
		FieldOperationSmokescreen,
		FieldTarabaQuarry,
		FieldGreatWall,
		FieldMidnightSun,
		FieldOperationRoadRage,
		FieldOperationHarvest,
		FieldRoadToJalalabad,
		// Highway Tampa may have been included in the original GameSpy response round the time patch 1.5 was released.
		// At least it does appear in the in-game BFHQ and is fully working.
		FieldHighwayTampa,
		// Operation Blue Pearl is not listed in the in-game BFHQ map list and was seemingly never added to the original
		// response keys, the game just ignores the values if present.
		FieldOperationBluePearl,
	}
)
