package getbackendinfo

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cetteup/gasp/internal/domain/unlock"
	"github.com/cetteup/gasp/internal/util"
	"github.com/cetteup/gasp/pkg/asp"
)

type Handler struct {
	unlockRepository unlock.Repository
}

func NewHandler(unlockRepository unlock.Repository) *Handler {
	return &Handler{
		unlockRepository: unlockRepository,
	}
}

func (h *Handler) HandleGET(c echo.Context) error {
	unlocks, err := h.unlockRepository.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("failed to find unlocks: %w", err))
	}

	resp := asp.NewOKResponse().
		WriteHeader("ver", "now").
		WriteData("0.1", asp.Timestamp()).
		WriteHeader("id", "kit", "name", "descr")

	for _, u := range unlocks {
		resp.WriteData(
			util.FormatUint(u.ID),
			util.FormatUint(u.KitID),
			u.Name,
			u.Description,
		)
	}

	return c.String(http.StatusOK, resp.Serialize())
}
