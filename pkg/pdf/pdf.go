package pdf

import "github.com/mrrizkin/omniscan/pkg/pdf/types"

type (
	PDFReader interface {
		NumPage() int
		Page(int) (PDFPage, error)
	}

	PDFPage interface {
		GetTextByRow(tolerance float64) (types.Rows, error)
	}
)
