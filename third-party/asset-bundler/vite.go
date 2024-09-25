package assetbundler

import (
	goviteparser "github.com/mrrizkin/go-vite-parser"
)

type vite struct {
	manifest *goviteparser.ViteManifestInfo
}

func newVite(config *goviteparser.Config) *vite {
	manifest := goviteparser.Parse(*config)
	return &vite{
		manifest: &manifest,
	}
}

func (v *vite) Entry(entries ...string) string {
	if v.manifest.IsDev() {
		return v.manifest.RenderDevEntriesTag(entries...)
	}

	return v.manifest.RenderEntriesTag(entries...)
}

func (v *vite) ReactRefresh() string {
	if v.manifest.IsDev() {
		return v.manifest.RenderReactRefreshTag()
	}

	return ""
}
