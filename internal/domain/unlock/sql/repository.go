package sql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/internal/sqlutil"
	"github.com/cetteup/gasp/internal/util"
)

const (
	unlockTable            = "unlock"
	unlockRecordTable      = "player_unlock"
	unlockRequirementTable = "unlock_requirement"
	playerTable            = "player"

	columnID          = "id"
	columnKitID       = "kit_id"
	columnName        = "name"
	columnDescription = "desc"

	columnPlayerID  = "player_id"
	columnUnlockID  = "unlock_id"
	columnTimestamp = "timestamp"

	columnParentID = "parent_id"
	columnChildID  = "child_id"

	columnRankID = "rank_id"

	virtualColumnUnlocked = "unlocked"
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

type RecordRepository struct {
	runner sq.BaseRunner
}

func NewRecordRepository(runner sq.BaseRunner) *RecordRepository {
	return &RecordRepository{
		runner: runner,
	}
}

func (r *RecordRepository) FindByPlayerID(ctx context.Context, playerID uint32) ([]unlock.Record, error) {
	const (
		playerUnlockCTEName = "pu"
		unlockCTEName       = "u"
	)

	// First, set up a CTE with the player's actual unlocks from the record table and the player's relevant details.
	// Due to the `RIGHT JOIN`, the CTE will contain a single row with the player details even if the player has not
	// yet unlocked any weapons. This row is required, since we need the player details to be able to tell if
	// any player with the given id exists.
	playerUnlockCTE := sq.
		Select(
			columnID,
			columnName,
			columnRankID,
			columnUnlockID,
			columnTimestamp,
		).
		From(unlockRecordTable).
		RightJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerTable,
			sqlutil.QualifyColumn(unlockRecordTable, columnPlayerID),
			sqlutil.QualifyColumn(playerTable, columnID),
		)).
		Where(sq.Eq{columnPlayerID: playerID}).
		// squirrel does not support CTEs, since they are not part of generic SQL,
		// so we need to use prefixes and suffixes to build and combine the expressions.
		Prefix(fmt.Sprintf("WITH %s AS (", playerUnlockCTEName)).
		Suffix(fmt.Sprintf("), %s AS (", unlockCTEName))

	// Next, create another CTE which combines the available unlocks with their requirements. This CTE is
	// primarily used to not have to join these two tables twice, since we need for both "halves" of the UNION later.
	unlockCTE := sq.
		Select(
			columnID,
			columnKitID,
			columnName,
			sqlutil.Quote(columnDescription),
			columnParentID,
		).
		From(sqlutil.Quote(unlockTable)).
		LeftJoin(fmt.Sprintf(
			"%s ON %s = %s",
			unlockRequirementTable,
			sqlutil.QualifyColumn(unlockTable, columnID),
			sqlutil.QualifyColumn(unlockRequirementTable, columnChildID),
		)).
		PrefixExpr(playerUnlockCTE).
		Suffix(")")

	// This will be second part of the UNION, adding all unlocks the player has yet to obtain.
	// Meaning we need to filter out any already obtained unlocks as well as any 2nd tier unlocks for which the 1st
	// tier has not been unlocked yet.
	availableUnlocksUnion := sq.
		Select(
			// We expect to find all 1st-tier unlocks for non-existent players, so make sure we default to safe values.
			fmt.Sprintf("COALESCE(%s, 0) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnID), "player_id"),
			fmt.Sprintf("COALESCE(%s, '') AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnName), "player_name"),
			fmt.Sprintf("COALESCE(%s, 0) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnRankID), "player_rank_id"),
			sqlutil.QualifyColumn(unlockCTEName, columnID),
			sqlutil.QualifyColumn(unlockCTEName, columnKitID),
			sqlutil.QualifyColumn(unlockCTEName, columnName),
			sqlutil.QualifyColumn(unlockCTEName, columnDescription),
			"0", // We're selecting non-obtained unlocks, so hard-set unlocked and timestamp to 0.
			"0",
		).
		From(unlockCTEName).
		LeftJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerUnlockCTEName,
			util.FormatUint(playerID), // "Static" join using the player id since there actually is no common column
			sqlutil.QualifyColumn(playerUnlockCTEName, columnID),
		)).
		Where(sq.And{
			// Exclude any unlocks already obtained by the player
			sq.Expr(fmt.Sprintf(
				// Need to COALESCE here to ensure the NOT IN works. Seems nothing is IN or NOT IN if the list only
				// contains NULL - which is the case for a player without any unlocks. For such a player,
				// pu only contains a single RIGHT JOIN-ed row in which the unlock_id is NULL.
				"%s NOT IN (SELECT COALESCE(%s, 0) FROM %s)",
				sqlutil.QualifyColumn(unlockCTEName, columnID),
				columnUnlockID,
				playerUnlockCTEName,
			)),
			// Only include unlocks that either don't have a parent or those for which the parent was already unlocked.
			sq.Or{
				sq.Expr(fmt.Sprintf("%s IS NULL", sqlutil.QualifyColumn(unlockCTEName, columnParentID))),
				sq.Expr(fmt.Sprintf(
					"%s IN (SELECT %s FROM %s)",
					sqlutil.QualifyColumn(unlockCTEName, columnParentID),
					columnUnlockID,
					playerUnlockCTEName,
				)),
			},
		}).
		OrderBy(columnID).
		Prefix("UNION DISTINCT")

	// Finally, the first part of the union returns all unlocks obtained by the player (if any).
	query := sq.
		Select(
			fmt.Sprintf("COALESCE(%s, 0) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnID), "player_id"),
			fmt.Sprintf("COALESCE(%s, '') AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnName), "player_name"),
			fmt.Sprintf("COALESCE(%s, 0) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnRankID), "player_rank_id"),
			sqlutil.QualifyColumn(unlockCTEName, columnID),
			sqlutil.QualifyColumn(unlockCTEName, columnKitID),
			sqlutil.QualifyColumn(unlockCTEName, columnName),
			sqlutil.QualifyColumn(unlockCTEName, columnDescription),
			fmt.Sprintf("NOT ISNULL(%s) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnUnlockID), virtualColumnUnlocked),
			fmt.Sprintf("COALESCE(%s, 0) AS %s", sqlutil.QualifyColumn(playerUnlockCTEName, columnTimestamp), columnTimestamp),
		).
		From(unlockCTEName).
		InnerJoin(fmt.Sprintf(
			"%s ON %s = %s",
			playerUnlockCTEName,
			sqlutil.QualifyColumn(unlockCTEName, columnID),
			sqlutil.QualifyColumn(playerUnlockCTEName, columnUnlockID),
		)).
		PrefixExpr(unlockCTE).
		SuffixExpr(availableUnlocksUnion)

	rows, err := query.RunWith(r.runner).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	records := make([]unlock.Record, 0)
	for rows.Next() {
		var record unlock.Record
		if err = rows.Scan(
			&record.Player.ID,
			&record.Player.Name,
			&record.Player.RankID,
			&record.Unlock.ID,
			&record.Unlock.KitID,
			&record.Unlock.Name,
			&record.Unlock.Description,
			&record.Unlocked,
			&record.Timestamp,
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
