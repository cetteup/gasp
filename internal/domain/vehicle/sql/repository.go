package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/vehicle"
)

const (
	vehicleRecordTable = "player_vehicle"

	columnPlayerID  = "player_id"
	columnVehicleID = "vehicle_id"
	columnTime      = "time"
	columnScore     = "score"
	columnKills     = "kills"
	columnDeaths    = "deaths"
	columnRoadKills = "roadkills"
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]vehicle.Record, error) {
	query := sq.
		Select(
			columnPlayerID,
			columnVehicleID,
			columnTime,
			columnScore,
			columnKills,
			columnDeaths,
			columnRoadKills,
		).
		From(vehicleRecordTable).
		Where(sq.Eq{columnPlayerID: playerID}).
		OrderBy(fmt.Sprintf("%s ASC", columnVehicleID))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]vehicle.Record, 0)
	for rows.Next() {
		var record vehicle.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Vehicle.ID,
			&record.Time,
			&record.Score,
			&record.Kills,
			&record.Deaths,
			&record.RoadKills,
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
