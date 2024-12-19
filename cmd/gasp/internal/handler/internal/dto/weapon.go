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
	// WeaponDummy This might be an equipments group (non-lethal, since nobody has any kills with this)
	// Some sources use it as zipline only, but some values don't quite line up when looking at stats from BF2Hub.
	// See https://ancientdev.com/bf2tech/bf2tech.org/index.php/BF2_Statistics.html#Function:_getplayerinfo
	WeaponDummy uint8 = 13
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
