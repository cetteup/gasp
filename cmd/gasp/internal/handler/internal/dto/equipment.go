package dto

const (
	// Values referenced as such at offset 0x78b39a in the BF2 amd64 Linux binary

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
