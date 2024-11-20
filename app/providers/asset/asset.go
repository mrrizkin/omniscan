package asset

import (
	goviteparser "github.com/mrrizkin/go-vite-parser"
)

type Asset struct {
	vite *goviteparser.ViteManifestInfo
}

func (*Asset) Construct() interface{} {
	return func() *Asset {
		manifest := goviteparser.Parse(goviteparser.Config{
			OutDir:       "/build/",
			ManifestPath: "public/build/manifest.json",
			HotFilePath:  "public/hot",
		})

		return &Asset{
			vite: &manifest,
		}
	}
}

func (a *Asset) Entry(entries ...string) string {
	if a.vite.IsDev() {
		return a.vite.RenderDevEntriesTag(entries...)
	}

	return a.vite.RenderEntriesTag(entries...)
}

func (a *Asset) ReactRefresh() string {
	if a.vite.IsDev() {
		return a.vite.RenderReactRefreshTag()
	}

	return ""
}
