package pdfextract

import (
	"math"
	"os"
	"path"
)

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && info.IsDir()
}

func isEqualTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func getPath(filepath string) string {
	return path.Join("./storage/app/pdf-extract/", filepath)
}
