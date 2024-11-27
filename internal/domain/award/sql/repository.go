package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/sqlutil"
)

const (
	awardRecordTable = "player_award"
	awardTable       = "award"
	roundTable       = "round"

	columnPlayerID = "player_id"
	columnAwardID  = "award_id"
	columnRoundID  = "round_id"
	columnLevel    = "level"

	columnID   = "id"
	columnType = "type"

	columnEnd = "time_end"
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]award.Record, error) {
	query := sq.
		Select(
			sqlutil.QualifyColumn(awardRecordTable, columnPlayerID),
			sqlutil.QualifyColumn(awardRecordTable, columnAwardID),
			sqlutil.QualifyColumn(awardTable, columnType),
			sqlutil.QualifyColumn(awardRecordTable, columnRoundID),
			sqlutil.QualifyColumn(roundTable, columnEnd),
			sqlutil.QualifyColumn(awardRecordTable, columnLevel),
		).
		From(awardRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			awardTable,
			sqlutil.QualifyColumn(awardRecordTable, columnAwardID),
			sqlutil.QualifyColumn(awardTable, columnID),
		)).
		LeftJoin(fmt.Sprintf(
			"%s ON %s = %s",
			roundTable,
			sqlutil.QualifyColumn(awardRecordTable, columnRoundID),
			sqlutil.QualifyColumn(roundTable, columnID),
		)).
		Where(sq.Eq{sqlutil.QualifyColumn(awardRecordTable, columnPlayerID): playerID}).
		OrderBy(
			fmt.Sprintf("%s ASC", sqlutil.QualifyColumn(awardRecordTable, columnAwardID)),
			fmt.Sprintf("%s ASC", sqlutil.QualifyColumn(awardRecordTable, columnLevel)),
		)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	as := make([]award.Record, 0)
	for rows.Next() {
		var a award.Record
		if err = rows.Scan(
			&a.Player.ID,
			&a.Award.ID,
			&a.Award.Type,
			&a.Round.ID,
			&a.Round.End,
			&a.Level,
		); err != nil {
			return nil, err
		}

		as = append(as, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return as, nil
}
