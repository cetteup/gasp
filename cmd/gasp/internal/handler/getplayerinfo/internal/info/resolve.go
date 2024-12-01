package info

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"strings"

	"github.com/cetteup/gasp/internal/util"
)

const (
	separator = ","

	bfhqInfoLen              = 236
	bfhqInfoChecksum         = "04486f0fd931261089878b02b7a482f10f458317568619f5f4c253f79735cb1a"
	bfhqNightVisionTimeIndex = 58
)

// Resolve Resolves the `info` query parameter value to the corresponding keys,
// handling groups and sets as needed.
func Resolve(info string, opts *ResolveOptions) []string {
	in := strings.Split(info, separator)
	out := make([]string, 0, len(in)+len(opts.DefaultKeys))
	out = append(out, opts.DefaultKeys...)
	for _, key := range in {
		switch key {
		case GroupArmyTime, GroupArmyWins, GroupArmyLosses, GroupArmyBestRoundScore:
			out = append(out, buildGroupKeys(opts.ArmyIDs, key)...)
		case GroupFieldTime, GroupFieldWins, GroupFieldLosses:
			out = append(out, buildGroupKeys(opts.FieldIDs, key)...)
		case GroupKitTime, GroupKitKills, GroupKitDeaths, GroupKitKillDeathRatio:
			out = append(out, buildGroupKeys(opts.KitIDs, key)...)
		case GroupVehicleTime, GroupVehicleKills, GroupVehicleDeaths, GroupVehicleKillDeathRatio, GroupVehicleRoadKills:
			out = append(out, buildGroupKeys(opts.VehicleIDs, key)...)
		case GroupWeaponTime, GroupWeaponKills, GroupWeaponDeaths, GroupWeaponAccuracy, GroupWeaponKillDeathRatio:
			out = append(out, buildGroupKeys(opts.WeaponIDs, key)...)
		case GroupEquipmentTimesDeployed:
			out = append(out, buildGroupKeys(opts.EquipmentIDs, GroupEquipmentTimesDeployed)...)
		case SetPersonalStats:
			out = append(
				out,
				KeyScore,
				KeyJoined,
				KeyWins,
				KeyLosses,
				KeyMode0,
				KeyMode1,
				KeyMode2,
				KeyTime,
				KeySMOC,
			)
		case SetCombatStats:
			out = append(
				out,
				KeyCombatScore,
				KeyAccuracy,
				KeyKills,
				KeyDamageAssists,
				KeyDeaths,
				KeySuicides,
				KeyKillStreak,
				KeyDeathStreak,
				KeyTopVictimID,
				KeyTopOpponentID,
				KeyKillsPerMinute,
				KeyDeathsPerMinute,
				KeyScorePreMinute,
				KeyKillsPerRound,
				KeyDeathsPerRound,
			)
		case SetTopVictim:
			out = append(out, KeyTopVictimName, KeyTopVictimRank)
		case SetTopOpponent:
			out = append(out, KeyTopOpponentName, KeyTopOpponentRank)
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
		out = slices.Insert(out, bfhqNightVisionTimeIndex, KeyNightVisionTime, KeyGasMaskTime)
		out = append(out, buildGroupKeys(opts.EquipmentIDs, GroupEquipmentTimesDeployed)...)
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
