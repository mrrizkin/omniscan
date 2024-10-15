package mutasi_scanner

import (
	"fmt"

	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func ExtractMutasiMetadataFile(file string) (*types.PDFMetadata, error) {
	conf := model.NewDefaultConfiguration()
	ctx, err := pdfcpu.ReadFile(file, conf)
	if err != nil {
		return nil, err
	}

	return extractMetadataToMap(ctx)
}

func extractMetadataToMap(ctx *model.Context) (*types.PDFMetadata, error) {
	if ctx == nil {
		return nil, fmt.Errorf("metadata-extractor: context model is nil")
	}

	var metadata types.PDFMetadata

	if ctx.XRefTable != nil {
		infoDict := ctx.XRefTable

		metadata.Title = infoDict.Title
		metadata.Author = infoDict.Author
		metadata.Subject = infoDict.Subject
		metadata.Keywords = infoDict.Keywords
		metadata.Creator = infoDict.Creator
		metadata.Producer = infoDict.Producer
		metadata.CreationDate = infoDict.CreationDate
		metadata.ModDate = infoDict.ModDate
	}

	if ctx.Root != nil {
		metadata.PageCount = ctx.PageCount

		if version := ctx.RootVersion; version != nil {
			metadata.PDFVersion = version.String()
		}
	}

	return &metadata, nil
}
