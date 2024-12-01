package weapon

var (
	WeaponIDs = []uint8{
		0,  // Assault Rifle
		1,  // Assault Grenade
		2,  // Carbine
		3,  // Light Machine Gun
		4,  // Sniper Rifle
		5,  // Pistol
		6,  // Anti-Tank / Anti-Air
		7,  // Sub Machine Gun
		8,  // Shotgun
		9,  // Knife
		10, // Defibrillator
		11, // C4
		12, // Hand Grenade
		13, // Claymore
		14, // Anti-Tank Mine
		15, // Grappling Hook
		16, // Zipline
		17, // Tactical
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
