package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/weapon"
	"github.com/cetteup/gasp/internal/sqlutil"
)

const (
	weaponRecordTable = "player_weapon"
	weaponTable       = "weapon"

	columnPlayerID      = "player_id"
	columnWeaponID      = "weapon_id"
	columnTime          = "time"
	columnScore         = "score"
	columnKills         = "kills"
	columnDeaths        = "deaths"
	columnShotsFired    = "fired"
	columnShotsHit      = "hits"
	columnTimesDeployed = "deployed"

	columnID          = "id"
	columnName        = "name"
	columnIsExplosive = "is_explosive"
	columnIsEquipment = "is_equipment"
)

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]weapon.Record, error) {
	query := sq.
		Select(
			sqlutil.QualifyColumn(weaponRecordTable, columnPlayerID),
			sqlutil.QualifyColumn(weaponRecordTable, columnWeaponID),
			sqlutil.QualifyColumn(weaponTable, columnName),
			sqlutil.QualifyColumn(weaponTable, columnIsExplosive),
			sqlutil.QualifyColumn(weaponTable, columnIsEquipment),
			sqlutil.QualifyColumn(weaponRecordTable, columnTime),
			sqlutil.QualifyColumn(weaponRecordTable, columnScore),
			sqlutil.QualifyColumn(weaponRecordTable, columnKills),
			sqlutil.QualifyColumn(weaponRecordTable, columnDeaths),
			sqlutil.QualifyColumn(weaponRecordTable, columnShotsFired),
			sqlutil.QualifyColumn(weaponRecordTable, columnShotsHit),
			sqlutil.QualifyColumn(weaponRecordTable, columnTimesDeployed),
		).
		From(weaponRecordTable).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			weaponTable,
			sqlutil.QualifyColumn(weaponRecordTable, columnWeaponID),
			sqlutil.QualifyColumn(weaponTable, columnID),
		)).
		Where(sq.Eq{sqlutil.QualifyColumn(weaponRecordTable, columnPlayerID): playerID}).
		OrderBy(fmt.Sprintf("%s ASC", sqlutil.QualifyColumn(weaponRecordTable, columnWeaponID)))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]weapon.Record, 0)
	for rows.Next() {
		var record weapon.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Weapon.ID,
			&record.Weapon.Name,
			&record.Weapon.IsExplosive,
			&record.Weapon.IsEquipment,
			&record.Time,
			&record.Score,
			&record.Kills,
			&record.Deaths,
			&record.ShotsFired,
			&record.ShotsHit,
			&record.TimesDeployed,
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
