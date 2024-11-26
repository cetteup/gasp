package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/pkg/asp"
)

const (
	dummyPID = 0

	prefixInvalid = "INVALID"
	prefixBanned  = "BANNED"
)

func (h *Handler) HandleGetVerifyPlayer(c echo.Context) error {
	params := struct {
		PID  uint32 `query:"pid" validate:"required"`
		Nick string `query:"SoldierNick" validate:"required"`
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
			// Use prefix nick and dummy PID to ensure the server will issue a kick/ban (either comparisons will fail)
			return c.String(http.StatusOK, buildResponse(
				addPrefix(params.Nick, prefixInvalid),
				params.Nick,
				dummyPID,
				params.PID,
			).Serialize())
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find player: %w", err))
	}

	if p.PermanentlyBanned {
		// Prefixed (this modified) name will trigger kick/ban on the server
		return c.String(http.StatusOK, buildResponse(
			addPrefix(p.Name, prefixBanned),
			params.Nick,
			p.ID,
			params.PID,
		).Serialize())
	}

	return c.String(http.StatusOK, buildResponse(
		p.Name,
		params.Nick,
		p.ID,
		params.PID,
	).Serialize())
}

// buildResponse Signature analog to default onPlayerNameValidated handler
func buildResponse(realNick, oldNick string, realPID, oldPID uint32) *asp.Response {
	resp := asp.NewOKResponse().
		WriteHeader("pid", "nick", "spid", "asof")

	if realNick == oldNick && realPID == oldPID {
		resp.
			// Using oldNick instead of realNick here to ensure we return the (determined matching) name as-is
			// The Python onPlayerNameValidated only receives and compares the old/real values (case-sensitive!)
			// Returning realNick would cause players with mismatched case to be banned, even if their login backend allows it
			WriteData(formatUint(realPID), oldNick, formatUint(oldPID), asp.Timestamp()).
			WriteHeader("result").
			WriteData("Ok")
	} else if realNick != oldNick && realPID != oldPID {
		resp.
			WriteData(formatUint(realPID), realNick, formatUint(oldPID), asp.Timestamp()).
			WriteHeader("result").
			// We obviously cannot validate the auth param, but neither value matching would indicate
			// that the player was not found and this is the closest to "completely invalid" there is
			// (no player can be logged into a profile that does not exist)
			WriteData("InvalidAuthProfileID")
	} else if realNick != oldNick {
		resp.
			WriteData(formatUint(realPID), realNick, formatUint(oldPID), asp.Timestamp()).
			WriteHeader("result").
			WriteData("InvalidReportedNick")
	} else {
		// Currently unused as realNick differs from oldNick for any non-ok response
		// Primarily here for completeness-sake
		resp.
			WriteData(formatUint(realPID), realNick, formatUint(oldPID), asp.Timestamp()).
			WriteHeader("result").
			WriteData("InvalidReportedProfileID")
	}

	return resp
}

func formatUint[T uint | uint8 | uint16 | uint32 | uint64](i T) string {
	// Converting to uint64 is always safe as i is *at most* 64-bit
	return strconv.FormatUint(uint64(i), 10)
}

func addPrefix(nick string, prefix string) string {
	// `[prefix] nick` usually get cut off after 23 characters in the game's client-server protocols
	// While the limit appears to not be applied to values returned by the validation,
	// it's probably best to follow that convention/limit
	prefixed := prefix + " " + nick
	return prefixed[:min(len(prefixed), 23)]
}
