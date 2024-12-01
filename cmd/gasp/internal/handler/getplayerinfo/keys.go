package getplayerinfo

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"strings"

	"github.com/cetteup/gasp/internal/util"
)

const (
	// Player-related keys

	keyID              = "pid"
	keyName            = "nick"
	keyScore           = "scor"
	keyJoined          = "jond"
	keyWins            = "wins"
	keyLosses          = "loss"
	keyMode0           = "mode0"
	keyMode1           = "mode1"
	keyMode2           = "mode2"
	keyTime            = "time"
	keySMOC            = "smoc"
	keyCombatScore     = "cmsc"
	keyKills           = "kill"
	keyDamageAssists   = "kila"
	keyDeaths          = "deth"
	keySuicides        = "suic"
	keyKillStreak      = "bksk"
	keyDeathStreak     = "wdsk"
	keyKillsPerMinute  = "klpm"
	keyDeathsPerMinute = "dtpm"
	keyScorePreMinute  = "ospm"
	keyKillsPerRound   = "klpr"
	keyDeathsPerRound  = "dtpr"
	keyTeamScore       = "twsc"
	keyCaptures        = "cpcp"
	keyCaptureAssists  = "cacp"
	keyDefends         = "dfcp"
	keyHeals           = "heal"
	keyRevives         = "rviv"
	keyResupplies      = "rsup"
	keyRepairs         = "rpar"
	keyTargetAssists   = "tgte"
	keyDriverAssists   = "dkas"
	keyDriverSpecials  = "dsab"
	keyCommandScore    = "cdsc"
	keyRankID          = "rank"
	keyKicks           = "kick"
	keyBestScore       = "bbrs"
	keyCommandTime     = "tcdr"
	keyBans            = "ban"
	keyLastOnline      = "lbtl"
	keySquadLeaderTime = "tsql"
	keySquadMemberTime = "tsqm"
	keyLoneWolfTime    = "tlwf"

	// Map-related keys

	keyFavoriteField = "fmap"

	// Kill-related keys

	keyTopVictimID      = "tvcr"
	keyTopOpponentID    = "topr"
	keyTopVictimKills   = "mvks"
	keyTopOpponentKills = "vmks"
	keyTopVictimName    = "mvns"
	keyTopVictimRank    = "mvrs"
	keyTopOpponentName  = "vmns"
	keyTopOpponentRank  = "vmrs"

	// Kit-related keys

	keyFavoriteKit = "fkit"

	// Vehicle-related keys

	keyRoadKills       = "vrk"
	keyFavoriteVehicle = "fveh"

	// Weapon-related keys

	keyAccuracy        = "osaa"
	keyFavoriteWeapon  = "fwea"
	keyNightVisionTime = "tnv"
	keyGasMaskTime     = "tgm"

	/*
		Groups allow requesting data for a group of identical type objects (kits, vehicles, weapons etc.)
	*/

	// Army groups

	groupArmyTime           = "atm-"
	groupArmyWins           = "awn-"
	groupArmyLosses         = "alo-"
	groupArmyBestRoundScore = "abr-"

	// Field groups

	groupFieldTime   = "mtm-"
	groupFieldWins   = "mwn-"
	groupFieldLosses = "mls-"

	// Kit groups

	groupKitTime           = "ktm-"
	groupKitKills          = "kkl-"
	groupKitDeaths         = "kdt-"
	groupKitKillDeathRatio = "kkd-"

	// Vehicle groups

	groupVehicleTime           = "vtm-"
	groupVehicleKills          = "vkl-"
	groupVehicleDeaths         = "vdt-"
	groupVehicleKillDeathRatio = "vkd-"
	groupVehicleRoadKills      = "vkr-"

	// Weapon groups

	groupWeaponTime           = "wtm-"
	groupWeaponKills          = "wkl-"
	groupWeaponDeaths         = "wdt-"
	groupWeaponAccuracy       = "wac-"
	groupWeaponKillDeathRatio = "wkd-"

	// Equipment groups (internally tracked as weapons)

	groupEquipmentTimesDeployed = "de-"

	/*
		Sets allow requesting data for multiple, potentially unrelated keys using a single keyword
	*/

	setPersonalStats = "per*"
	setCombatStats   = "cmb*"
	setTopVictim     = "mvn*"
	setTopOpponent   = "vmr*"

	separator = ","

	bfhqInfoLen              = 236
	bfhqInfoChecksum         = "04486f0fd931261089878b02b7a482f10f458317568619f5f4c253f79735cb1a"
	bfhqNightVisionTimeIndex = 58
)

var (
	armyIDs = []uint8{
		0, // United States Marines Corps
		1, // Middle Eastern Collation
		2, // Peoples Liberation Army
		3, // United States Navy Seals
		4, // British Special Air Service
		5, // Russian SPETZNAS
		6, // MEC Special Forces
		7, // Rebel Forces
		8, // Insurgent Forces
		9, // European Union
	}
	fieldIDs = []uint16{
		0,   // Kubra_Dam
		1,   // mashtuur_city
		2,   // Operation_Clean_Sweep
		3,   // zatar_wetlands
		4,   // strike_at_karkand
		5,   // Sharqi_Peninsula
		6,   // gulf_of_oman
		10,  // operationsmokescreen
		11,  // taraba_quarry
		12,  // road_to_jalalabad
		100, // daqing_oilfields
		101, // Dalian_Plant
		102, // dragon_valley
		103, // fushe_pass
		104, // hingan_hills
		105, // songhua_stalemate
		110, // greatwall
		200, // midnight_sun
		201, // OperationRoadRage
		202, // operationharvest
		300, // devils_perch
		301, // iron_gator
		302, // night_flight
		303, // warlord
		304, // leviathan
		305, // mass_destruction
		306, // surge
		307, // ghost_town
		601, // wake_island_2007
		602, // highway_tampa
		603, // operation_blue_pearl
	}
	kitIDs = []uint8{
		0, // Anti-Tank
		1, // Assault
		2, // Engineer
		3, // Medic
		4, // Special Ops
		5, // Support
		6, // Sniper
	}
	vehicleIDs = []uint8{
		0, // Armor
		1, // Aviator
		2, // Air Defense
		3, // Helicopter
		4, // Transport
		5, // Artillery
		6, // Ground Defense
	}
	// Weapon ids as used by the ASP world (explosives are grouped)
	weaponIDs = []uint8{
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
		11, // Explosives
		12, // Hand Grenade
		13, // Always empty (but required)
	}
	// Equipment ids as used by the ASP world
	equipmentIDs = []uint8{
		6, // Flashbangs and teargas
		7, // Grapling hook
		8, // Zip line
	}
)

// resolveInfoKeys Resolves the `info` query parameter value to the corresponding keys,
// handling groups and sets as needed.
func resolveInfoKeys(info string, defaults ...string) []string {
	in := strings.Split(info, separator)
	out := make([]string, 0, len(in)+len(defaults))
	out = append(out, defaults...)
	for _, key := range in {
		switch key {
		case groupArmyTime, groupArmyWins, groupArmyLosses, groupArmyBestRoundScore:
			out = append(out, buildGroupKeys(armyIDs, key)...)
		case groupFieldTime, groupFieldWins, groupFieldLosses:
			out = append(out, buildGroupKeys(fieldIDs, key)...)
		case groupKitTime, groupKitKills, groupKitDeaths, groupKitKillDeathRatio:
			out = append(out, buildGroupKeys(kitIDs, key)...)
		case groupVehicleTime, groupVehicleKills, groupVehicleDeaths, groupVehicleKillDeathRatio, groupVehicleRoadKills:
			out = append(out, buildGroupKeys(vehicleIDs, key)...)
		case groupWeaponTime, groupWeaponKills, groupWeaponDeaths, groupWeaponAccuracy, groupWeaponKillDeathRatio:
			out = append(out, buildGroupKeys(weaponIDs, key)...)
		case groupEquipmentTimesDeployed:
			out = append(out, buildGroupKeys(equipmentIDs, groupEquipmentTimesDeployed)...)
		case setPersonalStats:
			out = append(
				out,
				keyScore,
				keyJoined,
				keyWins,
				keyLosses,
				keyMode0,
				keyMode1,
				keyMode2,
				keyTime,
				keySMOC,
			)
		case setCombatStats:
			out = append(
				out,
				keyCombatScore,
				keyAccuracy,
				keyKills,
				keyDamageAssists,
				keyDeaths,
				keySuicides,
				keyKillStreak,
				keyDeathStreak,
				keyTopVictimID,
				keyTopOpponentID,
				keyKillsPerMinute,
				keyDeathsPerMinute,
				keyScorePreMinute,
				keyKillsPerRound,
				keyDeathsPerRound,
			)
		case setTopVictim:
			out = append(out, keyTopVictimName, keyTopVictimRank)
		case setTopOpponent:
			out = append(out, keyTopOpponentName, keyTopOpponentRank)
		default:
			out = append(out, key)
		}
	}

	// Even default `info` queries contain duplicates after resolving set placeholders
	// However, keys in the response are expected to be unique, so remove any duplicates
	seen := make(map[string]struct{}, len(out))
	out = slices.DeleteFunc(out, func(s string) bool {
		_, exists := seen[s]
		if exists {
			return true
		}
		seen[s] = struct{}{}
		return false
	})

	// Somehow GameSpy broke their own protocol, returning keys nobody asked for.
	// This is a clear violation of the protocol, but it is how the original backend handled it.
	// See https://web.archive.org/web/20111116174658/http://bf2web.gamespy.com:80/ASP/getplayerinfo.aspx?pid=43861616&info=per*,cmb*,twsc,cpcp,cacp,dfcp,kila,heal,rviv,rsup,rpar,tgte,dkas,dsab,cdsc,rank,cmsc,kick,kill,deth,suic,ospm,klpm,klpr,dtpr,bksk,wdsk,bbrs,tcdr,ban,dtpm,lbtl,osaa,vrk,tsql,tsqm,tlwf,mvks,vmks,mvn*,vmr*,fkit,fmap,fveh,fwea,wtm-,wkl-,wdt-,wac-,wkd-,vtm-,vkl-,vdt-,vkd-,vkr-,atm-,awn-,alo-,abr-,ktm-,kkl-,kdt-,kkd-
	// To maintain compatibility without tying these to some random key, insert them after the fact.
	// (no need to check seen here, since the keys cannot be in seen since, well, nobody requested them)
	if len(out) == bfhqInfoLen-5 && isBFHQInfo(info) {
		out = slices.Insert(out, bfhqNightVisionTimeIndex, keyNightVisionTime, keyGasMaskTime)
		out = append(out, buildGroupKeys(equipmentIDs, groupEquipmentTimesDeployed)...)
	}

	return out
}

func buildGroupKeys[T uint8 | uint16](ids []T, prefix string) []string {
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, prefix+util.FormatUint(id))
	}
	return keys
}

func isBFHQInfo(info string) bool {
	h := sha256.New()
	h.Write([]byte(info))
	return hex.EncodeToString(h.Sum(nil)) == bfhqInfoChecksum
}
