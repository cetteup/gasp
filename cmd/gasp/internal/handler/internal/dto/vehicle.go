package dto

const (
	VehicleArmor         uint8 = 0
	VehicleJet           uint8 = 1
	VehicleAirDefense    uint8 = 2
	VehicleHelicopter    uint8 = 3
	VehicleTransport     uint8 = 4
	VehicleArtillery     uint8 = 5
	VehicleGroundDefense uint8 = 6
)

var (
	VehicleIDs = []uint8{
		VehicleArmor,
		VehicleJet,
		VehicleAirDefense,
		VehicleHelicopter,
		VehicleTransport,
		VehicleArtillery,
		VehicleGroundDefense,
	}
)
