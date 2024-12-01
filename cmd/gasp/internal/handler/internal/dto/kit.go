package dto

const (
	KitAntiTank uint8 = 0
	KitAssault  uint8 = 1
	KitEngineer uint8 = 2
	KitMedic    uint8 = 3
	KitSpecOps  uint8 = 4
	KitSupport  uint8 = 5
	KitSniper   uint8 = 6
)

var (
	KitIDs = []uint8{
		KitAntiTank,
		KitAssault,
		KitEngineer,
		KitMedic,
		KitSpecOps,
		KitSupport,
		KitSniper,
	}
)
