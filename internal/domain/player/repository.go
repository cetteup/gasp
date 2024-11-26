package player

import (
	"context"
	"errors"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type Repository interface {
	FindByID(ctx context.Context, id uint32) (Player, error)
}
