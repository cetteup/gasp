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

const (
	totalPossibleUnlocks = 7 * 2 // 7 classes, 2 unlocks each
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
			if errors.Is(err2, player.ErrPlayerNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}
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
		// Return error as is so that any HTTPError returned by a task can be unwrapped and returned to the client.
		// Note: Only a single task may return an HTTPError, else we end up with a race condition/flakiness
		// (first task to return an HTTPError would set the response code).
		return err
	}

	availablePoints := determineAvailableUnlockPoints(p, unlockRecords, awardRecords)
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

func determineAvailableUnlockPoints(p player.Player, unlockRecords []unlock.Record, awardRecords []award.Record) int {
	usedPoints := 0
	for _, record := range unlockRecords {
		if record.Unlocked {
			usedPoints++
		}
	}

	// Player cannot have any unlock points available if they already unlocked everything
	if usedPoints >= totalPossibleUnlocks {
		return 0
	}

	// No more than 7 unlocks via rank, but don't let the number go negative
	rankPoints := max(min(int(p.RankID)-1, 7), 0)

	// One point per level two badge
	badgePoints := 0
	for _, record := range awardRecords {
		if record.Award.Type == award.TypeBadge && isKitBadge(record.Award.ID) && record.Level == 2 {
			badgePoints++
		}
	}
	// Unless the data in the db is inconsistent, more than 7 points should never be seen
	badgePoints = min(badgePoints, 7)

	return max(rankPoints+badgePoints-usedPoints, 0)
}

func isKitBadge(awardID uint32) bool {
	switch awardID {
	case
		1031119, // Assault
		1031120, // Anti-tank
		1031109, // Sniper
		1031115, // Spec-Ops
		1031121, // Support
		1031105, // Engineer
		1031113: // Medic
		return true
	default:
		return false
	}
}

func encodeUnlocked(unlocked bool) string {
	if unlocked {
		return "s"
	}
	return "n"
}
