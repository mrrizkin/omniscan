package estatementscanner

import (
	"fmt"

	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/bca"
	bcaledongthuc "github.com/mrrizkin/omniscan/pkg/e-statement-scanner/bca-ledongthuc"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

type scannerSet = map[string]types.EStatementScanner

type EStatementScanner struct {
	scanner scannerSet
}

func New() *EStatementScanner {
	scanner := scannerSet{
		"bca":            bca.New(),
		"bca-ledongthuc": bcaledongthuc.New(),
	}

	return &EStatementScanner{scanner}
}

func (ms *EStatementScanner) Scan(provider, filename string, input []byte) (*types.ScanResult, error) {
	scanner, ok := ms.scanner[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return scanner.ProcessFromBytes(filename, input)
}
