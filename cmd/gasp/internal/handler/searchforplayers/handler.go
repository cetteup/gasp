package searchforplayers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/asp"
)

const (
	whereContains   = "a" // "any"
	whereBeginsWith = "b"
	whereEndsWith   = "e"
	whereEquals     = "x" // "exactly"

	sortASC  = "a"
	sortDESC = "r" // "reverse"
)

type Handler struct {
	playerRepository player.Repository
}

func NewHandler(playerRepository player.Repository) *Handler {
	return &Handler{
		playerRepository: playerRepository,
	}
}

func (h *Handler) HandleGET(c echo.Context) error {
	params := struct {
		Nick  string `query:"nick" validate:"required"`
		Where string `query:"where" validate:"omitempty,oneof=a b e x"`
		Sort  string `query:"sort" validate:"omitempty,oneof=a r"`
	}{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	players, err := h.playerRepository.FindWithNameMatching(
		c.Request().Context(),
		params.Nick,
		toMatchCondition(params.Where),
		toSortOrder(params.Sort),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find players: %w", err))
	}

	resp := asp.NewOKResponse().
		WriteHeader("asof").
		WriteData(asp.Timestamp()).
		WriteHeader("n", "pid", "nick", "score")

	for i, p := range players {
		resp.WriteData(
			util.FormatInt(i+1),
			util.FormatUint(p.ID),
			p.Name,
			util.FormatInt(p.Score),
		)
	}

	return c.String(http.StatusOK, resp.Serialize())
}

// toMatchCondition Returns a default rather than an error for unmapped values
func toMatchCondition(where string) player.MatchCondition {
	switch where {
	case whereBeginsWith:
		return player.MatchConditionBeginsWith
	case whereEndsWith:
		return player.MatchConditionEndsWith
	case whereEquals:
		return player.MatchConditionEquals
	case whereContains:
		fallthrough
	default:
		return player.MatchConditionContains
	}
}

// toSortOrder Returns a default rather than an error for unmapped values
func toSortOrder(sort string) player.SortOrder {
	switch sort {
	case sortDESC:
		return player.SortOrderDESC
	case sortASC:
		fallthrough
	default:
		return player.SortOrderASC
	}
}
