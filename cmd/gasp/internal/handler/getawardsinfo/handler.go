package getawardsinfo

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/asp"
)

type Handler struct {
	awardRecordRepository award.RecordRepository
}

func NewHandler(awardRecordRepository award.RecordRepository) *Handler {
	return &Handler{
		awardRecordRepository: awardRecordRepository,
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

	records, err := h.awardRecordRepository.FindByPlayerID(c.Request().Context(), params.PID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find award records: %w", err))
	}

	resp := asp.NewOKResponse().
		WriteHeader("pid", "asof").
		WriteData(util.FormatUint(params.PID), asp.Timestamp()).
		WriteHeader("award", "level", "when", "first")

	dtos := EncodeRecords(records)
	for _, dto := range dtos {
		resp.WriteData(
			util.FormatUint(dto.Award),
			util.FormatUint(dto.Level),
			util.FormatUint(dto.When),
			util.FormatUint(dto.First),
		)
	}

	return c.String(http.StatusOK, resp.Serialize())
}
