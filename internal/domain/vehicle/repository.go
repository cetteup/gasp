package vehicle

import (
	"context"
)

type RecordRepository interface {
	FindByPlayerID(ctx context.Context, playerID uint32) ([]Record, error)
}
