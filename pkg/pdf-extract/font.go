package pdfextract

import (
	"fmt"
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf-extract/encoder"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type (
	Position struct {
		X float64
		Y float64
	}

	TextObject struct {
		FontName     string
		ResourceName string
		Text         string
		FontSize     float64
		Position     Position
	}

	Text struct {
		Font     string
		FontSize float64
		X        float64
		Y        float64
		S        string
	}

	TextEncoder interface {
		Decode(raw string) string
	}

	fontObject struct {
		*model.FontObject
		encoder TextEncoder
	}

	fontObjects map[string]*fontObject
)

func (fo *fontObject) Decode(raw string) string {
	if fo.encoder != nil {
		return fo.encoder.Decode(raw)
	}
	return raw
}

func (fo fontObjects) Get(resourceName string) (*fontObject, bool) {
	resourceName = strings.TrimPrefix(resourceName, "/")
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

func (fo fontObjects) Set(resourceName string, font *fontObject) {
	fo[resourceName] = font
}

func getCMap(ctx *model.Context, font *model.FontObject) (*encoder.CMap, error) {
	toUnicode, ok := font.FontDict.Find("ToUnicode")
	if !ok {
		return nil, nil
	}

	stream, valid, err := ctx.DereferenceStreamDict(toUnicode)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("invalid ToUnicode stream")
	}

	err = stream.Decode()
	if err != nil {
		return nil, err
	}

	cmap, err := encoder.ParseCmap(stream.Content)
	if err != nil {
		return nil, err
	}

	return cmap, nil
}

func GetFonts(ctx *model.Context) (fontObjects, error) {
	fonts := make(fontObjects)
	for _, font := range ctx.Optimize.FontObjects {
		var textEncoder TextEncoder
		encoding := ""
		if encodingObject, ok := font.FontDict.Find("Encoding"); ok {
			encoding = encodingObject.String()
		}

		switch encoding {
		case "WinAnsiEncoding":
			textEncoder = encoder.NewWinAnsiEncoding()
		case "MacRomanEncoding":
			textEncoder = encoder.NewMacRomanEncoding()
		case "Identity-H":
			cmap, err := getCMap(ctx, font)
			if err != nil {
				return nil, err
			}
			textEncoder = cmap
		}

		for _, resourceName := range font.ResourceNames {
			fo, ok := fonts.Get(resourceName)
			if !ok {
				fo = &fontObject{
					FontObject: font,
					encoder:    textEncoder,
				}
			} else {
				if cmap, ok := fo.encoder.(*encoder.CMap); ok {
					cmap.Merge(textEncoder.(*encoder.CMap))
					fo.encoder = cmap
				}
			}

			fonts.Set(resourceName, fo)
		}
	}

	return fonts, nil
}
