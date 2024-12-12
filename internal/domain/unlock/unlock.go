package unlock

import (
	"github.com/cetteup/gasp/internal/domain/award"
	"github.com/cetteup/gasp/internal/domain/player"
)

const (
	totalPossibleUnlocks = 7 * 2 // 7 classes, 2 unlocks each
)

type Unlock struct {
	ID          uint16
	Name        string
	Description string
	Kit         KitRef
}

type KitRef struct {
	ID uint8
}

type Record struct {
	Player    PlayerRef
	Unlock    Unlock
	Unlocked  bool
	Timestamp uint32 // zero if unlocked is false
}

type PlayerRef struct {
	ID uint32
}

func DetermineAvailablePoints(p player.Player, unlockRecords []Record, awardRecords []award.Record) int {
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
	rankPoints := max(min(int(p.Rank.ID)-1, 7), 0)

	// One point per level two badge
	badgePoints := 0
	for _, record := range awardRecords {
		if record.Award.Type == award.TypeBadge && award.IsKitBadge(record.Award.ID) && record.Level == 2 {
			badgePoints++
		}
	}
	// Unless the data in the db is inconsistent, more than 7 points should never be seen
	badgePoints = min(badgePoints, 7)

	return max(rankPoints+badgePoints-usedPoints, 0)
}
