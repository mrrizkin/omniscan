package mutasi_scanner

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
)

func ExtractMutasiMetadataFile(file string) (map[string]interface{}, error) {
	log.SetCLILogger(&logger{})
	err := api.ExtractContentFile(
		"storage/7771953151_Agu_2019.pdf",
		"storage",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]interface{})
	return metadata, nil
}

type logger struct{}

func (*logger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*logger) Println(args ...interface{}) {
	fmt.Println(args...)
}

func (*logger) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*logger) Fatalln(args ...interface{}) {
	fmt.Println(args...)
}
