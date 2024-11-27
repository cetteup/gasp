package getawardsinfo

import (
	"github.com/cetteup/gasp/internal/domain/award"
)

type RecordDTO struct {
	Award uint32
	Level uint64
	When  uint32
	First uint32
}

func EncodeRecords(records []award.Record) []RecordDTO {
	dtos := make([]RecordDTO, 0, len(records))
	medals := make(map[uint32]int)

	for _, record := range records {
		switch record.Award.Type {
		case award.TypeRibbon, award.TypeBadge:
			dtos = append(dtos, EncodeRecord(record))
		case award.TypeMedal:
			if i, exists := medals[record.Award.ID]; !exists {
				// Add medal entry if not yet present
				dtos = append(dtos, EncodeRecord(record))
				medals[record.Award.ID] = len(dtos) - 1
			} else {
				// Update medal entry if previously added
				dtos[i].Level += record.Level
				dtos[i].When = max(dtos[i].When, record.Round.End)
				dtos[i].First = min(dtos[i].First, record.Round.End)
			}
		}
	}

	return dtos
}

func EncodeRecord(record award.Record) RecordDTO {
	dto := RecordDTO{
		Award: record.Award.ID,
		Level: record.Level,
		When:  record.Round.End,
	}

	// Only medals have first set since they can be awarded multiple times
	if record.Award.Type == award.TypeMedal {
		dto.First = record.Round.End
	}

	return dto
}
