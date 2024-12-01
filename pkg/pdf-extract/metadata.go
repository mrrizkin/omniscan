package pdfextract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	xj "github.com/basgys/goxml2json"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Metadata struct {
	*pdfcpu.PDFInfo

	XMLMetadata []map[string]interface{} `json:"xmlmetadata"`
}

func NewMetadata(f []byte, filename string) (*Metadata, error) {
	reader := bytes.NewReader(f)
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXTRACTMETADATA
	ctx, err := api.ReadValidateAndOptimize(reader, conf)
	if err != nil {
		return nil, err
	}

	mm, err := pdfcpu.ExtractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	xmlMetadata := make([]map[string]interface{}, len(mm))
	for i, m := range mm {
		fname := fmt.Sprintf("%s_Metadata_%s_%d_%d.txt", filename, m.ParentType, m.ParentObjNr, m.ObjNr)
		metadata, err := parseXMLMetadata(m, fname)
		if err != nil {
			return nil, err
		}
		xmlMetadata[i] = metadata
	}

	pdfInfo, err := api.PDFInfo(reader, filename, nil, nil)
	if err != nil {
		return nil, err
	}

	return &Metadata{
		PDFInfo:     pdfInfo,
		XMLMetadata: xmlMetadata,
	}, nil
}

func parseXMLMetadata(r io.Reader, filename string) (map[string]interface{}, error) {
	// Create a generic map to store metadata
	metadata := make(map[string]interface{})
	jsonByte, err := xj.Convert(r)
	if err != nil {
		return nil, err
	}

	// Add filename to metadata for reference
	metadata["filename"] = filename

	// Unmarshal the JSON data into the generic map
	err = json.Unmarshal([]byte(jsonByte.String()), &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
