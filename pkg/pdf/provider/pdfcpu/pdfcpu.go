package pdfcpu

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/mrrizkin/omniscan/pkg/pdf"
	"github.com/mrrizkin/omniscan/pkg/pdf/types"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type (
	PDFCPU struct {
		ctx   *model.Context
		fonts fonts
	}

	Page struct {
		content []TextObject
	}
)

func NewReader(filename string, b []byte) (*PDFCPU, error) {
	r := bytes.NewReader(b)
	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadValidateAndOptimize(r, conf)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	fonts, err := GetFonts(ctx)
	if err != nil {
		return nil, err
	}

	return &PDFCPU{
		ctx:   ctx,
		fonts: fonts,
	}, nil
}

func (p *PDFCPU) Page(page int) (pdf.PDFPage, error) {
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

func (p *PDFCPU) NumPage() int {
	return p.ctx.PageCount
}

func (p *Page) GetTextByRow(tolerance float64) (types.Rows, error) {
	row := make(types.Rows, 0)
	var currentPosition int64 = 0
	rowIndex := -1
	for _, object := range p.content {
		if int64(object.Position.Y) == currentPosition {
			row = append(row, &types.Row{
				Content: types.TextHorizontal{{
					Font:     object.FontName,
					FontSize: object.FontSize,
					X:        object.Position.X,
					Y:        object.Position.Y,
					S:        object.Text,
				}},
				Position: int64(object.Position.Y),
			})
			currentPosition = int64(object.Position.Y)
			rowIndex++
		} else {
			if rowIndex == -1 {
				continue
			}
			row[rowIndex].Content = append(row[rowIndex].Content, types.Text{
				Font:     object.FontName,
				FontSize: object.FontSize,
				X:        object.Position.X,
				Y:        object.Position.Y,
				S:        object.Text,
			})
		}
	}

	sort.Sort(row)
	return row, nil
}
