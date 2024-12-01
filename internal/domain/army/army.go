package army

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
