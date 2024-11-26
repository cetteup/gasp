package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/sqlutil"
)

const (
	playerStatsTable = "player"

	columnID                = "id"
	columnName              = "name"
	columnJoined            = "joined"
	columnTime              = "time"
	columnRounds            = "rounds"
	columnRankID            = "rank_id"
	columnScore             = "score"
	columnCommandScore      = "cmdscore"
	columnCombatScore       = "skillscore"
	columnTeamScore         = "teamscore"
	columnKills             = "kills"
	columnDeaths            = "deaths"
	columnCaptures          = "captures"
	columnNeutralizes       = "neutralizes"
	columnCaptureAssists    = "captureassists"
	columnNeutralizeAssists = "neutralizeassists"
	columnDefends           = "defends"
	columnHeals             = "heals"
	columnRevives           = "revives"
	columnResupplies        = "resupplies"
	columnRepairs           = "repairs"
	columnDamageAssists     = "damageassists"
	columnTargetAssists     = "targetassists"
	columnDriverSpecials    = "driverspecials"
	columnTeamKills         = "teamkills"
	columnTeamDamage        = "teamdamage"
	columnTeamVehicleDamage = "teamvehicledamage"
	columnSuicides          = "suicides"
	columnKillStreak        = "killstreak"
	columnDeathStreak       = "deathstreak"
	columnCommandTime       = "cmdtime"
	columnSquadLeaderTime   = "sqltime"
	columnSquadMemberTime   = "sqmtime"
	columnLoneWolfTime      = "lwtime"
	columnTimeParachute     = "timepara"
	columnWins              = "wins"
	columnLosses            = "losses"
	columnBestScore         = "bestscore"
	columnMode0             = "mode0"
	columnMode1             = "mode1"
	columnMode2             = "mode2"
	columnPermanentlyBanned = "permban"

	maxResults = 20
	wildcard   = "%"
)

type Repository struct {
	runner sq.BaseRunner
}

func NewRepository(runner sq.BaseRunner) *Repository {
	return &Repository{
		runner: runner,
	}
}

func (r *Repository) FindByID(ctx context.Context, playerID uint32) (player.Player, error) {
	query := sq.
		Select(
			columnID,
			columnName,
			columnJoined,
			columnTime,
			columnRounds,
			columnRankID,
			columnScore,
			columnCommandScore,
			columnCombatScore,
			columnTeamScore,
			columnKills,
			columnDeaths,
			columnCaptures,
			columnNeutralizes,
			columnCaptureAssists,
			columnNeutralizeAssists,
			columnDefends,
			columnHeals,
			columnRevives,
			columnResupplies,
			columnRepairs,
			columnDamageAssists,
			columnTargetAssists,
			columnDriverSpecials,
			columnTeamKills,
			columnTeamDamage,
			columnTeamVehicleDamage,
			columnSuicides,
			columnKillStreak,
			columnDeathStreak,
			columnCommandTime,
			columnSquadLeaderTime,
			columnSquadMemberTime,
			columnLoneWolfTime,
			columnTimeParachute,
			columnWins,
			columnLosses,
			columnBestScore,
			columnMode0,
			columnMode1,
			columnMode2,
			columnPermanentlyBanned,
		).
		From(playerStatsTable).
		Where(sq.Eq{columnID: playerID})

	var p player.Player
	if err := query.RunWith(r.runner).QueryRowContext(ctx).Scan(
		&p.ID,
		&p.Name,
		&p.Joined,
		&p.Time,
		&p.Rounds,
		&p.RankID,
		&p.Score,
		&p.CommandScore,
		&p.CombatScore,
		&p.TeamScore,
		&p.Kills,
		&p.Deaths,
		&p.Captures,
		&p.Neutralizes,
		&p.CaptureAssists,
		&p.NeutralizeAssists,
		&p.Defends,
		&p.Heals,
		&p.Revives,
		&p.Resupplies,
		&p.Repairs,
		&p.DamageAssists,
		&p.TargetAssists,
		&p.DriverSpecials,
		&p.TeamKills,
		&p.TeamDamage,
		&p.TeamVehicleDamage,
		&p.Suicides,
		&p.KillStreak,
		&p.DeathStreak,
		&p.CommandTime,
		&p.SquadLeaderTime,
		&p.SquadMemberTime,
		&p.LoneWolfTime,
		&p.TimeParachute,
		&p.Wins,
		&p.Losses,
		&p.BestScore,
		&p.Mode0,
		&p.Mode1,
		&p.Mode2,
		&p.PermanentlyBanned,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return player.Player{}, player.ErrPlayerNotFound
		}
		return player.Player{}, err
	}

	return p, nil
}

func (r *Repository) FindWithNameMatching(ctx context.Context, name string, condition player.MatchCondition, order player.SortOrder) ([]player.Player, error) {
	query := sq.
		Select(
			columnID,
			columnName,
			columnJoined,
			columnTime,
			columnRounds,
			columnRankID,
			columnScore,
			columnCommandScore,
			columnCombatScore,
			columnTeamScore,
			columnKills,
			columnDeaths,
			columnCaptures,
			columnNeutralizes,
			columnCaptureAssists,
			columnNeutralizeAssists,
			columnDefends,
			columnHeals,
			columnRevives,
			columnResupplies,
			columnRepairs,
			columnDamageAssists,
			columnTargetAssists,
			columnDriverSpecials,
			columnTeamKills,
			columnTeamDamage,
			columnTeamVehicleDamage,
			columnSuicides,
			columnKillStreak,
			columnDeathStreak,
			columnCommandTime,
			columnSquadLeaderTime,
			columnSquadMemberTime,
			columnLoneWolfTime,
			columnTimeParachute,
			columnWins,
			columnLosses,
			columnBestScore,
			columnMode0,
			columnMode1,
			columnMode2,
			columnPermanentlyBanned,
		).
		From(playerStatsTable).
		Limit(maxResults)

	// LIKE values are parameterized and bound later, so string concatenation is build the *value* not the query
	// Thus the only thing we need to escape are placeholders (%, _) to avoid expensive pattern searches
	escaped := sqlutil.EscapeWildcards(name)
	switch condition {
	case player.MatchConditionContains:
		query = query.Where(sq.Like{columnName: wildcard + escaped + wildcard})
	case player.MatchConditionBeginsWith:
		query = query.Where(sq.Like{columnName: escaped + wildcard})
	case player.MatchConditionEndsWith:
		query = query.Where(sq.Like{columnName: wildcard + escaped})
	case player.MatchConditionEquals:
		query = query.Where(sq.Eq{columnName: escaped})
	default:
		return nil, fmt.Errorf("unknown match condition: %d", condition)
	}

	switch order {
	case player.SortOrderASC:
		query = query.OrderBy(fmt.Sprintf("%s ASC", columnName))
	case player.SortOrderDESC:
		query = query.OrderBy(fmt.Sprintf("%s DESC", columnName))
	default:
		return nil, fmt.Errorf("unkown sort order: %d", order)
	}

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	players := make([]player.Player, 0, maxResults)
	for rows.Next() {
		var p player.Player
		if err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Joined,
			&p.Time,
			&p.Rounds,
			&p.RankID,
			&p.Score,
			&p.CommandScore,
			&p.CombatScore,
			&p.TeamScore,
			&p.Kills,
			&p.Deaths,
			&p.Captures,
			&p.Neutralizes,
			&p.CaptureAssists,
			&p.NeutralizeAssists,
			&p.Defends,
			&p.Heals,
			&p.Revives,
			&p.Resupplies,
			&p.Repairs,
			&p.DamageAssists,
			&p.TargetAssists,
			&p.DriverSpecials,
			&p.TeamKills,
			&p.TeamDamage,
			&p.TeamVehicleDamage,
			&p.Suicides,
			&p.KillStreak,
			&p.DeathStreak,
			&p.CommandTime,
			&p.SquadLeaderTime,
			&p.SquadMemberTime,
			&p.LoneWolfTime,
			&p.TimeParachute,
			&p.Wins,
			&p.Losses,
			&p.BestScore,
			&p.Mode0,
			&p.Mode1,
			&p.Mode2,
			&p.PermanentlyBanned,
		); err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}
