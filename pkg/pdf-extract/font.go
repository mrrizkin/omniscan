package pdfextract

import (
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type (
	fontObject struct {
		*model.FontObject
		decoder *ToUnicodeDecoder
	}

	fontObjects map[string]*fontObject
)

func (fo *fontObject) ToUnicode(ctx *model.Context) error {
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

	decoder, err := NewToUnicodeDecoder(toUnicodeStream.Content)
	if err != nil {
		return err
	}

	fo.decoder = decoder

	return nil
}

func (fo *fontObject) Decode(raw string) (text string) {
	if fo.decoder != nil {
		return fo.decoder.Decode([]byte(raw))
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
