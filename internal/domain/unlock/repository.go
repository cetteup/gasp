package unlock

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Unlock, error)
}
