package unlock

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Unlock, error)
}

type RecordRepository interface {
	FindByPlayerID(ctx context.Context, playerID uint32) ([]Record, error)
}
