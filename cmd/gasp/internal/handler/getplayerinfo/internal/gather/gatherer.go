package gather

import (
	"context"
	"fmt"

	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getplayerinfo/internal/info"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/internal/dto"
	"github.com/cetteup/gasp/internal/constraints"
	"github.com/cetteup/gasp/internal/domain/army"
	"github.com/cetteup/gasp/internal/domain/field"
	"github.com/cetteup/gasp/internal/domain/kill"
	"github.com/cetteup/gasp/internal/domain/kit"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/rank"
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
		case dataSourcePlayer:
			runner.Append(g.gatherPlayerData(pid, basket))
		case dataSourceArmyRecords:
			runner.Append(g.gatherArmyRecordData(pid, basket))
		case dataSourceFieldRecords:
			runner.Append(g.gatherFieldRecordData(pid, basket))
		case dataSourceKillHistoryRecords:
			runner.Append(g.gatherKillHistoryRecordData(pid, basket))
		case dataSourceKitRecords:
			runner.Append(g.gatherKitRecordData(pid, basket))
		case dataSourceVehicleRecords:
			runner.Append(g.gatherVehicleRecordData(pid, basket))
		case dataSourceWeaponRecords:
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

		b.Store(info.KeyID, util.FormatUint(p.ID))
		b.Store(info.KeyName, p.Name)
		b.Store(info.KeyScore, util.FormatInt(p.Score))
		b.Store(info.KeyJoined, util.FormatUint(p.Joined))
		b.Store(info.KeyWins, util.FormatUint(p.Wins))
		b.Store(info.KeyLosses, util.FormatUint(p.Losses))
		b.Store(info.KeyMode0, util.FormatUint(p.Mode0))
		b.Store(info.KeyMode1, util.FormatUint(p.Mode1))
		b.Store(info.KeyMode2, util.FormatUint(p.Mode2))
		b.Store(info.KeyTime, util.FormatUint(p.Time))
		b.Store(info.KeySMOC, dto.FormatBool(p.Rank.ID == rank.SergeantMajorOfTheCorp))
		b.Store(info.KeyCombatScore, util.FormatInt(p.CombatScore))
		b.Store(info.KeyKills, util.FormatUint(p.Kills))
		b.Store(info.KeyDamageAssists, util.FormatUint(p.DamageAssists))
		b.Store(info.KeyDeaths, util.FormatUint(p.Deaths))
		b.Store(info.KeySuicides, util.FormatUint(p.Suicides))
		b.Store(info.KeyKillStreak, util.FormatUint(p.KillStreak))
		b.Store(info.KeyDeathStreak, util.FormatUint(p.DeathStreak))
		b.Store(info.KeyKillsPerMinute, util.FormatFloat(util.DivideFloat(p.Kills*60, p.Time)))
		b.Store(info.KeyDeathsPerMinute, util.FormatFloat(util.DivideFloat(p.Deaths*60, p.Time)))
		b.Store(info.KeyScorePreMinute, util.FormatFloat(util.DivideFloat(p.Score*60, p.Time)))
		b.Store(info.KeyKillsPerRound, util.FormatFloat(util.DivideFloat(p.Kills, p.Rounds)))
		b.Store(info.KeyDeathsPerRound, util.FormatFloat(util.DivideFloat(p.Deaths, p.Rounds)))
		b.Store(info.KeyTeamScore, util.FormatInt(p.TeamScore))
		b.Store(info.KeyCaptures, util.FormatUint(p.Captures))
		b.Store(info.KeyCaptureAssists, util.FormatUint(p.CaptureAssists))
		b.Store(info.KeyDefends, util.FormatUint(p.Defends))
		b.Store(info.KeyHeals, util.FormatUint(p.Heals))
		b.Store(info.KeyRevives, util.FormatUint(p.Revives))
		b.Store(info.KeyResupplies, util.FormatUint(p.Resupplies))
		b.Store(info.KeyRepairs, util.FormatUint(p.Repairs))
		b.Store(info.KeyTargetAssists, util.FormatUint(p.TargetAssists))
		b.Store(info.KeyDriverAssists, util.FormatUint(p.DriverAssists))
		b.Store(info.KeyDriverSpecials, util.FormatUint(p.DriverSpecials))
		b.Store(info.KeyCommandScore, util.FormatInt(p.CommandScore))
		b.Store(info.KeyRankID, util.FormatUint(p.Rank.ID))
		b.Store(info.KeyKicks, util.FormatUint(p.TimesKicked))
		b.Store(info.KeyBestScore, util.FormatUint(p.BestScore))
		b.Store(info.KeyCommandTime, util.FormatUint(p.CommandTime))
		b.Store(info.KeyBans, util.FormatUint(p.TimesBanned))
		b.Store(info.KeyLastOnline, util.FormatUint(p.LastOnline))
		b.Store(info.KeySquadLeaderTime, util.FormatUint(p.SquadLeaderTime))
		b.Store(info.KeySquadMemberTime, util.FormatUint(p.SquadMemberTime))
		b.Store(info.KeyLoneWolfTime, util.FormatUint(p.LoneWolfTime))

		// Add required empty/dummy values
		b.Store(info.KeyNightVisionTime, "0")
		b.Store(info.KeyGasMaskTime, "0")

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
		for _, id := range army.IDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(info.GroupArmyTime+suffix, util.FormatUint(record.Time))
			b.Store(info.GroupArmyWins+suffix, util.FormatUint(record.Wins))
			b.Store(info.GroupArmyLosses+suffix, util.FormatUint(record.Losses))
			b.Store(info.GroupArmyBestRoundScore+suffix, util.FormatInt(record.BestRoundScore))
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
		for _, id := range field.IDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(fromDomainFieldID(id))
			b.Store(info.GroupFieldTime+suffix, util.FormatUint(record.Time))
			b.Store(info.GroupFieldWins+suffix, util.FormatUint(record.Wins))
			b.Store(info.GroupFieldLosses+suffix, util.FormatUint(record.Losses))

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(info.KeyFavoriteField, util.FormatUint(favorite.Field.ID))

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
		b.Store(info.KeyTopVictimID, util.FormatUint(victim.Other.ID))
		b.Store(info.KeyTopVictimName, victim.Other.Name)
		b.Store(info.KeyTopVictimRank, util.FormatUint(victim.Other.RankID))
		b.Store(info.KeyTopVictimKills, util.FormatUint(victim.Kills))

		opponent := catalog[kill.RelationTypeAttacker]
		b.Store(info.KeyTopOpponentID, util.FormatUint(opponent.Other.ID))
		b.Store(info.KeyTopOpponentName, opponent.Other.Name)
		b.Store(info.KeyTopOpponentRank, util.FormatUint(opponent.Other.RankID))
		b.Store(info.KeyTopOpponentKills, util.FormatUint(opponent.Kills))

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
		for _, id := range kit.IDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(info.GroupKitTime+suffix, util.FormatUint(record.Time))
			b.Store(info.GroupKitKills+suffix, util.FormatUint(record.Kills))
			b.Store(info.GroupKitDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(info.GroupKitKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(info.KeyFavoriteKit, util.FormatUint(favorite.Kit.ID))

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
		for _, id := range vehicle.IDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]
			suffix := util.FormatUint(id)
			b.Store(info.GroupVehicleTime+suffix, util.FormatUint(record.Time))
			b.Store(info.GroupVehicleKills+suffix, util.FormatUint(record.Kills))
			b.Store(info.GroupVehicleDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(info.GroupVehicleKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
			b.Store(info.GroupVehicleRoadKills+suffix, util.FormatUint(record.RoadKills))
			b.Store(info.GroupVehicleAccuracy+suffix, "0") // Always zero

			// Add road kills
			roadKills += uint64(record.RoadKills)

			// Update favorite if needed
			if record.Time > favorite.Time {
				favorite = record
			}
		}

		b.Store(info.KeyFavoriteVehicle, util.FormatUint(favorite.Vehicle.ID))
		b.Store(info.KeyRoadKills, util.FormatUint(roadKills))

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
				ID: dto.WeaponExplosives,
			},
		}
		// Using uint64 just to be safe, uint32 would probably be plenty
		var shotsFired, shotsHit uint64
		// Looping over backend domain weapon ids here, since we need all to be able to "translate" to ASP world
		for _, id := range weapon.IDs {
			// If id is not present in catalog, we get a zero entry - which is perfect
			record := catalog[id]

			// Not using record.IsExplosive or record.IsEquipment here since we might be working based on zero-records
			// from the catalog. Also, equipment means something else in the ASP world than in the internal domain.
			if !isExplosiveID(id) && !isEquipmentID(id) {
				suffix := util.FormatUint(id)
				b.Store(info.GroupWeaponTime+suffix, util.FormatUint(record.Time))
				b.Store(info.GroupWeaponKills+suffix, util.FormatUint(record.Kills))
				b.Store(info.GroupWeaponDeaths+suffix, util.FormatUint(record.Deaths))
				b.Store(info.GroupWeaponAccuracy+suffix, util.FormatUint(util.DivideUint(record.ShotsHit, record.ShotsFired)*100))
				b.Store(info.GroupWeaponKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
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
				b.Store(info.GroupEquipmentTimesDeployed+util.FormatUint(equipmentID), util.FormatUint(record.TimesDeployed))
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
				ID: dto.WeaponDummy,
			},
		}} {
			suffix := util.FormatUint(record.Weapon.ID)
			b.Store(info.GroupWeaponTime+suffix, util.FormatUint(record.Time))
			b.Store(info.GroupWeaponKills+suffix, util.FormatUint(record.Kills))
			b.Store(info.GroupWeaponDeaths+suffix, util.FormatUint(record.Deaths))
			b.Store(info.GroupWeaponAccuracy+suffix, util.FormatUint(util.DivideUint(record.ShotsHit, record.ShotsFired)*100))
			b.Store(info.GroupWeaponKillDeathRatio+suffix, formatRatio(ratio(record.Kills, record.Deaths)))
		}

		b.Store(info.KeyFavoriteWeapon, util.FormatUint(favorite.Weapon.ID))
		b.Store(info.KeyAccuracy, util.FormatUint(util.DivideUint(shotsHit, shotsFired)*100))

		return nil
	}
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

func fromDomainFieldID(fieldID uint16) uint16 {
	switch fieldID {
	case field.OperationBluePearl:
		// Currently mismatched in the database
		return dto.FieldOperationBluePearl
	default:
		return fieldID
	}
}

func isExplosiveID(weaponID uint8) bool {
	return weaponID == weapon.C4 || weaponID == weapon.Claymore || weaponID == weapon.AntiTankMine
}

func isEquipmentID(weaponID uint8) bool {
	return weaponID == weapon.GrapplingHook || weaponID == weapon.Zipline || weaponID == weapon.Tactical
}

func weaponToEquipmentID(weaponID uint8) (uint8, error) {
	switch weaponID {
	case weapon.GrapplingHook:
		return dto.EquipmentGraplingHook, nil
	case weapon.Zipline:
		return dto.EquipmentZipline, nil
	case weapon.Tactical:
		return dto.EquipmentTactical, nil
	default:
		return 0, fmt.Errorf("equipment id not known for weapon id: %d", weaponID)
	}
}
