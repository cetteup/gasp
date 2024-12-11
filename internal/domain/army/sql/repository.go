package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/army"
)

const (
	armyRecordTable = "player_army"

	columnPlayerID        = "player_id"
	columnArmyID          = "army_id"
	columnTime            = "time"
	columnWins            = "wins"
	columnLosses          = "losses"
	columnScore           = "score"
	columnBestRoundScore  = "best"
	columnWorstRoundScore = "worst"
	columnBestRounds      = "brnd"
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]army.Record, error) {
	query := sq.
		Select(
			columnPlayerID,
			columnArmyID,
			columnTime,
			columnWins,
			columnLosses,
			columnScore,
			columnBestRoundScore,
			columnWorstRoundScore,
			columnBestRounds,
		).
		From(armyRecordTable).
		Where(sq.Eq{columnPlayerID: playerID}).
		OrderBy(fmt.Sprintf("%s ASC", columnArmyID))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]army.Record, 0)
	for rows.Next() {
		var record army.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Army.ID,
			&record.Time,
			&record.Wins,
			&record.Losses,
			&record.Score,
			&record.BestRoundScore,
			&record.WorstRoundScore,
			&record.BestRounds,
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
