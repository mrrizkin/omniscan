package estatementscanner

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/bca"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

type scannerSet = map[string]types.EStatementScanner

type EStatementScanner struct {
	scanner scannerSet
}

func New() *EStatementScanner {
	scanner := scannerSet{
		"bca": bca.New(),
	}

	return &EStatementScanner{scanner}
}

func (ms *EStatementScanner) Scan(provider string, input interface{}) (*types.ScanResult, error) {
	scanner, ok := ms.scanner[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	switch inputType.Kind() {
	case reflect.String:
		return scanner.ProcessFromPath(input.(string))
	case reflect.Slice:
		if inputType.Elem().Kind() == reflect.Uint8 {
			return scanner.ProcessFromBytes(input.([]byte))
		}
	case reflect.Ptr:
		if inputType.Elem() == reflect.TypeOf((*os.File)(nil)).Elem() {
			return scanner.ProcessFromReader(input.(io.Reader))
		}
	}

	return nil, fmt.Errorf("unsupported input type: %v", inputType)
}
