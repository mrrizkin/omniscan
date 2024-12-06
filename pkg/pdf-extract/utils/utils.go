package utils

import (
	"math"
	"strconv"
)

func Hex2Bytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	result := make([]byte, len(hexStr)/2)
	for i := 0; i < len(result); i++ {
		val, err := strconv.ParseInt(hexStr[i*2:i*2+2], 16, 16)
		if err != nil {
			return nil, err
		}
		result[i] = byte(val)
	}
	return result, nil
}

func IsEqualTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func IsUTF8(s string) bool {
	for _, r := range s {
		if r < 0x80 {
			continue
		}
		return false
	}
	return true
}

func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for chunkSize < len(slice) {
		slice, chunks = slice[chunkSize:], append(chunks, slice[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, slice)
	return chunks
}
