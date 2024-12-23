package info

const (
	// Player-related keys

	KeyID              = "pid"
	KeyName            = "nick"
	KeyScore           = "scor"
	KeyJoined          = "jond"
	KeyWins            = "wins"
	KeyLosses          = "loss"
	KeyMode0           = "mode0"
	KeyMode1           = "mode1"
	KeyMode2           = "mode2"
	KeyTime            = "time"
	KeySMOC            = "smoc"
	KeyCombatScore     = "cmsc"
	KeyKills           = "kill"
	KeyDamageAssists   = "kila"
	KeyDeaths          = "deth"
	KeySuicides        = "suic"
	KeyKillStreak      = "bksk"
	KeyDeathStreak     = "wdsk"
	KeyKillsPerMinute  = "klpm"
	KeyDeathsPerMinute = "dtpm"
	KeyScorePreMinute  = "ospm"
	KeyKillsPerRound   = "klpr"
	KeyDeathsPerRound  = "dtpr"
	KeyTeamScore       = "twsc"
	KeyCaptures        = "cpcp"
	KeyCaptureAssists  = "cacp"
	KeyDefends         = "dfcp"
	KeyHeals           = "heal"
	KeyRevives         = "rviv"
	KeyResupplies      = "rsup"
	KeyRepairs         = "rpar"
	KeyTargetAssists   = "tgte"
	KeyDriverAssists   = "dkas"
	KeyDriverSpecials  = "dsab"
	KeyCommandScore    = "cdsc"
	KeyRankID          = "rank"
	KeyKicks           = "kick"
	KeyBestScore       = "bbrs"
	KeyCommandTime     = "tcdr"
	KeyBans            = "ban"
	KeyLastOnline      = "lbtl"
	KeySquadLeaderTime = "tsql"
	KeySquadMemberTime = "tsqm"
	KeyLoneWolfTime    = "tlwf"

	// Map-related keys

	KeyFavoriteField = "fmap"

	// Kill-related keys

	KeyTopVictimID      = "tvcr"
	KeyTopOpponentID    = "topr"
	KeyTopVictimKills   = "mvks" // "my victim kills" as per BF2 amd64 Linux binary 0x78b4c9
	KeyTopOpponentKills = "vmks" // "victimized me kills" as per BF2 amd64 Linux binary 0x78b4f6
	KeyTopVictimName    = "mvns" // "my victim name" as per BF2 amd64 Linux binary 0x78b4b5
	KeyTopOpponentName  = "vmns" // "victimized me name" as per BF2 amd64 Linux binary 0x78b4de
	KeyTopVictimRank    = "mvrs"
	KeyTopOpponentRank  = "vmrs"

	// Kit-related keys

	KeyFavoriteKit = "fkit"

	// Vehicle-related keys

	KeyRoadKills       = "vrk"
	KeyFavoriteVehicle = "fveh"

	// Weapon-related keys

	KeyAccuracy        = "osaa"
	KeyFavoriteWeapon  = "fwea"
	KeyNightVisionTime = "tnv" // "time in night vision" as per BF2 amd64 Linux binary 0x78aad5
	KeyGasMaskTime     = "tgm" // "time in gas mask" as per BF2 amd64 Linux binary 0x78aaee

	/*
		Groups allow requesting data for a group of identical type objects (kits, vehicles, weapons etc.)
	*/

	// Army groups

	GroupArmyTime           = "atm-"
	GroupArmyWins           = "awn-"
	GroupArmyLosses         = "alo-"
	GroupArmyBestRoundScore = "abr-"

	// Field groups

	GroupFieldTime   = "mtm-"
	GroupFieldWins   = "mwn-"
	GroupFieldLosses = "mls-"

	// Kit groups

	GroupKitTime           = "ktm-"
	GroupKitKills          = "kkl-"
	GroupKitDeaths         = "kdt-"
	GroupKitKillDeathRatio = "kkd-"

	// Vehicle groups

	GroupVehicleTime           = "vtm-"
	GroupVehicleKills          = "vkl-"
	GroupVehicleDeaths         = "vdt-"
	GroupVehicleKillDeathRatio = "vkd-"
	GroupVehicleRoadKills      = "vkr-"
	GroupVehicleAccuracy       = "vac-"

	// Weapon groups

	GroupWeaponTime           = "wtm-"
	GroupWeaponKills          = "wkl-"
	GroupWeaponDeaths         = "wdt-"
	GroupWeaponAccuracy       = "wac-"
	GroupWeaponKillDeathRatio = "wkd-"

	// Equipment groups (internally tracked as weapons)

	GroupEquipmentTimesDeployed = "de-"

	/*
		Sets allow requesting data for multiple, potentially unrelated keys using a single keyword
	*/

	SetPersonalStats = "per*"
	SetCombatStats   = "cmb*"
	SetTopVictim     = "mvn*"
	SetTopOpponent   = "vmr*"
)
