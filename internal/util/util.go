package util

import (
	"strconv"
)

func FormatUint[T uint | uint8 | uint16 | uint32 | uint64](i T) string {
	// Converting to uint64 is always safe as i is *at most* 64-bit
	return strconv.FormatUint(uint64(i), 10)
}

func FormatInt[T int | int8 | int16 | int32 | int64](i T) string {
	// Converting to int64 is always safe as i is *at most* 64-bit
	return strconv.FormatInt(int64(i), 10)
}
