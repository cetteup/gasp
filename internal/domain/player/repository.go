package player

import (
	"context"
	"errors"
)

type MatchCondition int
type SortOrder int

const (
	MatchConditionContains MatchCondition = iota
	MatchConditionBeginsWith
	MatchConditionEndsWith
	MatchConditionEquals

	SortOrderASC  SortOrder = 1
	SortOrderDESC SortOrder = -1
)

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type Repository interface {
	FindByID(ctx context.Context, id uint32) (Player, error)
	FindWithNameMatching(ctx context.Context, name string, condition MatchCondition, order SortOrder) ([]Player, error)
}
