package weapon

const (
	AssaultRifle    uint8 = 0
	AssaultGrenade  uint8 = 1
	Carbine         uint8 = 2
	LightMachineGun uint8 = 3
	SniperRifle     uint8 = 4
	Pistol          uint8 = 5
	AntiTankAntiAir uint8 = 6
	SubMachineGun   uint8 = 7
	Shotgun         uint8 = 8
	Knife           uint8 = 9
	Defibrillator   uint8 = 10
	C4              uint8 = 11
	HandGrenade     uint8 = 12
	Claymore        uint8 = 13
	AntiTankMine    uint8 = 14
	GrapplingHook   uint8 = 15
	Zipline         uint8 = 16
	Tactical        uint8 = 17
)

var (
	IDs = []uint8{
		AssaultRifle,
		AssaultGrenade,
		Carbine,
		LightMachineGun,
		SniperRifle,
		Pistol,
		AntiTankAntiAir,
		SubMachineGun,
		Shotgun,
		Knife,
		Defibrillator,
		C4,
		HandGrenade,
		Claymore,
		AntiTankMine,
		GrapplingHook,
		Zipline,
		Tactical,
	}
)

type Weapon struct {
	ID          uint8
	Name        string
	IsExplosive bool
	IsEquipment bool
}

type Record struct {
	Player PlayerRef
	// Full Weapon rather than the usual reference, since we need the weapon details for grouping weapons
	Weapon        Weapon
	Time          uint32
	Score         int
	Kills         uint32
	Deaths        uint32
	ShotsFired    uint32
	ShotsHit      uint32
	TimesDeployed uint16
}

type PlayerRef struct {
	ID uint32
}
