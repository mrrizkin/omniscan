package pdfextract

import (
	"math"
	"path"
)

func isEqualTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func getPath(filepath string) string {
	return path.Join("./storage/app/pdf-extract/", filepath)
}
