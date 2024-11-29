package getunlocksinfo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/asp"
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

	unlockRecords, err := h.unlockRecordRepository.FindByPlayerID(c.Request().Context(), params.PID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find unlock records: %w", err))
	}

	// To save an extra database query, the player is only loaded with the unlock records. If the player exists,
	// every record contains the player details. If the player was not found, all player details will be zero values.
	usedPoints := 0
	var p unlock.PlayerStub
	for _, record := range unlockRecords {
		if record.Unlocked {
			p = record.Player
			usedPoints++
		}
		// Set player if not set yet and not a zero value
		if p.ID == 0 && record.Player.ID != 0 {
			p = record.Player
		}
	}

	// If player id is still zero, no record container player details - indicating that the player was not found.
	if p.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	availablePoints, err := h.determineAvailableUnlockPoints(c.Request().Context(), p, usedPoints)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to determine availabe unlock points: %w", err))
	}

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

func (h *Handler) determineAvailableUnlockPoints(ctx context.Context, p unlock.PlayerStub, usedPoints int) (int, error) {
	// Player cannot have any unlock points available if they already unlocked everything
	if usedPoints >= totalPossibleUnlocks {
		return 0, nil
	}

	// No more than 7 unlocks via rank, but don't let the number go negative
	rankPoints := max(min(int(p.RankID)-1, 7), 0)

	// Filtering by award type/level here is slower, not faster
	awardRecords, err := h.awardRecordRepository.FindByPlayerID(ctx, p.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to find award records: %w", err)
	}

	// One point per level two badge
	badgePoints := 0
	for _, record := range awardRecords {
		if record.Award.Type == award.TypeBadge && isKitBadge(record.Award.ID) && record.Level == 2 {
			badgePoints++
		}
	}
	// Unless the data in the db is inconsistent, more than 7 points should never be seen
	badgePoints = min(badgePoints, 7)

	return max(rankPoints+badgePoints-usedPoints, 0), nil
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
