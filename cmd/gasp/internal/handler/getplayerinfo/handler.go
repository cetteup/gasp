package getplayerinfo

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getplayerinfo/internal/gather"
	"github.com/cetteup/gasp/cmd/gasp/internal/handler/getplayerinfo/internal/info"
	"github.com/cetteup/gasp/internal/domain/army"
	"github.com/cetteup/gasp/internal/domain/field"
	"github.com/cetteup/gasp/internal/domain/kill"
	"github.com/cetteup/gasp/internal/domain/kit"
	"github.com/cetteup/gasp/internal/domain/player"
	"github.com/cetteup/gasp/internal/domain/vehicle"
	"github.com/cetteup/gasp/internal/domain/weapon"
	"github.com/cetteup/gasp/pkg/asp"
)

type Handler struct {
	gatherer *gather.Gatherer
}

func NewHandler(
	playerRepository player.Repository,
	armyRecordRepository army.RecordRepository,
	fieldRecordRepository field.RecordRepository,
	killHistoryRecordRepository kill.HistoryRecordRepository,
	kitRecordRepository kit.RecordRepository,
	vehicleRecordRepository vehicle.RecordRepository,
	weaponRecordRepository weapon.RecordRepository,
) *Handler {
	return &Handler{
		// Gatherer is "hidden" to only pass repositories to handlers (completely arbitrary design decision)
		gatherer: gather.NewGatherer(
			playerRepository,
			armyRecordRepository,
			fieldRecordRepository,
			killHistoryRecordRepository,
			kitRecordRepository,
			vehicleRecordRepository,
			weaponRecordRepository,
		),
	}
}

func (h *Handler) HandleGET(c echo.Context) error {
	params := struct {
		PID     uint32  `query:"pid" validate:"required"`
		Info    string  `query:"info" validate:"required"`
		Field   *uint16 `query:"map" validate:"omitempty,oneof=0 1 2 3 4 5 6 10 11 12 100 101 102 103 104 105 110 120 200 201 202 300 301 302 303 304 305 306 307 601 602"`
		Kit     *uint8  `query:"kit" validate:"omitempty,oneof=0 1 2 3 4 5 6"`
		Vehicle *uint8  `query:"vehicle" validate:"omitempty,oneof=0 1 2 3 4 5 6"`
		Weapon  *uint8  `query:"weapon" validate:"omitempty,oneof=0 1 2 3 4 5 6 7 8 9 10 11 12 13"`
	}{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	// Info query may contain wildcard "groups" such as "cmb*", which need to be resolved to the underlying keys
	opts := info.NewResolveOptions().
		// Add player id and name as defaults, since those should *always* be the first two keys in the response
		SetDefaultKeys(info.KeyID, info.KeyName).
		// Maybe only replaces the default values if at least one value is non-nil
		MaybeSetFieldIDs(params.Field).
		MaybeSetKitIDs(params.Kit).
		MaybeSetVehicleIDs(params.Vehicle).
		MaybeSetWeaponIDs(params.Weapon)
	keys := info.Resolve(params.Info, opts)

	values, err := h.gatherer.Gather(c.Request().Context(), params.PID, keys)
	if err != nil {
		if errors.Is(err, player.ErrPlayerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to gather values: %w", err))
	}

	resp, err := buildResponse(keys, values)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to build response: %w", err))
	}

	return c.String(http.StatusOK, resp.Serialize())
}

func buildResponse(keys []string, values map[string]string) (*asp.Response, error) {
	resp := asp.NewOKResponse().
		WriteHeader("asof").
		WriteData(asp.Timestamp()).
		// Init empty header and data line to append keys/values to
		WriteHeader().
		WriteData()

	// Order of values is important, as it is expected to match the order of the given keys
	// Hence we iterate over keys, not over values (keyword: random map order)
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			return nil, fmt.Errorf("key is missing from gathered values: %s", key)
		}

		resp.
			AppendHeader(key).
			AppendData(value)

	}

	return resp, nil
}
