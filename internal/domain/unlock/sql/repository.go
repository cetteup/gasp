package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/internal/sqlutil"
)

const (
	unlockTable = "unlock"

	columnID          = "id"
	columnKitID       = "kit_id"
	columnName        = "name"
	columnDescription = "desc"
)

type Repository struct {
	runner sq.BaseRunner
}

func NewRepository(runner sq.BaseRunner) *Repository {
	return &Repository{
		runner: runner,
	}
}

func (r *Repository) FindAll(ctx context.Context) ([]unlock.Unlock, error) {
	query := sq.
		Select(
			columnID,
			columnKitID,
			columnName,
			sqlutil.Quote(columnDescription), // DESC is a reserved keyword
		).
		From(sqlutil.Quote(unlockTable)). // UNLOCK is a reserved keyword
		OrderBy(fmt.Sprintf("%s ASC", columnID))

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	unlocks := make([]unlock.Unlock, 0)
	for rows.Next() {
		var u unlock.Unlock
		if err = rows.Scan(
			&u.ID,
			&u.KitID,
			&u.Name,
			&u.Description,
		); err != nil {
			return nil, err
		}

		unlocks = append(unlocks, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return unlocks, nil
}
