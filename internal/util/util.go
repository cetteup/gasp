package util

import (
	"fmt"
	"strconv"

	"github.com/cetteup/gasp/internal/constraints"
)

func FormatUint[T constraints.UnsignedInteger](i T) string {
	// Converting to uint64 is always safe as i is *at most* 64-bit
	return strconv.FormatUint(uint64(i), 10)
}

func FormatInt[T constraints.SignedInteger](i T) string {
	// Converting to int64 is always safe as i is *at most* 64-bit
	return strconv.FormatInt(int64(i), 10)
}

func FormatFloat[T constraints.Float](f T) string {
	return fmt.Sprintf("%.2f", f)
}
