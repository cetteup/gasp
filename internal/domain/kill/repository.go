package kill

import (
	"context"
)

type HistoryRecordRepository interface {
	FindTopRelatedByPlayerID(ctx context.Context, playerID uint32) ([]HistoryRecord, error)
}
