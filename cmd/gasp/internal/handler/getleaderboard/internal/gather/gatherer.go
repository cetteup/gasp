package gather

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cetteup/gasp/cmd/gasp/internal/handler/internal/dto"
	"github.com/cetteup/gasp/internal/domain/kit"
	"github.com/cetteup/gasp/internal/domain/leaderboard"
	"github.com/cetteup/gasp/internal/domain/vehicle"
	"github.com/cetteup/gasp/internal/domain/weapon"
	"github.com/cetteup/gasp/internal/util"
)

const (
	typeScore      = "score"
	typeKit        = "kit"
	typeVehicle    = "vehicle"
	typeWeapon     = "weapon"
	typeRisingStar = "risingstar"

	scoreIDOverall = "overall"
	scoreIDCombat  = "combat"
	scoreIDCommand = "commander"
	scoreIDTeam    = "team"
)

var (
	ErrInvalidLeaderboardType = errors.New("invalid leaderboard type")
	ErrInvalidLeaderboardID   = errors.New("invalid leaderboard id")
)

type GatheredData struct {
	Keys    []string
	Entries []map[string]string
	Size    int
	AsOf    uint32
}

type Gatherer struct {
	leaderboardRepository leaderboard.Repository
}

func NewGatherer(leaderboardRepository leaderboard.Repository) *Gatherer {
	return &Gatherer{
		leaderboardRepository: leaderboardRepository,
	}
}

func (g *Gatherer) Gather(ctx context.Context, t, id string, position, before, after uint32, pid *uint32) (GatheredData, error) {
	var filter leaderboard.Filter
	if pid != nil {
		filter = leaderboard.NewPIDFilter(*pid)
	} else {
		first := uint32(max(int(position)-int(before)-1, 0))
		// Limit number of returned entries to 20
		last := min(position+after, first+20)
		filter = leaderboard.NewPositionFilter(first, last)
	}

	switch t {
	case typeScore:
		return g.gatherScoreData(ctx, id, filter)
	case typeKit:
		return g.gatherKitData(ctx, id, filter)
	case typeVehicle:
		return g.gatherVehicleData(ctx, id, filter)
	case typeWeapon:
		return g.gatherWeaponData(ctx, id, filter)
	case typeRisingStar:
		return g.gatherRisingStarData(ctx, filter)
	default:
		return GatheredData{}, ErrInvalidLeaderboardType
	}
}

func (g *Gatherer) gatherScoreData(ctx context.Context, id string, filter leaderboard.Filter) (GatheredData, error) {
	scoreType, err := toScoreType(id)
	if err != nil {
		return GatheredData{}, err
	}

	entries, size, err := g.leaderboardRepository.FindTopPlayersByScore(ctx, scoreType, filter)
	if err != nil {
		return GatheredData{}, err
	}

	// Even though all ids are some type of score, the leaderboards all have different keys
	var keys []string
	switch scoreType {
	case leaderboard.ScoreTypeOverall:
		keys = []string{"n", "pid", "nick", "score", "totaltime", "playerrank", "countrycode"}
	case leaderboard.ScoreTypeCommand:
		keys = []string{"n", "pid", "nick", "coscore", "cotime", "playerrank", "countrycode"}
	case leaderboard.ScoreTypeTeam:
		keys = []string{"n", "pid", "nick", "teamscore", "totaltime", "playerrank", "countrycode"}
	case leaderboard.ScoreTypeCombat:
		keys = []string{"n", "pid", "nick", "score", "totalkills", "totaltime", "playerrank", "countrycode"}
	}

	formatted := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		// Always add all possibly required values to the map (avoid another switch or similar)
		formatted = append(formatted, map[string]string{
			"n":           util.FormatUint(entry.Position),
			"pid":         util.FormatUint(entry.Data.ID),
			"nick":        entry.Data.Name,
			"score":       util.FormatInt(entry.Data.Score),
			"coscore":     util.FormatInt(entry.Data.CommandScore),
			"teamscore":   util.FormatInt(entry.Data.TeamScore),
			"totalkills":  util.FormatUint(entry.Data.Kills),
			"totaltime":   util.FormatUint(entry.Data.Time),
			"cotime":      util.FormatUint(entry.Data.CommandTime),
			"playerrank":  util.FormatUint(entry.Data.Rank.ID),
			"countrycode": strings.ToUpper(entry.Data.Country),
		})
	}

	return GatheredData{
		Keys:    keys,
		Entries: formatted,
		Size:    size,
		// Will overflow on 7 February 2106 at 06:28:15 UTC
		AsOf: uint32(time.Now().UTC().Unix()),
	}, nil
}

func (g *Gatherer) gatherKitData(ctx context.Context, id string, filter leaderboard.Filter) (GatheredData, error) {
	kitID, err := toKitID(id)
	if err != nil {
		return GatheredData{}, err
	}

	entries, size, err := g.leaderboardRepository.FindTopPlayersByKit(ctx, kitID, filter)
	if err != nil {
		return GatheredData{}, err
	}

	// "deathsby" was spelled correctly in original responses
	// See https://web.archive.org/web/20070523064955/http://bf2web.gamespy.com:80/ASP/getleaderboard.aspx?type=kit&id=0&pos=10&before=9&after=10&debug=tx&nocache=633148776697742028
	keys := []string{"n", "pid", "nick", "killswith", "deathsby", "timeused", "playerrank", "countrycode"}
	formatted := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		formatted = append(formatted, map[string]string{
			"n":         util.FormatUint(entry.Position),
			"pid":       util.FormatUint(entry.Data.Player.ID),
			"nick":      entry.Data.Player.Name,
			"killswith": util.FormatUint(entry.Data.Kills),
			// "deathsby" is NOT "deaths", but the actual value is currently not available, so just return zero
			// See https://github.com/startersclan/asp/issues/86 for details
			"deathsby":    "0",
			"timeused":    util.FormatUint(entry.Data.Time),
			"playerrank":  util.FormatUint(entry.Data.Player.Rank.ID),
			"countrycode": strings.ToUpper(entry.Data.Player.Country),
		})
	}

	return GatheredData{
		Keys:    keys,
		Entries: formatted,
		Size:    size,
		// Will overflow on 7 February 2106 at 06:28:15 UTC
		AsOf: uint32(time.Now().UTC().Unix()),
	}, nil
}

func (g *Gatherer) gatherVehicleData(ctx context.Context, id string, filter leaderboard.Filter) (GatheredData, error) {
	vehicleID, err := toVehicleID(id)
	if err != nil {
		return GatheredData{}, err
	}

	entries, size, err := g.leaderboardRepository.FindTopPlayersByVehicle(ctx, vehicleID, filter)
	if err != nil {
		return GatheredData{}, err
	}

	// Yes, "deathsby" was actually misspelled "detahsby" in the original responses (and never fixed it seems)
	// See  https://web.archive.org/web/20070523065223/http://bf2web.gamespy.com:80/ASP/getleaderboard.aspx?type=vehicle&id=0&pos=10&before=9&after=10&debug=tx&nocache=633148776697742028
	keys := []string{"n", "pid", "nick", "killswith", "detahsby", "timeused", "playerrank", "countrycode"}
	formatted := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		formatted = append(formatted, map[string]string{
			"n":         util.FormatUint(entry.Position),
			"pid":       util.FormatUint(entry.Data.Player.ID),
			"nick":      entry.Data.Player.Name,
			"killswith": util.FormatUint(entry.Data.Kills),
			// "deathsby" is NOT "deaths", but the actual value is currently not available, so just return zero
			// See https://github.com/startersclan/asp/issues/86 for details
			"detahsby":    "0",
			"timeused":    util.FormatUint(entry.Data.Time),
			"playerrank":  util.FormatUint(entry.Data.Player.Rank.ID),
			"countrycode": strings.ToUpper(entry.Data.Player.Country),
		})
	}

	return GatheredData{
		Keys:    keys,
		Entries: formatted,
		Size:    size,
		// Will overflow on 7 February 2106 at 06:28:15 UTC
		AsOf: uint32(time.Now().UTC().Unix()),
	}, nil
}

func (g *Gatherer) gatherWeaponData(ctx context.Context, id string, filter leaderboard.Filter) (GatheredData, error) {
	weaponID, err := toWeaponID(id)
	if err != nil {
		return GatheredData{}, err
	}

	entries, size, err := g.leaderboardRepository.FindTopPlayersByWeapon(ctx, weaponID, filter)
	if err != nil {
		return GatheredData{}, err
	}

	// Yes, "deathsby" was actually misspelled "detahsby" in the original responses (and never fixed it seems)
	// See https://web.archive.org/web/20081003204115/http://bf2web.gamespy.com:80/ASP/getleaderboard.aspx?type=weapon&id=6&pos=10&before=9&after=10&debug=tx&nocache=633585835850837712
	keys := []string{"n", "pid", "nick", "killswith", "detahsby", "timeused", "accuracy", "playerrank", "countrycode"}
	formatted := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		formatted = append(formatted, map[string]string{
			"n":         util.FormatUint(entry.Position),
			"pid":       util.FormatUint(entry.Data.Player.ID),
			"nick":      entry.Data.Player.Name,
			"killswith": util.FormatUint(entry.Data.Kills),
			// "deathsby" is NOT "deaths", but the actual value is currently not available, so just return zero
			// See https://github.com/startersclan/asp/issues/86 for details
			"detahsby":    "0",
			"timeused":    util.FormatUint(entry.Data.Time),
			"accuracy":    util.FormatUint(util.DivideUint(entry.Data.ShotsHit*100, entry.Data.ShotsFired)),
			"playerrank":  util.FormatUint(entry.Data.Player.Rank.ID),
			"countrycode": strings.ToUpper(entry.Data.Player.Country),
		})
	}

	return GatheredData{
		Keys:    keys,
		Entries: formatted,
		Size:    size,
		// Will overflow on 7 February 2106 at 06:28:15 UTC
		AsOf: uint32(time.Now().UTC().Unix()),
	}, nil
}

func (g *Gatherer) gatherRisingStarData(ctx context.Context, filter leaderboard.Filter) (GatheredData, error) {
	entries, size, err := g.leaderboardRepository.FindRisingStars(ctx, filter)
	if err != nil {
		return GatheredData{}, err
	}

	keys := []string{"n", "pid", "nick", "weeklyscore", "totaltime", "date", "playerrank", "countrycode"}
	formatted := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		joined := time.Unix(int64(entry.Data.Player.Joined), 0)
		formatted = append(formatted, map[string]string{
			"n":           util.FormatUint(entry.Position),
			"pid":         util.FormatUint(entry.Data.Player.ID),
			"nick":        entry.Data.Player.Name,
			"weeklyscore": util.FormatFloat(float64(entry.Data.WeeklyScore) / 10000),
			"totaltime":   util.FormatUint(entry.Data.Player.Time),
			"date":        joined.UTC().Format("01/02/06 03:04:00 PM"), // "rounded" to minutes
			"playerrank":  util.FormatUint(entry.Data.Player.Rank.ID),
			"countrycode": strings.ToUpper(entry.Data.Player.Country),
		})
	}

	// While all other leaderboards are always "asof" now, the rising star leaderboard is only updated manually
	// due to how heavy of a computation it is.
	timestamp, err := g.leaderboardRepository.GetRisingStarUpdateTimestamp(ctx)
	if err != nil {
		return GatheredData{}, err
	}

	return GatheredData{
		Keys:    keys,
		Entries: formatted,
		Size:    size,
		AsOf:    timestamp,
	}, nil
}

func toScoreType(id string) (leaderboard.ScoreType, error) {
	switch id {
	case scoreIDOverall:
		return leaderboard.ScoreTypeOverall, nil
	case scoreIDCombat:
		return leaderboard.ScoreTypeCombat, nil
	case scoreIDCommand:
		return leaderboard.ScoreTypeCommand, nil
	case scoreIDTeam:
		return leaderboard.ScoreTypeTeam, nil
	default:
		return 0, ErrInvalidLeaderboardID
	}
}

func toKitID(id string) (uint8, error) {
	i, err := strconv.ParseUint(id, 10, 8)
	if err != nil {
		return 0, ErrInvalidLeaderboardID
	}

	// We parse with bit size, so casting from uint64 to uint8 is safe here
	switch uint8(i) {
	case dto.KitAntiTank:
		return kit.AntiTank, nil
	case dto.KitAssault:
		return kit.Assault, nil
	case dto.KitEngineer:
		return kit.Engineer, nil
	case dto.KitMedic:
		return kit.Medic, nil
	case dto.KitSpecOps:
		return kit.SpecOps, nil
	case dto.KitSupport:
		return kit.Support, nil
	case dto.KitSniper:
		return kit.Sniper, nil
	default:
		return 0, ErrInvalidLeaderboardID
	}
}

func toVehicleID(id string) (uint8, error) {
	i, err := strconv.ParseUint(id, 10, 8)
	if err != nil {
		return 0, ErrInvalidLeaderboardID
	}

	// We parse with bit size, so casting from uint64 to uint8 is safe here
	switch uint8(i) {
	case dto.VehicleArmor:
		return vehicle.Armor, nil
	case dto.VehicleJet:
		return vehicle.Jet, nil
	case dto.VehicleAirDefense:
		return vehicle.AirDefense, nil
	case dto.VehicleHelicopter:
		return vehicle.Helicopter, nil
	case dto.VehicleTransport:
		return vehicle.Transport, nil
	case dto.VehicleGroundDefense:
		return vehicle.GroundDefense, nil
	default:
		// Not all vehicles are supported
		return 0, ErrInvalidLeaderboardID
	}
}

func toWeaponID(id string) (uint8, error) {
	i, err := strconv.ParseUint(id, 10, 8)
	if err != nil {
		return 0, ErrInvalidLeaderboardID
	}

	// We parse with bit size, so casting from uint64 to uint8 is safe here
	switch uint8(i) {
	case dto.WeaponAssaultRifle:
		return weapon.AssaultRifle, nil
	case dto.WeaponAssaultGrenade:
		return weapon.AssaultGrenade, nil
	case dto.WeaponCarbine:
		return weapon.Carbine, nil
	case dto.WeaponLightMachineGun:
		return weapon.LightMachineGun, nil
	case dto.WeaponSniperRifle:
		return weapon.SniperRifle, nil
	case dto.WeaponPistol:
		return weapon.Pistol, nil
	case dto.WeaponAntiTankAntiAir:
		return weapon.AntiTankAntiAir, nil
	case dto.WeaponSubMachineGun:
		return weapon.SubMachineGun, nil
	case dto.WeaponShotgun:
		return weapon.Shotgun, nil
	default:
		// Not all weapons are supported
		return 0, ErrInvalidLeaderboardID
	}
}
