package pdfextract

import (
	"bytes"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Reader struct {
	ctx   *model.Context
	fonts fontObjects
}

func NewReader(f []byte) (*Reader, error) {
	r := bytes.NewReader(f)
	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadValidateAndOptimize(r, conf)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	fonts := make(fontObjects)
	for _, font := range ctx.Optimize.FontObjects {
		for _, resName := range font.ResourceNames {
			fo, ok := fonts[resName]
			if !ok {
				fo = &fontObject{FontObject: font}
			} else {
				fo.FontObject.FontDict = font.FontDict
			}

			err := fo.GetCharacterMap(ctx)
			if err != nil {
				return nil, err
			}

			fonts[resName] = fo
		}
	}

	return &Reader{
		ctx:   ctx,
		fonts: fonts,
	}, nil
}

func (p *Reader) Page(page int) (*Page, error) {
	r, err := pdfcpu.ExtractPageContent(p.ctx, page)
	if err != nil {
		return nil, err
	}

	content, err := p.parse(r)
	if err != nil {
		return nil, err
	}

	return &Page{content: content}, nil
}

func (p *Reader) NumPage() int {
	return p.ctx.PageCount
}
