package utils

import gonanoid "github.com/matoous/go-nanoid"

func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for chunkSize < len(slice) {
		slice, chunks = slice[chunkSize:], append(chunks, slice[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, slice)
	return chunks
}

func RandomStr(length int) string {
	str, _ := gonanoid.Generate(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
		length,
	)
	return str
}

func InArrayString(val string, array []string) bool {
	for _, v := range array {
		if v == val {
			return true
		}
	}
	return false
}
