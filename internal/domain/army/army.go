package army

const (
	// USMC United States Marines Corps
	USMC uint8 = 0
	// MEC Middle Eastern Collation
	MEC uint8 = 1
	// PLA Peoples Liberation Army
	PLA uint8 = 2
	// NavySeals United States Navy Seals
	NavySeals uint8 = 3
	// SAS British Special Air Service
	SAS uint8 = 4
	// Spetznas Russian SPETZNAS
	Spetznas uint8 = 5
	// MECSF MEC Special Forces
	MECSF uint8 = 6
	// Rebels Rebel Forces
	Rebels uint8 = 7
	// Insurgents Insurgent Forces
	Insurgents uint8 = 8
	// EU European Union
	EU uint8 = 9
)

var (
	IDs = []uint8{
		USMC,
		MEC,
		PLA,
		NavySeals,
		SAS,
		Spetznas,
		MECSF,
		Rebels,
		Insurgents,
		EU,
	}
)

type Record struct {
	Player          PlayerRef
	Army            ArmyRef
	Time            uint32
	Wins            uint16
	Losses          uint16
	Score           int
	BestRoundScore  int16
	WorstRoundScore int16
	BestRounds      int16
}

type ArmyRef struct {
	ID uint8
}

type PlayerRef struct {
	ID uint32
}
