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

func DivideUint[A, B constraints.UnsignedInteger](a A, b B) uint64 {
	if b == 0 {
		return uint64(a)
	}
	return uint64(a) / uint64(b)
}

func DivideFloat[A, B constraints.Integer](a A, b B) float64 {
	// Checking for 0 explicitly rather than using max(b, 1) since b could be negative
	if b == 0 {
		return float64(a)
	}
	return float64(a) / float64(b)
}
