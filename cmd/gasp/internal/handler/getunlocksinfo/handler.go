package getunlocksinfo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/internal/util"
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

	availablePoints := unlock.DetermineAvailablePoints(p, unlockRecords, awardRecords)
	resp := asp.NewOKResponse().
		WriteHeader("pid", "nick", "asof").
		WriteData(util.FormatUint(p.ID), p.Name, asp.Timestamp()).
		WriteHeader("enlisted", "officer").
		WriteData(util.FormatInt(availablePoints), "0").
		WriteHeader("id", "state")

	for _, record := range unlockRecords {
		resp.WriteData(util.FormatUint(record.Unlock.ID), encodeUnlocked(record.Unlocked))
	}

	return c.String(http.StatusOK, resp.Serialize())
}

func encodeUnlocked(unlocked bool) string {
	if unlocked {
		return "s"
	}
	return "n"
}
