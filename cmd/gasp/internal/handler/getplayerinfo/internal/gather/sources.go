package gather

import (
	"fmt"
	"strings"

	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getplayerinfo/internal/info"
)

type dataSource int

const (
	dataSourceUnknown dataSource = iota
	dataSourcePlayer
	dataSourceArmyRecords
	dataSourceFieldRecords
	dataSourceKillHistoryRecords
	dataSourceKitRecords
	dataSourceVehicleRecords
	dataSourceWeaponRecords
)

var keyToSource = map[string]dataSource{
	info.KeyID:              dataSourcePlayer,
	info.KeyName:            dataSourcePlayer,
	info.KeyScore:           dataSourcePlayer,
	info.KeyJoined:          dataSourcePlayer,
	info.KeyWins:            dataSourcePlayer,
	info.KeyLosses:          dataSourcePlayer,
	info.KeyMode0:           dataSourcePlayer,
	info.KeyMode1:           dataSourcePlayer,
	info.KeyMode2:           dataSourcePlayer,
	info.KeyTime:            dataSourcePlayer,
	info.KeySMOC:            dataSourcePlayer,
	info.KeyCombatScore:     dataSourcePlayer,
	info.KeyKills:           dataSourcePlayer,
	info.KeyDamageAssists:   dataSourcePlayer,
	info.KeyDeaths:          dataSourcePlayer,
	info.KeySuicides:        dataSourcePlayer,
	info.KeyKillStreak:      dataSourcePlayer,
	info.KeyDeathStreak:     dataSourcePlayer,
	info.KeyKillsPerMinute:  dataSourcePlayer,
	info.KeyDeathsPerMinute: dataSourcePlayer,
	info.KeyScorePreMinute:  dataSourcePlayer,
	info.KeyKillsPerRound:   dataSourcePlayer,
	info.KeyDeathsPerRound:  dataSourcePlayer,
	info.KeyTeamScore:       dataSourcePlayer,
	info.KeyCaptures:        dataSourcePlayer,
	info.KeyCaptureAssists:  dataSourcePlayer,
	info.KeyDefends:         dataSourcePlayer,
	info.KeyHeals:           dataSourcePlayer,
	info.KeyRevives:         dataSourcePlayer,
	info.KeyResupplies:      dataSourcePlayer,
	info.KeyRepairs:         dataSourcePlayer,
	info.KeyTargetAssists:   dataSourcePlayer,
	info.KeyDriverAssists:   dataSourcePlayer,
	info.KeyDriverSpecials:  dataSourcePlayer,
	info.KeyCommandScore:    dataSourcePlayer,
	info.KeyRankID:          dataSourcePlayer,
	info.KeyKicks:           dataSourcePlayer,
	info.KeyBestScore:       dataSourcePlayer,
	info.KeyCommandTime:     dataSourcePlayer,
	info.KeyBans:            dataSourcePlayer,
	info.KeyLastOnline:      dataSourcePlayer,
	info.KeySquadLeaderTime: dataSourcePlayer,
	info.KeySquadMemberTime: dataSourcePlayer,
	info.KeyLoneWolfTime:    dataSourcePlayer,

	info.KeyFavoriteField: dataSourceFieldRecords,

	info.KeyFavoriteKit: dataSourceKitRecords,

	info.KeyTopVictimID:      dataSourceKillHistoryRecords,
	info.KeyTopOpponentID:    dataSourceKillHistoryRecords,
	info.KeyTopVictimKills:   dataSourceKillHistoryRecords,
	info.KeyTopOpponentKills: dataSourceKillHistoryRecords,
	info.KeyTopVictimName:    dataSourceKillHistoryRecords,
	info.KeyTopVictimRank:    dataSourceKillHistoryRecords,
	info.KeyTopOpponentName:  dataSourceKillHistoryRecords,
	info.KeyTopOpponentRank:  dataSourceKillHistoryRecords,

	info.KeyRoadKills:       dataSourceVehicleRecords,
	info.KeyFavoriteVehicle: dataSourceVehicleRecords,

	info.KeyAccuracy:        dataSourceWeaponRecords,
	info.KeyFavoriteWeapon:  dataSourceWeaponRecords,
	info.KeyNightVisionTime: dataSourceWeaponRecords,
	info.KeyGasMaskTime:     dataSourceWeaponRecords,

	info.GroupArmyTime:           dataSourceArmyRecords,
	info.GroupArmyWins:           dataSourceArmyRecords,
	info.GroupArmyLosses:         dataSourceArmyRecords,
	info.GroupArmyBestRoundScore: dataSourceArmyRecords,

	info.GroupFieldTime:   dataSourceFieldRecords,
	info.GroupFieldWins:   dataSourceFieldRecords,
	info.GroupFieldLosses: dataSourceFieldRecords,

	info.GroupKitTime:           dataSourceKitRecords,
	info.GroupKitKills:          dataSourceKitRecords,
	info.GroupKitDeaths:         dataSourceKitRecords,
	info.GroupKitKillDeathRatio: dataSourceKitRecords,

	info.GroupVehicleTime:           dataSourceVehicleRecords,
	info.GroupVehicleKills:          dataSourceVehicleRecords,
	info.GroupVehicleDeaths:         dataSourceVehicleRecords,
	info.GroupVehicleKillDeathRatio: dataSourceVehicleRecords,
	info.GroupVehicleRoadKills:      dataSourceVehicleRecords,

	info.GroupWeaponTime:           dataSourceWeaponRecords,
	info.GroupWeaponKills:          dataSourceWeaponRecords,
	info.GroupWeaponDeaths:         dataSourceWeaponRecords,
	info.GroupWeaponAccuracy:       dataSourceWeaponRecords,
	info.GroupWeaponKillDeathRatio: dataSourceWeaponRecords,

	info.GroupEquipmentTimesDeployed: dataSourceWeaponRecords,
}

func determineDataSources(keys []string) (map[dataSource]struct{}, error) {
	sources := make(map[dataSource]struct{}, 7)
	for _, key := range keys {
		source, err := determineDataSource(key)
		if err != nil {
			return nil, err
		}
		if _, exists := sources[source]; !exists {
			sources[source] = struct{}{}
		}
	}

	return sources, nil
}

func determineDataSource(key string) (dataSource, error) {
	// Split after "-" to be able to easily handle both individual and group keys
	// (eliminates the need to add every resolved group key to keyToSource)
	k := strings.SplitAfter(key, "-")
	source, ok := keyToSource[k[0]]
	if !ok {
		return dataSourceUnknown, fmt.Errorf("unknown key: %s", key)
	}
	return source, nil
}
