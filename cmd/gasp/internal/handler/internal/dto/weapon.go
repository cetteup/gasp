package dto

const (
	WeaponAssaultRifle    uint8 = 0
	WeaponAssaultGrenade  uint8 = 1
	WeaponCarbine         uint8 = 2
	WeaponLightMachineGun uint8 = 3
	WeaponSniperRifle     uint8 = 4
	WeaponPistol          uint8 = 5
	WeaponAntiTankAntiAir uint8 = 6
	WeaponSubMachineGun   uint8 = 7
	WeaponShotgun         uint8 = 8
	WeaponKnife           uint8 = 9
	WeaponDefibrillator   uint8 = 10
	WeaponExplosives      uint8 = 11
	WeaponHandGrenade     uint8 = 12
	WeaponDummy           uint8 = 13
)

var (
	// WeaponIDs Weapon ids as used by the ASP world (explosives are grouped)
	WeaponIDs = []uint8{
		WeaponAssaultRifle,
		WeaponAssaultGrenade,
		WeaponCarbine,
		WeaponLightMachineGun,
		WeaponSniperRifle,
		WeaponPistol,
		WeaponAntiTankAntiAir,
		WeaponSubMachineGun,
		WeaponShotgun,
		WeaponKnife,
		WeaponDefibrillator,
		WeaponExplosives,
		WeaponHandGrenade,
		WeaponDummy,
	}
)
