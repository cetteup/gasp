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
			sqlutil.Qualify(awardRecordTable, columnPlayerID),
			sqlutil.Qualify(awardRecordTable, columnAwardID),
			sqlutil.Qualify(awardTable, columnType),
			sqlutil.Qualify(awardRecordTable, columnRoundID),
			sqlutil.Qualify(roundTable, columnEnd),
			sqlutil.Qualify(awardRecordTable, columnLevel),
		).
		From(awardRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			awardTable,
			sqlutil.Qualify(awardRecordTable, columnAwardID),
			sqlutil.Qualify(awardTable, columnID),
		)).
		LeftJoin(fmt.Sprintf(
			"%s ON %s = %s",
			roundTable,
			sqlutil.Qualify(awardRecordTable, columnRoundID),
			sqlutil.Qualify(roundTable, columnID),
		)).
		Where(sq.Eq{sqlutil.Qualify(awardRecordTable, columnPlayerID): playerID}).
		OrderBy(
			fmt.Sprintf("%s ASC", sqlutil.Qualify(awardRecordTable, columnAwardID)),
			fmt.Sprintf("%s ASC", sqlutil.Qualify(awardRecordTable, columnLevel)),
		)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]award.Record, 0)
	for rows.Next() {
		var record award.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Award.ID,
			&record.Award.Type,
			&record.Round.ID,
			&record.Round.End,
			&record.Level,
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
