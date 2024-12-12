package ranknotification

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/pkg/asp"
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
		PID uint32 `query:"pid" validate:"required"`
	}{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	p, err := h.playerRepository.FindByID(c.Request().Context(), params.PID)
	if err != nil {
		if errors.Is(err, player.ErrPlayerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find player: %w", err))
	}

	// Only reset flags if required
	if p.RankChanged || p.RankDecreased {
		err2 := h.playerRepository.ResetRankChangeFlags(c.Request().Context(), p.ID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to reset rank change flags: %w", err2))
		}
	}

	return c.String(http.StatusOK, asp.NewOKResponse().Serialize())
}
