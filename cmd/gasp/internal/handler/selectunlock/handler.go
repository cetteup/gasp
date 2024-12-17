package selectunlock

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/pkg/asp"
	"github.com/cetteup/gasp/pkg/task"
)

type Handler struct {
	playerRepository       player.Repository
	awardRecordRepository  award.RecordRepository
	unlockRecordRepository unlock.RecordRepository
}

func NewHandler(
	playerRepository player.Repository,
	awardRecordRepository award.RecordRepository,
	unlockRecordRepository unlock.RecordRepository,
) *Handler {
	return &Handler{
		playerRepository:       playerRepository,
		awardRecordRepository:  awardRecordRepository,
		unlockRecordRepository: unlockRecordRepository,
	}
}

func (h *Handler) HandlePOST(c echo.Context) error {
	params := struct {
		PID      uint32 `form:"pid" validate:"required"`
		UnlockID uint16 `form:"id" validate:"required,oneof=11 22 33 44 55 66 77 88 99 111 222 333 444 555"`
	}{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	var p player.Player
	var unlockRecords []unlock.Record
	var awardRecords []award.Record
	var runner task.AsyncRunner
	runner.Append(func(ctx context.Context) error {
		var err2 error
		p, err2 = h.playerRepository.FindByID(ctx, params.PID)
		if err2 != nil {
			return fmt.Errorf("failed to find player: %w", err2)
		}
		return nil
	})
	runner.Append(func(ctx context.Context) error {
		var err2 error
		unlockRecords, err2 = h.unlockRecordRepository.FindByPlayerID(ctx, params.PID)
		if err2 != nil {
			return fmt.Errorf("failed to find unlock records: %w", err2)
		}
		return nil
	})
	runner.Append(func(ctx context.Context) error {
		var err2 error
		awardRecords, err2 = h.awardRecordRepository.FindByPlayerID(ctx, params.PID)
		if err2 != nil {
			return fmt.Errorf("failed to find award records: %w", err2)
		}
		return nil
	})

	if err := runner.Run(c.Request().Context()); err != nil {
		if errors.Is(err, player.ErrPlayerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return err
	}

	// Ensure selected unlock has not yet been unlocked
	// Could also consider this a noop, but the endpoint is POST not PUT, not being idempotent is fine
	if slices.ContainsFunc(unlockRecords, func(record unlock.Record) bool {
		// Records may contain non-unlocked entries
		return record.Unlock.ID == params.UnlockID && record.Unlocked
	}) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity).SetInternal(errors.New("unlock already unlocked"))
	}

	// Ensure players has points available
	if unlock.DetermineAvailablePoints(p, unlockRecords, awardRecords) < 1 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity).SetInternal(errors.New("no unlock points available"))
	}

	record := unlock.Record{
		Player: unlock.PlayerRef{
			ID: p.ID,
		},
		Unlock: unlock.Unlock{
			ID: params.UnlockID,
		},
		Unlocked: true,
		// Will overflow on 7 February 2106 at 06:28:15 UTC
		Timestamp: uint32(time.Now().UTC().Unix()),
	}

	err := h.unlockRecordRepository.Insert(c.Request().Context(), record)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to insert unlock record: %w", err))
	}

	return c.String(http.StatusOK, asp.NewOKResponse().Serialize())
}
