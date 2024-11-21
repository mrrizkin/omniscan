package pdfextract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type PDFMetadata struct {
	outdirpath string
	filename   string
	metadata   Metadata
}

type Metadata struct {
	*pdfcpu.PDFInfo

	XMLMetadata []map[string]interface{} `json:"xmlmetadata"`
}

func NewPDFMetadata(f []byte, filename string) (*PDFMetadata, error) {
	if path.Ext(filename) != ".pdf" {
		return nil, fmt.Errorf("invalid file extension")
	}

	outdirpath := getPath(path.Join("metadata", filename))
	exist := directoryExists(outdirpath)
	if !exist {
		err := os.MkdirAll(outdirpath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	reader := bytes.NewReader(f)
	err := api.ExtractMetadata(reader, outdirpath, filename, nil)
	if err != nil {
		return nil, err
	}

	pdfInfo, err := api.PDFInfo(reader, filename, nil, nil)
	if err != nil {
		return nil, err
	}

	return &PDFMetadata{
		filename:   strings.TrimSuffix(filename, ".pdf"),
		outdirpath: outdirpath,
		metadata: Metadata{
			PDFInfo:     pdfInfo,
			XMLMetadata: make([]map[string]interface{}, 0),
		},
	}, nil
}

func (p *PDFMetadata) ParseXMLMetadata() error {
	files, err := os.ReadDir(p.outdirpath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("skipping directory: %s\n", file.Name())
			continue
		}

		metadata, err := p.parseXMLMetadata(file.Name())
		if err != nil {
			return err
		}

		p.metadata.XMLMetadata = append(p.metadata.XMLMetadata, metadata)
	}

	return nil
}

func (p *PDFMetadata) Metadata() Metadata {
	return p.metadata
}

func (p *PDFMetadata) Close() error {
	err := os.RemoveAll(p.outdirpath)
	if err != nil {
		return err
	}

	return nil
}

func (p *PDFMetadata) parseXMLMetadata(filename string) (map[string]interface{}, error) {
	// Construct full file path
	fullPath := path.Join(p.outdirpath, filename)

	// Read the entire file content
	xmlContent, err := os.ReadFile(fullPath)
	if err != nil {
		// Log the error but continue (return empty map)
		return nil, err
	}

	// Create a generic map to store metadata
	metadata := make(map[string]interface{})

	reader := bytes.NewReader(xmlContent)
	jsonByte, err := xj.Convert(reader)
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
