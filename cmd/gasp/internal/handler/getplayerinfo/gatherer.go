package getplayerinfo

import (
	"context"
	"fmt"

	"github.com/cetteup/gasp/internal/constraints"
	"github.com/cetteup/gasp/internal/domain/army"
	"github.com/cetteup/gasp/internal/domain/field"
	"github.com/cetteup/gasp/internal/domain/kill"
	"github.com/cetteup/gasp/internal/domain/kit"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/vehicle"
	"github.com/cetteup/gasp/internal/domain/weapon"
	"github.com/cetteup/gasp/internal/sync"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/task"
)

type Values struct {
	Individual map[string]string
	Groups     map[string][]GroupValue
}

type GroupValue struct {
	Key   string
	Value string
}

type Gatherer struct {
	playerRepository            player.Repository
	armyRecordRepository        army.RecordRepository
	fieldRecordRepository       field.RecordRepository
	killHistoryRecordRepository kill.HistoryRecordRepository
	kitRecordRepository         kit.RecordRepository
	vehicleRecordRepository     vehicle.RecordRepository
	weaponRecordRepository      weapon.RecordRepository
}

func NewGatherer(
	playerRepository player.Repository,
	armyRecordRepository army.RecordRepository,
	fieldRecordRepository field.RecordRepository,
	killHistoryRecordRepository kill.HistoryRecordRepository,
	kitRecordRepository kit.RecordRepository,
	vehicleRecordRepository vehicle.RecordRepository,
	weaponRecordRepository weapon.RecordRepository,
) *Gatherer {
	return &Gatherer{
		playerRepository:            playerRepository,
		armyRecordRepository:        armyRecordRepository,
		fieldRecordRepository:       fieldRecordRepository,
		killHistoryRecordRepository: killHistoryRecordRepository,
		kitRecordRepository:         kitRecordRepository,
		vehicleRecordRepository:     vehicleRecordRepository,
		weaponRecordRepository:      weaponRecordRepository,
	}
}

func (g *Gatherer) Gather(ctx context.Context, pid uint32, keys []string) (map[string]string, error) {
	sources, err := determineDataSources(keys)
	if err != nil {
		return nil, err
	}

	basket := &sync.Map[string, string]{}
	var runner task.AsyncRunner
	for source := range sources {
		switch source {
		case DataSourcePlayer:
			runner.Append(g.gatherPlayerData(pid, basket))
		case DataSourceArmyRecords:
			runner.Append(g.gatherArmyRecordData(pid, basket))
		case DataSourceFieldRecords:
			runner.Append(g.gatherFieldRecordData(pid, basket))
		case DataSourceKillHistoryRecords:
			runner.Append(g.gatherKillHistoryRecordData(pid, basket))
		case DataSourceKitRecords:
			runner.Append(g.gatherKitRecordData(pid, basket))
		case DataSourceVehicleRecords:
			runner.Append(g.gatherVehicleRecordData(pid, basket))
		case DataSourceWeaponRecords:
			runner.Append(g.gatherWeaponRecordData(pid, basket))
		default:
			return nil, fmt.Errorf("unkown data source: %d", source)
		}
	}

	if err = runner.Run(ctx); err != nil {
		return nil, err
	}

	values := make(map[string]string, len(keys))
	basket.Range(func(key, value string) bool {
		values[key] = value
		return true
	})

	return values, nil
}

func (g *Gatherer) gatherPlayerData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		p, err := g.playerRepository.FindByID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find player: %w", err)
		}

		b.Store(keyID, util.FormatUint(p.ID))
		b.Store(keyName, p.Name)
		b.Store(keyScore, util.FormatInt(p.Score))
		b.Store(keyJoined, util.FormatUint(p.Joined))
		b.Store(keyWins, util.FormatUint(p.Wins))
		b.Store(keyLosses, util.FormatUint(p.Losses))
		b.Store(keyMode0, util.FormatUint(p.Mode0))
		b.Store(keyMode1, util.FormatUint(p.Mode1))
		b.Store(keyMode2, util.FormatUint(p.Mode2))
		b.Store(keyTime, util.FormatUint(p.Time))
		// TODO use constant for rank
		b.Store(keySMOC, formatBool(p.RankID == 11))
		b.Store(keyCombatScore, util.FormatInt(p.CombatScore))
		b.Store(keyKills, util.FormatUint(p.Kills))
		b.Store(keyDamageAssists, util.FormatUint(p.DamageAssists))
		b.Store(keyDeaths, util.FormatUint(p.Deaths))
		b.Store(keySuicides, util.FormatUint(p.Suicides))
		b.Store(keyKillStreak, util.FormatUint(p.KillStreak))
		b.Store(keyDeathStreak, util.FormatUint(p.DeathStreak))
		b.Store(keyKillsPerMinute, util.FormatFloat(divide(p.Kills*60, p.Time)))
		b.Store(keyDeathsPerMinute, util.FormatFloat(divide(p.Deaths*60, p.Time)))
		b.Store(keyScorePreMinute, util.FormatFloat(divide(p.Score*60, p.Time)))
		b.Store(keyKillsPerRound, util.FormatFloat(divide(p.Kills, p.Rounds)))
		b.Store(keyDeathsPerRound, util.FormatFloat(divide(p.Deaths, p.Rounds)))
		b.Store(keyTeamScore, util.FormatInt(p.TeamScore))
		b.Store(keyCaptures, util.FormatUint(p.Captures))
		b.Store(keyCaptureAssists, util.FormatUint(p.CaptureAssists))
		b.Store(keyDefends, util.FormatUint(p.Defends))
		b.Store(keyHeals, util.FormatUint(p.Heals))
		b.Store(keyRevives, util.FormatUint(p.Revives))
		b.Store(keyResupplies, util.FormatUint(p.Resupplies))
		b.Store(keyRepairs, util.FormatUint(p.Repairs))
		b.Store(keyTargetAssists, util.FormatUint(p.TargetAssists))
		b.Store(keyDriverAssists, util.FormatUint(p.DriverAssists))
		b.Store(keyDriverSpecials, util.FormatUint(p.DriverSpecials))
		b.Store(keyCommandScore, util.FormatInt(p.CommandScore))
		b.Store(keyRankID, util.FormatUint(p.RankID))
		b.Store(keyKicks, util.FormatUint(p.TimesKicked))
		b.Store(keyBestScore, util.FormatUint(p.BestScore))
		b.Store(keyCommandTime, util.FormatUint(p.CommandTime))
		b.Store(keyBans, util.FormatUint(p.TimesBanned))
		b.Store(keyLastOnline, util.FormatUint(p.LastOnline))
		b.Store(keySquadLeaderTime, util.FormatUint(p.SquadLeaderTime))
		b.Store(keySquadMemberTime, util.FormatUint(p.SquadMemberTime))
		b.Store(keyLoneWolfTime, util.FormatUint(p.LoneWolfTime))

		// Add required empty/dummy values
		b.Store(keyNightVisionTime, "0")
		b.Store(keyGasMaskTime, "0")

		return nil
	}
}

func (g *Gatherer) gatherArmyRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.armyRecordRepository.FindByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find army records: %w", err)
		}

		catalog := make(map[uint8]army.Record, len(records))
		for _, record := range records {
			catalog[record.Army.ID] = record
		}

		// Records are added "lazily", so we may only have records for some armies
		// Since we cannot leave any gaps, we fill every item with record data or zeroes
		for _, id := range armyIDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(groupArmyTime+suffix, util.FormatUint(record.Time))
			b.Store(groupArmyWins+suffix, util.FormatUint(record.Wins))
			b.Store(groupArmyLosses+suffix, util.FormatUint(record.Losses))
			b.Store(groupArmyBestRoundScore+suffix, util.FormatInt(record.BestRoundScore))
		}

		return nil
	}
}

func (g *Gatherer) gatherFieldRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.fieldRecordRepository.FindByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find field records: %w", err)
		}

		catalog := make(map[uint16]field.Record, len(records))
		for _, record := range records {
			catalog[record.Field.ID] = record
		}

		// Records are added "lazily", so we may only have records for some maps
		// Since we cannot leave any gaps, we fill every item with record data or zeroes
		var favorite field.Record
		for _, id := range fieldIDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(groupFieldTime+suffix, util.FormatUint(record.Time))
			b.Store(groupFieldWins+suffix, util.FormatUint(record.Wins))
			b.Store(groupFieldLosses+suffix, util.FormatUint(record.Losses))

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(keyFavoriteField, util.FormatUint(favorite.Field.ID))

		return nil
	}
}

func (g *Gatherer) gatherKillHistoryRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.killHistoryRecordRepository.FindTopRelatedByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find kill history records: %w", err)
		}

		// We want (and should only ever get) a maximum of 2 records
		catalog := make(map[kill.RelationType]kill.HistoryRecord, 2)
		for _, record := range records {
			// Comparing kills just to guarantee that we *can* properly handle more than 2 records
			if record.Kills > catalog[record.RelationType].Kills {
				catalog[record.RelationType] = record
			}
		}

		// We may not have records for both types in two cases:
		// - player has not played/killed/been killed yet ("lazy")
		// - top victim/opponent is an unknown player (e.g. from another provider)
		// Accessing the catalog will return a zero record in case the entry is missing - which is perfect
		victim := catalog[kill.RelationTypeVictim]
		b.Store(keyTopVictimID, util.FormatUint(victim.Other.ID))
		b.Store(keyTopVictimName, victim.Other.Name)
		b.Store(keyTopVictimRank, util.FormatUint(victim.Other.RankID))
		b.Store(keyTopVictimKills, util.FormatUint(victim.Kills))

		opponent := catalog[kill.RelationTypeAttacker]
		b.Store(keyTopOpponentID, util.FormatUint(opponent.Other.ID))
		b.Store(keyTopOpponentName, opponent.Other.Name)
		b.Store(keyTopOpponentRank, util.FormatUint(opponent.Other.RankID))
		b.Store(keyTopOpponentKills, util.FormatUint(opponent.Kills))

		return nil
	}
}

func (g *Gatherer) gatherKitRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.kitRecordRepository.FindByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find kit records: %w", err)
		}

		catalog := make(map[uint8]kit.Record, len(records))
		for _, record := range records {
			catalog[record.Kit.ID] = record
		}

		// Records are added "lazily", so we may only have records for some kits
		// Since we cannot leave any gaps, we fill every item with record data or zeroes
		var favorite kit.Record
		for _, id := range kitIDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(groupKitTime+suffix, util.FormatUint(record.Time))
			b.Store(groupKitKills+suffix, util.FormatUint(record.Kills))
			b.Store(groupKitDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(groupKitKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(keyFavoriteKit, util.FormatUint(favorite.Kit.ID))

		return nil
	}
}

func (g *Gatherer) gatherVehicleRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.vehicleRecordRepository.FindByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find vehicle records: %w", err)
		}

		catalog := make(map[uint8]vehicle.Record, len(records))
		for _, record := range records {
			catalog[record.Vehicle.ID] = record
		}

		// Records are added "lazily", so we may only have records for some vehicles
		// Since we cannot leave any gaps, we fill every item with record data or zeroes
		var favorite vehicle.Record
		// Using uint64 just to be safe, uint32 would probably be plenty
		var roadKills uint64
		for _, id := range kitIDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(groupVehicleTime+suffix, util.FormatUint(record.Time))
			b.Store(groupVehicleKills+suffix, util.FormatUint(record.Kills))
			b.Store(groupVehicleDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(groupVehicleKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
			b.Store(groupVehicleRoadKills+suffix, util.FormatUint(record.RoadKills))

			// Add road kills
			roadKills += uint64(record.RoadKills)

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(keyFavoriteVehicle, util.FormatUint(favorite.Vehicle.ID))
		b.Store(keyRoadKills, util.FormatUint(roadKills))

		return nil
	}
}

func (g *Gatherer) gatherWeaponRecordData(pid uint32, b *sync.Map[string, string]) task.Task {
	return func(ctx context.Context) error {
		records, err := g.weaponRecordRepository.FindByPlayerID(ctx, pid)
		if err != nil {
			return fmt.Errorf("failed to find vehicle records: %w", err)
		}

		catalog := make(map[uint8]weapon.Record, len(records))
		for _, record := range records {
			catalog[record.Weapon.ID] = record
		}

		// Records are added "lazily", so we may only have records for some vehicles
		// Since we cannot leave any gaps, we fill every item with record data or zeroes
		var favorite weapon.Record
		// "Virtual" record, since explosives are grouped into one in ASP domain
		explosives := weapon.Record{
			Weapon: weapon.Weapon{
				ID: 11,
			},
		}
		// Using uint64 just to be safe, uint32 would probably be plenty
		var shotsFired, shotsHit uint64
		// Looping over backend domain weapon ids here, since we need all to be able to "translate" to ASP world
		for _, id := range weapon.WeaponIDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]

			// Not using record.IsExplosive or record.IsEquipment here since we might be working based on zero-records
			// from the catalog. Also, equipment means something else in the ASP world than in the internal domain.
			if !isExplosiveID(id) && !isEquipmentID(id) {
				suffix := util.FormatUint(id)
				b.Store(groupWeaponTime+suffix, util.FormatUint(record.Time))
				b.Store(groupWeaponKills+suffix, util.FormatUint(record.Kills))
				b.Store(groupWeaponDeaths+suffix, util.FormatUint(record.Deaths))
				b.Store(groupWeaponAccuracy+suffix, util.FormatFloat(divide(record.ShotsHit, record.ShotsFired)*100))
				b.Store(groupWeaponKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
			} else if isExplosiveID(id) {
				explosives.Time += record.Time
				explosives.Kills += record.Kills
				explosives.Deaths += record.Deaths
				explosives.ShotsFired += record.ShotsFired
				explosives.ShotsHit += record.ShotsHit
			} else {
				equipmentID, err2 := weaponToEquipmentID(id)
				if err2 != nil {
					return err2
				}
				b.Store(groupEquipmentTimesDeployed+util.FormatUint(equipmentID), util.FormatUint(record.TimesDeployed))
			}

			// Add shots fired/hit
			shotsFired += uint64(record.ShotsFired)
			shotsHit += uint64(record.ShotsHit)

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		// Add cumulated explosives values and weapon 13 dummy (empty, but must be present)
		for _, record := range []weapon.Record{explosives, {
			Weapon: weapon.Weapon{
				ID: 13,
			},
		}} {
			suffix := util.FormatUint(record.Weapon.ID)
			b.Store(groupWeaponTime+suffix, util.FormatUint(record.Time))
			b.Store(groupWeaponKills+suffix, util.FormatUint(record.Kills))
			b.Store(groupWeaponDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(groupWeaponAccuracy+suffix, util.FormatFloat(divide(record.ShotsHit, record.ShotsFired)*100))
			b.Store(groupWeaponKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
		}

		b.Store(keyFavoriteWeapon, util.FormatUint(favorite.Weapon.ID))
		b.Store(keyAccuracy, util.FormatFloat(divide(shotsHit, shotsFired)*100))

		return nil
	}
}

func divide[A, B constraints.Integer](a A, b B) float64 {
	// Checking for 0 explicitly rather than using max(b, 1) since b could be negative
	if b == 0 {
		return float64(a)
	}
	return float64(a) / float64(b)
}

func gcd[A, B constraints.Integer](a A, b B) int64 {
	x, y := int64(a), int64(b)
	for y != 0 {
		r := x % y
		x = y
		y = r
	}
	return x
}

func ratio[A, B constraints.Integer](a A, b B) (int64, int64) {
	if b == 0 {
		return int64(a), int64(b)
	}

	d := gcd(a, b)
	return int64(a) / d, int64(b) / d
}

func formatRatio[A, B constraints.Integer](a A, b B) string {
	return fmt.Sprintf("%d:%d", a, b)
}

func formatBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func isExplosiveID(weaponID uint8) bool {
	return weaponID == 11 || weaponID == 13 || weaponID == 14
}

func isEquipmentID(weaponID uint8) bool {
	return weaponID >= 15 && weaponID <= 17
}

func weaponToEquipmentID(weaponID uint8) (uint8, error) {
	switch weaponID {
	case 15:
		return 7, nil
	case 16:
		return 8, nil
	case 17:
		return 6, nil
	default:
		return 0, fmt.Errorf("equipment id not known for weapon id: %d", weaponID)
	}
}
