package getplayerinfo

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

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
	gatherer *Gatherer
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
		gatherer: NewGatherer(
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
		PID  uint32 `query:"pid" validate:"required"`
		Info string `query:"info" validate:"required"`
	}{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("failed to bind request parameters: %w", err))
	}

	if err := validator.New().StructCtx(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(fmt.Errorf("invalid parameters: %w", err))
	}

	// Info query may contain wildcard "groups" such as "cmb*", which need to be resolved to the underlying keys
	// Add player id and name as defaults, since those should *always* be the first two keys in the response
	keys := resolveInfoKeys(params.Info, keyID, keyName)

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

	// var p player.Player
	// var armyRecords []army.Record
	// var fieldRecords []field.Record
	// var kitRecords []kit.Record
	// var vehicleRecords []vehicle.Record
	// var weaponRecords []weapon.Record
	// var runner task.AsyncRunner
	// pp := task.WithRunner(&runner, func() (*player.Player, task.Task) {
	//	var t *player.Player
	//	return t, func(ctx context.Context) error {
	//		v, err2 := h.playerRepository.FindByID(ctx, params.PID)
	//		if err2 != nil {
	//			if errors.Is(err2, player.ErrPlayerNotFound) {
	//				return echo.NewHTTPError(http.StatusNotFound)
	//			}
	//			return fmt.Errorf("failed to find player: %w", err2)
	//		}
	//		t = &v
	//		return nil
	//	}
	// })
	//
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	p, err2 = h.playerRepository.FindByID(ctx, params.PID)
	//	if err2 != nil {
	//		if errors.Is(err2, player.ErrPlayerNotFound) {
	//			return echo.NewHTTPError(http.StatusNotFound)
	//		}
	//		return fmt.Errorf("failed to find player: %w", err2)
	//	}
	//	return nil
	// })
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	armyRecords, err2 = h.armyRecordRepository.FindByPlayerID(ctx, params.PID)
	//	if err2 != nil {
	//		return fmt.Errorf("failed to find army records: %w", err2)
	//	}
	//	return nil
	// })
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	fieldRecords, err2 = h.fieldRecordRepository.FindByPlayerID(ctx, params.PID)
	//	if err2 != nil {
	//		return fmt.Errorf("failed to find field records: %w", err2)
	//	}
	//	return nil
	// })
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	kitRecords, err2 = h.kitRecordRepository.FindByPlayerID(ctx, params.PID)
	//	if err2 != nil {
	//		return fmt.Errorf("failed to find kit records: %w", err2)
	//	}
	//	return nil
	// })
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	vehicleRecords, err2 = h.vehicleRecordRepository.FindByPlayerID(ctx, params.PID)
	//	if err2 != nil {
	//		return fmt.Errorf("failed to find vehicle records: %w", err2)
	//	}
	//	return nil
	// })
	// runner.Append(func(ctx context.Context) error {
	//	var err2 error
	//	weaponRecords, err2 = h.weaponRecordRepository.FindByPlayerID(ctx, params.PID)
	//	if err2 != nil {
	//		return fmt.Errorf("failed to find weapon records: %w", err2)
	//	}
	//	return nil
	// })
	//
	// if err = runner.Run(c.Request().Context()); err != nil {
	//	// Return error as is so that any HTTPError returned by a task can be unwrapped and returned to the client.
	//	// Note: Only a single task may return an HTTPError, else we end up with a race condition/flakiness
	//	// (first task to return an HTTPError would set the response code).
	//	return err
	// }
	// println(pp.Name)
	//
	// resp := asp.NewOKResponse().
	//	WriteHeader("asof").
	//	WriteData(asp.Timestamp()).
	//	WriteHeader("pid", "nick").
	//	WriteData(util.FormatUint(p.ID), p.Name)
	//
	// resp.WriteHeader("army", "time", "wins", "losses", "best-rounds")
	// for _, record := range armyRecords {
	//	resp.WriteData(
	//		util.FormatUint(record.Army.ID),
	//		util.FormatUint(record.Time),
	//		util.FormatUint(record.Wins),
	//		util.FormatUint(record.Losses),
	//		util.FormatInt(record.BestRounds),
	//	)
	// }
	//
	// resp.WriteHeader("map", "time", "wins", "losses")
	// for _, record := range fieldRecords {
	//	resp.WriteData(
	//		util.FormatUint(record.Field.ID),
	//		util.FormatUint(record.Time),
	//		util.FormatUint(record.Wins),
	//		util.FormatUint(record.Losses),
	//	)
	// }
	//
	// resp.WriteHeader("kit", "time", "kills", "deaths")
	// for _, record := range kitRecords {
	//	resp.WriteData(
	//		util.FormatUint(record.Kit.ID),
	//		util.FormatUint(record.Time),
	//		util.FormatUint(record.Kills),
	//		util.FormatUint(record.Deaths),
	//	)
	// }
	//
	// resp.WriteHeader("vehicle", "time", "kills", "deaths", "road-kills")
	// for _, record := range vehicleRecords {
	//	resp.WriteData(
	//		util.FormatUint(record.Vehicle.ID),
	//		util.FormatUint(record.Time),
	//		util.FormatUint(record.Kills),
	//		util.FormatUint(record.Deaths),
	//		util.FormatUint(record.RoadKills),
	//	)
	// }
	//
	// resp.WriteHeader("weapon", "time", "kills", "deaths", "shots-fired")
	// for _, record := range weaponRecords {
	//	resp.WriteData(
	//		util.FormatUint(record.Weapon.ID),
	//		util.FormatUint(record.Time),
	//		util.FormatUint(record.Kills),
	//		util.FormatUint(record.Deaths),
	//		util.FormatUint(record.ShotsFired),
	//		util.FormatUint(record.ShotsHit),
	//	)
	// }
	//
	// return c.String(http.StatusOK, resp.Serialize())
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
