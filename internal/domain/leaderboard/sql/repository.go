package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/leaderboard"
	"github.com/cetteup/gasp/internal/sqlutil"
)

const (
	playerTable        = "player"
	kitRecordTable     = "player_kit"
	vehicleRecordTable = "player_vehicle"
	weaponRecordTable  = "player_weapon"
	risingStarTable    = "risingstar"

	columnID           = "id"
	columnName         = "name"
	columnJoined       = "joined"
	columnCountry      = "country"
	columnTime         = "time"
	columnRankID       = "rank_id"
	columnScore        = "score"
	columnCommandScore = "cmdscore"
	columnCombatScore  = "skillscore"
	columnTeamScore    = "teamscore"
	columnKills        = "kills"
	columnCommandTime  = "cmdtime"

	columnKitID  = "kit_id"
	columnDeaths = "deaths"

	columnVehicleID = "vehicle_id"
	columnRoadKills = "roadkills"

	columnWeaponID      = "weapon_id"
	columnShotsFired    = "fired"
	columnShotsHit      = "hits"
	columnTimesDeployed = "deployed"

	columnPosition    = "pos"
	columnPlayerID    = "player_id"
	columnWeeklyScore = "weeklyscore"

	virtualColumnPosition = "position"

	maxResults = 10000
)

type Repository struct {
	runner sq.BaseRunner
}

func NewRepository(runner sq.BaseRunner) *Repository {
	return &Repository{
		runner: runner,
	}
}

func (r *Repository) FindTopPlayersByScore(ctx context.Context, scoreType leaderboard.ScoreType, filter leaderboard.Filter) ([]leaderboard.Entry[leaderboard.PlayerStub], int, error) {
	const cteName = "l"

	// Need to use different columns to rank/filter by
	scoreColumn, err := scoreTypeToColumn(scoreType)
	if err != nil {
		return nil, 0, err
	}

	cte := sq.
		Select(
			columnID,
			columnName,
			columnJoined,
			columnCountry,
			columnTime,
			columnRankID,
			columnScore,
			columnCommandScore,
			columnCombatScore,
			columnTeamScore,
			columnKills,
			columnCommandTime,
			buildPositionColumnExpr(scoreColumn, columnName),
		).
		From(playerTable).
		Where(sq.Gt{scoreColumn: 0}).
		Limit(maxResults)

	count, err := r.getEntryCount(ctx, cte)
	if err != nil {
		return nil, 0, err
	}

	query := sq.
		Select("*").
		FromSelect(cte, cteName)

	query = addFilter(query, columnID, filter)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]leaderboard.Entry[leaderboard.PlayerStub], 0, max(filter.Last-filter.First, 1))
	for rows.Next() {
		var entry leaderboard.Entry[leaderboard.PlayerStub]
		if err = rows.Scan(
			&entry.Data.ID,
			&entry.Data.Name,
			&entry.Data.Joined,
			&entry.Data.Country,
			&entry.Data.Time,
			&entry.Data.Rank.ID,
			&entry.Data.Score,
			&entry.Data.CommandScore,
			&entry.Data.CombatScore,
			&entry.Data.TeamScore,
			&entry.Data.Kills,
			&entry.Data.CommandTime,
			&entry.Position,
		); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

func (r *Repository) FindTopPlayersByKit(ctx context.Context, kitID uint8, filter leaderboard.Filter) ([]leaderboard.Entry[leaderboard.KitRecord], int, error) {
	const cteName = "l"

	cte := sq.
		Select(
			// Using QualifyAlias here since some column names are not unique, e.g. "kills"
			sqlutil.QualifyAlias(kitRecordTable, columnPlayerID),
			sqlutil.QualifyAlias(playerTable, columnName),
			sqlutil.QualifyAlias(playerTable, columnJoined),
			sqlutil.QualifyAlias(playerTable, columnCountry),
			sqlutil.QualifyAlias(playerTable, columnTime),
			sqlutil.QualifyAlias(playerTable, columnRankID),
			sqlutil.QualifyAlias(playerTable, columnScore),
			sqlutil.QualifyAlias(playerTable, columnCommandScore),
			sqlutil.QualifyAlias(playerTable, columnCombatScore),
			sqlutil.QualifyAlias(playerTable, columnTeamScore),
			sqlutil.QualifyAlias(playerTable, columnKills),
			sqlutil.QualifyAlias(playerTable, columnCommandTime),
			sqlutil.QualifyAlias(kitRecordTable, columnKitID),
			sqlutil.QualifyAlias(kitRecordTable, columnTime),
			sqlutil.QualifyAlias(kitRecordTable, columnScore),
			sqlutil.QualifyAlias(kitRecordTable, columnKills),
			sqlutil.QualifyAlias(kitRecordTable, columnDeaths),
			buildPositionColumnExpr(sqlutil.Qualify(kitRecordTable, columnKills), sqlutil.Qualify(playerTable, columnName)),
		).
		From(kitRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerTable,
			sqlutil.Qualify(kitRecordTable, columnPlayerID),
			sqlutil.Qualify(playerTable, columnID),
		)).
		Where(sq.And{
			sq.Eq{sqlutil.Qualify(kitRecordTable, columnKitID): kitID},
			sq.Gt{sqlutil.Qualify(kitRecordTable, columnKills): 0},
		}).
		Limit(maxResults)

	count, err := r.getEntryCount(ctx, cte)
	if err != nil {
		return nil, 0, err
	}

	query := sq.
		Select("*").
		FromSelect(cte, cteName)

	query = addFilter(query, sqlutil.Predicate(kitRecordTable, columnPlayerID), filter)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]leaderboard.Entry[leaderboard.KitRecord], 0, max(filter.Last-filter.First, 1))
	for rows.Next() {
		var entry leaderboard.Entry[leaderboard.KitRecord]
		if err = rows.Scan(
			&entry.Data.Player.ID,
			&entry.Data.Player.Name,
			&entry.Data.Player.Joined,
			&entry.Data.Player.Country,
			&entry.Data.Player.Time,
			&entry.Data.Player.Rank.ID,
			&entry.Data.Player.Score,
			&entry.Data.Player.CommandScore,
			&entry.Data.Player.CombatScore,
			&entry.Data.Player.TeamScore,
			&entry.Data.Player.Kills,
			&entry.Data.Player.CommandTime,
			&entry.Data.Kit.ID,
			&entry.Data.Time,
			&entry.Data.Score,
			&entry.Data.Kills,
			&entry.Data.Deaths,
			&entry.Position,
		); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

func (r *Repository) FindTopPlayersByVehicle(ctx context.Context, vehicleID uint8, filter leaderboard.Filter) ([]leaderboard.Entry[leaderboard.VehicleRecord], int, error) {
	const cteName = "l"

	cte := sq.
		Select(
			// Using QualifyAlias here since some column names are not unique, e.g. "kills"
			sqlutil.QualifyAlias(vehicleRecordTable, columnPlayerID),
			sqlutil.QualifyAlias(playerTable, columnName),
			sqlutil.QualifyAlias(playerTable, columnJoined),
			sqlutil.QualifyAlias(playerTable, columnCountry),
			sqlutil.QualifyAlias(playerTable, columnTime),
			sqlutil.QualifyAlias(playerTable, columnRankID),
			sqlutil.QualifyAlias(playerTable, columnScore),
			sqlutil.QualifyAlias(playerTable, columnCommandScore),
			sqlutil.QualifyAlias(playerTable, columnCombatScore),
			sqlutil.QualifyAlias(playerTable, columnTeamScore),
			sqlutil.QualifyAlias(playerTable, columnKills),
			sqlutil.QualifyAlias(playerTable, columnCommandTime),
			sqlutil.QualifyAlias(vehicleRecordTable, columnVehicleID),
			sqlutil.QualifyAlias(vehicleRecordTable, columnTime),
			sqlutil.QualifyAlias(vehicleRecordTable, columnScore),
			sqlutil.QualifyAlias(vehicleRecordTable, columnKills),
			sqlutil.QualifyAlias(vehicleRecordTable, columnDeaths),
			sqlutil.QualifyAlias(vehicleRecordTable, columnRoadKills),
			buildPositionColumnExpr(sqlutil.Qualify(vehicleRecordTable, columnKills), sqlutil.Qualify(playerTable, columnName)),
		).
		From(vehicleRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerTable,
			sqlutil.Qualify(vehicleRecordTable, columnPlayerID),
			sqlutil.Qualify(playerTable, columnID),
		)).
		Where(sq.And{
			sq.Eq{sqlutil.Qualify(vehicleRecordTable, columnVehicleID): vehicleID},
			sq.Gt{sqlutil.Qualify(vehicleRecordTable, columnKills): 0},
		}).
		Limit(maxResults)

	count, err := r.getEntryCount(ctx, cte)
	if err != nil {
		return nil, 0, err
	}

	query := sq.
		Select("*").
		FromSelect(cte, cteName)

	query = addFilter(query, sqlutil.Predicate(vehicleRecordTable, columnPlayerID), filter)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]leaderboard.Entry[leaderboard.VehicleRecord], 0, max(filter.Last-filter.First, 1))
	for rows.Next() {
		var entry leaderboard.Entry[leaderboard.VehicleRecord]
		if err = rows.Scan(
			&entry.Data.Player.ID,
			&entry.Data.Player.Name,
			&entry.Data.Player.Joined,
			&entry.Data.Player.Country,
			&entry.Data.Player.Time,
			&entry.Data.Player.Rank.ID,
			&entry.Data.Player.Score,
			&entry.Data.Player.CommandScore,
			&entry.Data.Player.CombatScore,
			&entry.Data.Player.TeamScore,
			&entry.Data.Player.Kills,
			&entry.Data.Player.CommandTime,
			&entry.Data.Vehicle.ID,
			&entry.Data.Time,
			&entry.Data.Score,
			&entry.Data.Kills,
			&entry.Data.Deaths,
			&entry.Data.RoadKills,
			&entry.Position,
		); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

func (r *Repository) FindTopPlayersByWeapon(ctx context.Context, weaponID uint8, filter leaderboard.Filter) ([]leaderboard.Entry[leaderboard.WeaponRecord], int, error) {
	const cteName = "l"

	cte := sq.
		Select(
			// Using QualifyAlias here since some column names are not unique, e.g. "kills"
			sqlutil.QualifyAlias(weaponRecordTable, columnPlayerID),
			sqlutil.QualifyAlias(playerTable, columnName),
			sqlutil.QualifyAlias(playerTable, columnJoined),
			sqlutil.QualifyAlias(playerTable, columnCountry),
			sqlutil.QualifyAlias(playerTable, columnTime),
			sqlutil.QualifyAlias(playerTable, columnRankID),
			sqlutil.QualifyAlias(playerTable, columnScore),
			sqlutil.QualifyAlias(playerTable, columnCommandScore),
			sqlutil.QualifyAlias(playerTable, columnCombatScore),
			sqlutil.QualifyAlias(playerTable, columnTeamScore),
			sqlutil.QualifyAlias(playerTable, columnKills),
			sqlutil.QualifyAlias(playerTable, columnCommandTime),
			sqlutil.QualifyAlias(weaponRecordTable, columnWeaponID),
			sqlutil.QualifyAlias(weaponRecordTable, columnTime),
			sqlutil.QualifyAlias(weaponRecordTable, columnScore),
			sqlutil.QualifyAlias(weaponRecordTable, columnKills),
			sqlutil.QualifyAlias(weaponRecordTable, columnDeaths),
			sqlutil.QualifyAlias(weaponRecordTable, columnShotsFired),
			sqlutil.QualifyAlias(weaponRecordTable, columnShotsHit),
			sqlutil.QualifyAlias(weaponRecordTable, columnTimesDeployed),
			buildPositionColumnExpr(sqlutil.Qualify(weaponRecordTable, columnKills), sqlutil.Qualify(playerTable, columnName)),
		).
		From(weaponRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerTable,
			sqlutil.Qualify(weaponRecordTable, columnPlayerID),
			sqlutil.Qualify(playerTable, columnID),
		)).
		Where(sq.And{
			sq.Eq{sqlutil.Qualify(weaponRecordTable, columnWeaponID): weaponID},
			sq.Gt{sqlutil.Qualify(weaponRecordTable, columnKills): 0},
		}).
		Limit(maxResults)

	count, err := r.getEntryCount(ctx, cte)
	if err != nil {
		return nil, 0, err
	}

	query := sq.
		Select("*").
		FromSelect(cte, cteName)

	query = addFilter(query, sqlutil.Predicate(weaponRecordTable, columnPlayerID), filter)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]leaderboard.Entry[leaderboard.WeaponRecord], 0, max(filter.Last-filter.First, 1))
	for rows.Next() {
		var entry leaderboard.Entry[leaderboard.WeaponRecord]
		if err = rows.Scan(
			&entry.Data.Player.ID,
			&entry.Data.Player.Name,
			&entry.Data.Player.Joined,
			&entry.Data.Player.Country,
			&entry.Data.Player.Time,
			&entry.Data.Player.Rank.ID,
			&entry.Data.Player.Score,
			&entry.Data.Player.CommandScore,
			&entry.Data.Player.CombatScore,
			&entry.Data.Player.TeamScore,
			&entry.Data.Player.Kills,
			&entry.Data.Player.CommandTime,
			&entry.Data.Weapon.ID,
			&entry.Data.Time,
			&entry.Data.Score,
			&entry.Data.Kills,
			&entry.Data.Deaths,
			&entry.Data.ShotsFired,
			&entry.Data.ShotsHit,
			&entry.Data.TimesDeployed,
			&entry.Position,
		); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

func (r *Repository) FindRisingStars(ctx context.Context, filter leaderboard.Filter) ([]leaderboard.Entry[leaderboard.RisingStar], int, error) {
	const cteName = "l"

	cte := sq.
		Select(
			sqlutil.Qualify(risingStarTable, columnPlayerID),
			sqlutil.Qualify(playerTable, columnName),
			sqlutil.Qualify(playerTable, columnJoined),
			sqlutil.Qualify(playerTable, columnCountry),
			sqlutil.Qualify(playerTable, columnTime),
			sqlutil.Qualify(playerTable, columnRankID),
			sqlutil.Qualify(playerTable, columnScore),
			sqlutil.Qualify(playerTable, columnCommandScore),
			sqlutil.Qualify(playerTable, columnCombatScore),
			sqlutil.Qualify(playerTable, columnTeamScore),
			sqlutil.Qualify(playerTable, columnKills),
			sqlutil.Qualify(playerTable, columnCommandTime),
			sqlutil.Qualify(risingStarTable, columnWeeklyScore),
			sqlutil.Qualify(risingStarTable, columnPosition),
		).
		From(risingStarTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerTable,
			sqlutil.Qualify(risingStarTable, columnPlayerID),
			sqlutil.Qualify(playerTable, columnID),
		)).
		Limit(maxResults)

	count, err := r.getEntryCount(ctx, cte)
	if err != nil {
		return nil, 0, err
	}

	query := sq.
		Select("*").
		FromSelect(cte, cteName)

	query = addFilter(query, columnPlayerID, filter)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]leaderboard.Entry[leaderboard.RisingStar], 0, max(filter.Last-filter.First, 1))
	for rows.Next() {
		var entry leaderboard.Entry[leaderboard.RisingStar]
		if err = rows.Scan(
			&entry.Data.Player.ID,
			&entry.Data.Player.Name,
			&entry.Data.Player.Joined,
			&entry.Data.Player.Country,
			&entry.Data.Player.Time,
			&entry.Data.Player.Rank.ID,
			&entry.Data.Player.Score,
			&entry.Data.Player.CommandScore,
			&entry.Data.Player.CombatScore,
			&entry.Data.Player.TeamScore,
			&entry.Data.Player.Kills,
			&entry.Data.Player.CommandTime,
			&entry.Data.WeeklyScore,
			&entry.Position,
		); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

func (r *Repository) GetRisingStarUpdateTimestamp(ctx context.Context) (uint32, error) {
	query := sq.
		// The table is ALTER-ed when re-building the rising star leaderboard, causing create_time to change
		Select("UNIX_TIMESTAMP(CREATE_TIME)").
		From("INFORMATION_SCHEMA.TABLES").
		Where(sq.And{
			sq.Expr("TABLE_SCHEMA = (SELECT DATABASE())"),
			sq.Eq{
				"TABLE_NAME": risingStarTable,
			},
		})

	var timestamp uint32
	if err := query.RunWith(r.runner).QueryRowContext(ctx).Scan(&timestamp); err != nil {
		return 0, err
	}

	return timestamp, nil
}

func (r *Repository) getEntryCount(ctx context.Context, cte sq.SelectBuilder) (int, error) {
	const cteName = "l"

	query := sq.
		Select("COUNT(*)").
		FromSelect(cte, cteName)

	var count int
	if err := query.RunWith(r.runner).QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func addFilter(query sq.SelectBuilder, column string, filter leaderboard.Filter) sq.SelectBuilder {
	if filter.PID != nil {
		return query.Where(sq.Eq{column: *filter.PID})
	}

	return query.
		Offset(uint64(filter.First)).
		Limit(uint64(filter.Last - filter.First))
}

func buildPositionColumnExpr(column, tiebreaker string) string {
	return fmt.Sprintf("RANK() OVER (ORDER BY %s DESC, %s ASC) AS %s", column, tiebreaker, virtualColumnPosition)
}

func scoreTypeToColumn(scoreType leaderboard.ScoreType) (string, error) {
	switch scoreType {
	case leaderboard.ScoreTypeOverall:
		return columnScore, nil
	case leaderboard.ScoreTypeCommand:
		return columnCommandScore, nil
	case leaderboard.ScoreTypeTeam:
		return columnTeamScore, nil
	case leaderboard.ScoreTypeCombat:
		return columnCombatScore, nil
	default:
		return "", fmt.Errorf("unknown score type: %d", scoreType)
	}
}
