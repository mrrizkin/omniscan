package estatementscanner

import (
	"fmt"

	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/bca"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/mandiri"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
	"github.com/mrrizkin/omniscan/pkg/pdf"
	"github.com/mrrizkin/omniscan/pkg/pdf/provider/pdfcpu"
	"github.com/mrrizkin/omniscan/pkg/pdf/provider/rscpdf"
)

type EStatementScanner struct{}

func New() *EStatementScanner {
	return &EStatementScanner{}
}

func (ms *EStatementScanner) Scan(bank, library, filename string, input []byte) (*types.ScanResult, error) {
	var pdfReader pdf.PDFReader
	var err error

	switch library {
	case "pdfcpu":
		pdfReader, err = pdfcpu.NewReader(filename, input)
		if err != nil {
			return nil, err
		}
	case "rscpdf":
		pdfReader, err = rscpdf.NewReader(filename, input)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported pdf lib: %s", library)
	}

	result := &types.ScanResult{}

	switch bank {
	case "bca":
		result, err = bca.ScanFromBytes(filename, pdfReader)
		if err != nil {
			return nil, err
		}
	case "mandiri":
		result, err = mandiri.ScanFromBytes(filename, pdfReader)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported bank: %s", bank)
	}

	result.Metadata, err = pdf.ExtractMetadata(input, filename)
	if err != nil {
		return nil, err
	}

	return result, nil
}
