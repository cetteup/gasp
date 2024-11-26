package handler

import (
	"github.com/cetteup/gasp/internal/domain/player"
)

type Handler struct {
	playerRepository player.Repository
}

func NewHandler(playerRepository player.Repository) *Handler {
	return &Handler{
		playerRepository: playerRepository,
	}
}
