package getplayerinfo

import (
	"fmt"
	"strings"
)

type DataSource int

const (
	DataSourceUnknown DataSource = iota
	DataSourcePlayer
	DataSourceArmyRecords
	DataSourceFieldRecords
	DataSourceKillHistoryRecords
	DataSourceKitRecords
	DataSourceVehicleRecords
	DataSourceWeaponRecords
)

var keyToSource = map[string]DataSource{
	keyID:              DataSourcePlayer,
	keyName:            DataSourcePlayer,
	keyScore:           DataSourcePlayer,
	keyJoined:          DataSourcePlayer,
	keyWins:            DataSourcePlayer,
	keyLosses:          DataSourcePlayer,
	keyMode0:           DataSourcePlayer,
	keyMode1:           DataSourcePlayer,
	keyMode2:           DataSourcePlayer,
	keyTime:            DataSourcePlayer,
	keySMOC:            DataSourcePlayer,
	keyCombatScore:     DataSourcePlayer,
	keyKills:           DataSourcePlayer,
	keyDamageAssists:   DataSourcePlayer,
	keyDeaths:          DataSourcePlayer,
	keySuicides:        DataSourcePlayer,
	keyKillStreak:      DataSourcePlayer,
	keyDeathStreak:     DataSourcePlayer,
	keyKillsPerMinute:  DataSourcePlayer,
	keyDeathsPerMinute: DataSourcePlayer,
	keyScorePreMinute:  DataSourcePlayer,
	keyKillsPerRound:   DataSourcePlayer,
	keyDeathsPerRound:  DataSourcePlayer,
	keyTeamScore:       DataSourcePlayer,
	keyCaptures:        DataSourcePlayer,
	keyCaptureAssists:  DataSourcePlayer,
	keyDefends:         DataSourcePlayer,
	keyHeals:           DataSourcePlayer,
	keyRevives:         DataSourcePlayer,
	keyResupplies:      DataSourcePlayer,
	keyRepairs:         DataSourcePlayer,
	keyTargetAssists:   DataSourcePlayer,
	keyDriverAssists:   DataSourcePlayer,
	keyDriverSpecials:  DataSourcePlayer,
	keyCommandScore:    DataSourcePlayer,
	keyRankID:          DataSourcePlayer,
	keyKicks:           DataSourcePlayer,
	keyBestScore:       DataSourcePlayer,
	keyCommandTime:     DataSourcePlayer,
	keyBans:            DataSourcePlayer,
	keyLastOnline:      DataSourcePlayer,
	keySquadLeaderTime: DataSourcePlayer,
	keySquadMemberTime: DataSourcePlayer,
	keyLoneWolfTime:    DataSourcePlayer,

	keyFavoriteField: DataSourceFieldRecords,

	keyFavoriteKit: DataSourceKitRecords,

	keyTopVictimID:      DataSourceKillHistoryRecords,
	keyTopOpponentID:    DataSourceKillHistoryRecords,
	keyTopVictimKills:   DataSourceKillHistoryRecords,
	keyTopOpponentKills: DataSourceKillHistoryRecords,
	keyTopVictimName:    DataSourceKillHistoryRecords,
	keyTopVictimRank:    DataSourceKillHistoryRecords,
	keyTopOpponentName:  DataSourceKillHistoryRecords,
	keyTopOpponentRank:  DataSourceKillHistoryRecords,

	keyRoadKills:       DataSourceVehicleRecords,
	keyFavoriteVehicle: DataSourceVehicleRecords,

	keyAccuracy:        DataSourceWeaponRecords,
	keyFavoriteWeapon:  DataSourceWeaponRecords,
	keyNightVisionTime: DataSourceWeaponRecords,
	keyGasMaskTime:     DataSourceWeaponRecords,

	groupArmyTime:           DataSourceArmyRecords,
	groupArmyWins:           DataSourceArmyRecords,
	groupArmyLosses:         DataSourceArmyRecords,
	groupArmyBestRoundScore: DataSourceArmyRecords,

	groupFieldTime:   DataSourceFieldRecords,
	groupFieldWins:   DataSourceFieldRecords,
	groupFieldLosses: DataSourceFieldRecords,

	groupKitTime:           DataSourceKitRecords,
	groupKitKills:          DataSourceKitRecords,
	groupKitDeaths:         DataSourceKitRecords,
	groupKitKillDeathRatio: DataSourceKitRecords,

	groupVehicleTime:           DataSourceVehicleRecords,
	groupVehicleKills:          DataSourceVehicleRecords,
	groupVehicleDeaths:         DataSourceVehicleRecords,
	groupVehicleKillDeathRatio: DataSourceVehicleRecords,
	groupVehicleRoadKills:      DataSourceVehicleRecords,

	groupWeaponTime:           DataSourceWeaponRecords,
	groupWeaponKills:          DataSourceWeaponRecords,
	groupWeaponDeaths:         DataSourceWeaponRecords,
	groupWeaponAccuracy:       DataSourceWeaponRecords,
	groupWeaponKillDeathRatio: DataSourceWeaponRecords,

	groupEquipmentTimesDeployed: DataSourceWeaponRecords,
}

func determineDataSources(keys []string) (map[DataSource]struct{}, error) {
	sources := make(map[DataSource]struct{}, 7)
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

func determineDataSource(key string) (DataSource, error) {
	// Split after "-" to be able to easily handle both individual and group keys
	// (eliminates the need to add every resolved group key to keyToSource)
	k := strings.SplitAfter(key, "-")
	source, ok := keyToSource[k[0]]
	if !ok {
		return DataSourceUnknown, fmt.Errorf("unknown key: %s", key)
	}
	return source, nil
}
