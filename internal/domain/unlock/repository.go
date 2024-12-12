package unlock

import (
	"context"
	"errors"
)

var (
	ErrRecordNotUnlocked = errors.New("record not unlocked")
)

type Repository interface {
	FindAll(ctx context.Context) ([]Unlock, error)
}

type RecordRepository interface {
	Insert(ctx context.Context, record Record) error
	FindByPlayerID(ctx context.Context, playerID uint32) ([]Record, error)
}
