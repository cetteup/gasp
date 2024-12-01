package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/kit"
)

const (
	kitRecordTable = "player_kit"

	columnPlayerID = "player_id"
	columnKitID    = "kit_id"
	columnTime     = "time"
	columnScore    = "score"
	columnKills    = "kills"
	columnDeaths   = "deaths"
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]kit.Record, error) {
	query := sq.
		Select(
			columnPlayerID,
			columnKitID,
			columnTime,
			columnScore,
			columnKills,
			columnDeaths,
		).
		From(kitRecordTable).
		Where(sq.Eq{columnPlayerID: playerID}).
		OrderBy(fmt.Sprintf("%s ASC", columnKitID))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]kit.Record, 0)
	for rows.Next() {
		var record kit.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Kit.ID,
			&record.Time,
			&record.Score,
			&record.Kills,
			&record.Deaths,
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
