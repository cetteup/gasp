package getleaderboard

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getleaderboard/internal/gather"
	"github.com/cetteup/gasp/internal/domain/leaderboard"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/asp"
)

type Gatherer interface {
	Gather(ctx context.Context, t, id string, position, before, after uint32, pid *uint32) (gather.GatheredData, error)
}

type Handler struct {
	gatherer Gatherer
}

func NewHandler(leaderboardRepository leaderboard.Repository) *Handler {
	return &Handler{
		// Gatherer is "hidden" to only pass repositories to handlers (completely arbitrary design decision)
		gatherer: gather.NewGatherer(leaderboardRepository),
	}
}

func (h *Handler) HandleGET(c echo.Context) error {
	params := struct {
		Type     string  `query:"type" validate:"required,oneof=score kit vehicle weapon risingstar"`
		ID       string  `query:"id" validate:"required_unless=Type risingstar,omitempty,oneof=overall combat commander team 0 1 2 3 4 5 6 7 8"`
		Position uint32  `query:"pos"`
		Before   uint32  `query:"before"`
		After    uint32  `query:"after"`
		PID      *uint32 `query:"pid"`
	}{
		// Default values
		Position: 1,
		Before:   0,
		After:    19,
	}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	data, err := h.gatherer.Gather(
		c.Request().Context(),
		params.Type,
		params.ID,
		params.Position,
		params.Before,
		params.After,
		params.PID,
	)
	if err != nil {
		if errors.Is(err, gather.ErrInvalidLeaderboardType) || errors.Is(err, gather.ErrInvalidLeaderboardID) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to gather data: %w", err))
	}

	resp, err := buildResponse(data.Keys, data.Entries, data.Size, data.AsOf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to build response: %w", err))
	}

	return c.String(http.StatusOK, resp.Serialize())
}

func buildResponse(keys []string, entries []map[string]string, size int, asOf uint32) (*asp.Response, error) {
	resp := asp.NewOKResponse().
		WriteHeader("size", "asof").
		WriteData(util.FormatInt(size), util.FormatUint(asOf)).
		WriteHeader(keys...)

	for _, entry := range entries {
		// Init empty data line to append values to
		resp.WriteData()
		for _, key := range keys {
			value, ok := entry[key]
			if !ok {
				return nil, fmt.Errorf("key is missing from gathered values: %s", key)
			}

			resp.AppendData(value)
		}
	}

	return resp, nil
}
