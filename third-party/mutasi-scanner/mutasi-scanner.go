package mutasi_scanner

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/bca"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

type scannerSet = map[string]types.MutasiScanner

type MutasiScanner struct {
	scanner scannerSet
}

func New() *MutasiScanner {
	scanner := scannerSet{
		"bca": bca.New(),
	}

	return &MutasiScanner{scanner}
}

func (ms *MutasiScanner) Scan(provider string, input interface{}) (*types.Transactions, error) {
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
