package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/field"
)

const (
	fieldRecordTable = "player_map"

	columnPlayerID = "player_id"
	columnFieldID  = "map_id"
	columnTime     = "time"
	columnWins     = "wins"
	columnLosses   = "losses"

	highestOfficialFieldID = 603 // Operation Blue Pearl
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]field.Record, error) {
	query := sq.
		Select(
			columnPlayerID,
			columnFieldID,
			columnTime,
			columnWins,
			columnLosses,
		).
		From(fieldRecordTable).
		Where(sq.And{
			sq.Eq{columnPlayerID: playerID},
			// Only fetch records for official maps
			sq.LtOrEq{columnFieldID: highestOfficialFieldID},
		}).
		OrderBy(fmt.Sprintf("%s ASC", columnFieldID))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]field.Record, 0)
	for rows.Next() {
		var record field.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Field.ID,
			&record.Time,
			&record.Wins,
			&record.Losses,
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
