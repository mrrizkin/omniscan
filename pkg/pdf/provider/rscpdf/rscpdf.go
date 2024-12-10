package rscpdf

import (
	"bytes"
	"sort"

	"github.com/ledongthuc/pdf"
	pdff "github.com/mrrizkin/omniscan/pkg/pdf"
	"github.com/mrrizkin/omniscan/pkg/pdf/types"
)

type RSCPDF struct {
	pdfReader *pdf.Reader
}

type RSCPDFPage struct {
	pdfPage pdf.Page
}

func NewReader(filename string, b []byte) (*RSCPDF, error) {
	reader := bytes.NewReader(b)
	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return nil, err
	}
	return &RSCPDF{
		pdfReader: pdfReader,
	}, nil
}

func (r *RSCPDF) NumPage() int {
	return r.pdfReader.NumPage()
}

func (r *RSCPDF) Page(page int) (pdff.PDFPage, error) {
	return &RSCPDFPage{pdfPage: r.pdfReader.Page(page)}, nil
}

func (r *RSCPDFPage) GetTextByRow(tolerance float64) (types.Rows, error) {
	unsortedRows, _ := r.pdfPage.GetTextByRow()
	sortedRows := make(types.Rows, len(unsortedRows))
	for i, row := range unsortedRows {
		content := make(types.TextHorizontal, len(row.Content))
		for j, text := range row.Content {
			content[j] = types.Text{
				Font:     text.Font,
				FontSize: text.FontSize,
				X:        text.X,
				Y:        text.Y,
				S:        text.S,
			}
		}

		sortedRows[i] = &types.Row{
			Position: row.Position,
			Content:  content,
		}
	}
	sort.Sort(sortedRows)
	return sortedRows, nil
}
