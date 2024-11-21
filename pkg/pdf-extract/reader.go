package pdfextract

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFReader struct {
	outdirpath string
	filename   string
	totalPage  int
}

func NewPDFReader(f []byte, filename string) (*PDFReader, error) {
	if path.Ext(filename) != ".pdf" {
		return nil, fmt.Errorf("invalid file extension")
	}

	outdirpath := getPath(path.Join("reader", filename))
	exist := directoryExists(outdirpath)
	if !exist {
		err := os.MkdirAll(outdirpath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	reader := bytes.NewReader(f)
	err := api.ExtractContent(reader, outdirpath, filename, nil, nil)
	if err != nil {
		return nil, err
	}

	totalPage, err := api.PageCount(reader, nil)
	if err != nil {
		return nil, err
	}

	return &PDFReader{
		filename:   strings.TrimSuffix(filename, ".pdf"),
		outdirpath: outdirpath,
		totalPage:  totalPage,
	}, nil
}

func (p *PDFReader) Page(page int) (*Page, error) {
	file, err := os.Open(path.Join(p.outdirpath, fmt.Sprintf("%s_Content_page_%d.txt", p.filename, page)))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := parser(file)
	if err != nil {
		return nil, err
	}

	return &Page{content: content}, nil
}

func (p *PDFReader) NumPage() int {
	return p.totalPage
}

func (p *PDFReader) Close() error {
	err := os.RemoveAll(p.outdirpath)
	if err != nil {
		return err
	}

	return nil
}
