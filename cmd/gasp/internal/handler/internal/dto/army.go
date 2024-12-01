package dto

const (
	// ArmyUSMC United States Marines Corps
	ArmyUSMC uint8 = 0
	// ArmyMEC Middle Eastern Collation
	ArmyMEC uint8 = 1
	// ArmyPLA Peoples Liberation Army
	ArmyPLA uint8 = 2
	// ArmyNavySeals United States Navy Seals
	ArmyNavySeals uint8 = 3
	// ArmySAS British Special Air Service
	ArmySAS uint8 = 4
	// ArmySpetznas Russian SPETZNAS
	ArmySpetznas uint8 = 5
	// ArmyMECSF MEC Special Forces
	ArmyMECSF uint8 = 6
	// ArmyRebels Rebel Forces
	ArmyRebels uint8 = 7
	// ArmyInsurgents Insurgent Forces
	ArmyInsurgents uint8 = 8
	// ArmyEU European Union
	ArmyEU uint8 = 9
)

var (
	ArmyIDs = []uint8{
		ArmyUSMC,
		ArmyMEC,
		ArmyPLA,
		ArmyNavySeals,
		ArmySAS,
		ArmySpetznas,
		ArmyMECSF,
		ArmyRebels,
		ArmyInsurgents,
		ArmyEU,
	}
)
