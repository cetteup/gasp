package dto

const (
	// EquipmentTactical Flashbangs *and* teargas
	EquipmentTactical     uint8 = 6
	EquipmentGraplingHook uint8 = 7
	EquipmentZipline      uint8 = 8
)

var (
	EquipmentIDs = []uint8{
		EquipmentTactical,
		EquipmentGraplingHook,
		EquipmentZipline,
	}
)
