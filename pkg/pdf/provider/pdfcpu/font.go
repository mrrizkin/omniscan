package pdfcpu

import (
	"fmt"
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf/provider/pdfcpu/encoder"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type (
	Position struct {
		X float64
		Y float64
	}

	TextObject struct {
		FontName string
		Text     string
		FontSize float64
		Position Position
	}

	fontObject struct {
		*model.FontObject
		enc encoder.TextEncoding
	}

	fonts map[string]*fontObject
)

func (fo *fontObject) Decode(raw string) string {
	if fo.enc != nil {
		return fo.enc.Decode(raw)
	}
	return raw
}

func (fs fonts) Get(resourceName string) (*fontObject, bool) {
	resourceName = strings.TrimPrefix(resourceName, "/")
	font, ok := fs[resourceName]
	if !ok {
		for name, f := range fs {
			if strings.Contains(name, resourceName) {
				font = f
				ok = true
				break
			}
		}
	}
	return font, ok
}

func (fo fonts) Set(resourceName string, font *fontObject) {
	fo[resourceName] = font
}

func GetFonts(ctx *model.Context) (fonts, error) {
	fonts := make(fonts)
	for _, font := range ctx.Optimize.FontObjects {
		encoding := ""
		var enc encoder.TextEncoding
		if encodingDict, ok := font.FontDict.Find("Encoding"); ok {
			encoding = encodingDict.String()
		}

		switch encoding {
		case "WinAnsiEncoding":
			enc = encoder.NewWinAnsiEncoding()
		case "MacRomanEncoding":
			enc = encoder.NewMacRomanEncoding()
		case "Identity-H":
			cmap, err := getCMap(ctx, font.FontDict)
			if err != nil {
				return nil, err
			}

			enc = cmap
		}

		for _, resName := range font.ResourceNames {
			fo, ok := fonts[resName]
			if !ok {
				fo = &fontObject{
					FontObject: font,
					enc:        enc,
				}
			} else {
				if fo.enc == nil {
					fo.enc = enc
				}

				if cmap, ok := enc.(*encoder.CMap); ok {
					cmap.Merge(fo.enc.(*encoder.CMap))
					fo.enc = cmap
				}
			}

			fonts[resName] = fo
		}
	}

	return fonts, nil
}

func getCMap(ctx *model.Context, fontDict types.Dict) (*encoder.CMap, error) {
	toUnicode, ok := fontDict.Find("ToUnicode")
	if !ok {
		return nil, fmt.Errorf("ToUnicode not found")
	}

	stream, valid, err := ctx.DereferenceStreamDict(toUnicode)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("invalid ToUnicode stream")
	}

	cmap, err := encoder.ParseCmap(stream.Content)
	if err != nil {
		return nil, err
	}

	return cmap, nil
}
