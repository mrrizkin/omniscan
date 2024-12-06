package pdfextract

import (
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf-extract/encoder"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type (
	fontObject struct {
		*model.FontObject
		cmap *encoder.CMap
	}

	fontObjects map[string]*fontObject
)

func (fo *fontObject) GetCharacterMap(ctx *model.Context) error {
	toUnicode, ok := fo.FontDict.Find("ToUnicode")
	if !ok {
		return nil
	}

	toUnicodeStream, valid, err := ctx.DereferenceStreamDict(toUnicode)
	if err != nil {
		return err
	}

	if !valid {
		return nil
	}

	err = toUnicodeStream.Decode()
	if err != nil {
		return err
	}

	cmap, err := encoder.ParseCmap(toUnicodeStream.Content)
	if err != nil {
		return err
	}

	if fo.cmap == nil {
		fo.cmap = cmap
	} else {
		fo.cmap.Merge(cmap)
	}

	return nil
}

func (fo *fontObject) Decode(raw string) string {
	if fo.cmap != nil {
		return fo.cmap.Decode(raw)
	}
	return raw
}

func (fo fontObjects) Get(resourceName string) (*fontObject, bool) {
	if strings.Contains(resourceName, "/") {
		resourceName = strings.ReplaceAll(resourceName, "/", "")
	}

	font, ok := fo[resourceName]
	if !ok {
		for name, f := range fo {
			if strings.Contains(name, resourceName) {
				font = f
				ok = true
				break
			}
		}
	}
	return font, ok
}
