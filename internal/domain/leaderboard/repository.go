package leaderboard

import (
	"context"
)

type ScoreType int

const (
	ScoreTypeOverall ScoreType = iota
	ScoreTypeCommand
	ScoreTypeTeam
	ScoreTypeCombat
)

type Repository interface {
	FindTopPlayersByScore(ctx context.Context, scoreType ScoreType, filter Filter) ([]Entry[PlayerStub], int, error)
	FindTopPlayersByKit(ctx context.Context, kitID uint8, filter Filter) ([]Entry[KitRecord], int, error)
	FindTopPlayersByVehicle(ctx context.Context, vehicleID uint8, filter Filter) ([]Entry[VehicleRecord], int, error)
	FindTopPlayersByWeapon(ctx context.Context, weaponID uint8, filter Filter) ([]Entry[WeaponRecord], int, error)
	FindRisingStars(ctx context.Context, filter Filter) ([]Entry[RisingStar], int, error)
	GetRisingStarUpdateTimestamp(ctx context.Context) (uint32, error)
}
