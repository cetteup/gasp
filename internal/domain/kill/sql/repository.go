package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/kill"
	"github.com/cetteup/gasp/internal/sqlutil"
	"github.com/cetteup/gasp/internal/util"
)

const (
	killHistoryRecordTable = "player_kill_history"
	playerTable            = "player"

	columnAttacker = "attacker"
	columnVictim   = "victim"
	columnKills    = "count"

	columnID     = "id"
	columnName   = "name"
	columnRankID = "rank_id"

	virtualColumnPlayerID = "player_id"
	virtualColumnOtherID  = "other_id"
)

type HistoryRecordRepository struct {
	runner sq.BaseRunner
}

func NewHistoryRecordRepository(runner sq.BaseRunner) *HistoryRecordRepository {
	return &HistoryRecordRepository{
		runner: runner,
	}
}

func (r *HistoryRecordRepository) FindTopRelatedByPlayerID(ctx context.Context, playerID uint32) ([]kill.HistoryRecord, error) {
	const (
		victimDTName   = "v"
		attackerDTName = "a"
	)

	attackerUnion := sq.
		Select("*").
		FromSelect(
			sq.
				Select(
					fmt.Sprintf("%s AS %s", sqlutil.QualifyColumn(killHistoryRecordTable, columnVictim), virtualColumnPlayerID),
					fmt.Sprintf("%s AS %s", sqlutil.QualifyColumn(killHistoryRecordTable, columnAttacker), virtualColumnOtherID),
					sqlutil.QualifyColumn(playerTable, columnName),
					sqlutil.QualifyColumn(playerTable, columnRankID),
					sqlutil.QualifyColumn(killHistoryRecordTable, columnKills),
					util.FormatInt(int(kill.RelationTypeAttacker)), // Hard-set virtual type column value to victim (*other* player is attacker)
				).
				From(killHistoryRecordTable).
				InnerJoin(fmt.Sprintf(
					"%s ON %s = %s",
					playerTable,
					sqlutil.QualifyColumn(killHistoryRecordTable, columnAttacker),
					sqlutil.QualifyColumn(playerTable, columnID),
				)).
				Where(sq.Eq{sqlutil.QualifyColumn(killHistoryRecordTable, columnVictim): playerID}).
				OrderBy(fmt.Sprintf("%s DESC", sqlutil.QualifyColumn(killHistoryRecordTable, columnKills))).
				Limit(1),
			attackerDTName,
		).
		Suffix("UNION ALL")

	query := sq.
		Select("*").
		FromSelect(
			sq.
				Select(
					fmt.Sprintf("%s AS %s", sqlutil.QualifyColumn(killHistoryRecordTable, columnAttacker), virtualColumnPlayerID),
					fmt.Sprintf("%s AS %s", sqlutil.QualifyColumn(killHistoryRecordTable, columnVictim), virtualColumnOtherID),
					sqlutil.QualifyColumn(playerTable, columnName),
					sqlutil.QualifyColumn(playerTable, columnRankID),
					sqlutil.QualifyColumn(killHistoryRecordTable, columnKills),
					util.FormatInt(int(kill.RelationTypeVictim)), // Hard-set virtual type column value to victim (*other* player is victim)
				).
				From(killHistoryRecordTable).
				InnerJoin(fmt.Sprintf(
					"%s ON %s = %s",
					playerTable,
					sqlutil.QualifyColumn(killHistoryRecordTable, columnVictim),
					sqlutil.QualifyColumn(playerTable, columnID),
				)).
				Where(sq.Eq{sqlutil.QualifyColumn(killHistoryRecordTable, columnAttacker): playerID}).
				OrderBy(fmt.Sprintf("%s DESC", sqlutil.QualifyColumn(killHistoryRecordTable, columnKills))).
				Limit(1),
			victimDTName,
		).
		PrefixExpr(attackerUnion)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]kill.HistoryRecord, 0)
	for rows.Next() {
		var record kill.HistoryRecord
		if err = rows.Scan(
			&record.Player.ID,
			&record.Other.ID,
			&record.Other.Name,
			&record.Other.RankID,
			&record.Kills,
			&record.RelationType,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
