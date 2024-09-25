package assetbundler

import (
	goviteparser "github.com/mrrizkin/go-vite-parser"
)

func Vite(config *goviteparser.Config) *vite {
	return newVite(config)
}
