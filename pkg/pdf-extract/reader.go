package pdfextract

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type fontObjects map[string]*model.FontObject

func (fo fontObjects) Get(resourceName string) (*model.FontObject, bool) {
	if strings.Contains(resourceName, "/") {
		resourceName = strings.ReplaceAll(resourceName, "/", "")
	}

	font, ok := fo[resourceName]
	if !ok {
		for name, f := range fo {
			if strings.Contains(name, resourceName) {
				font = f
				ok = true
				break
			}
		}
	}
	return font, ok
}

type PDFReader struct {
	outdirpath string
	filename   string
	ctx        *model.Context
	fonts      fontObjects
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
	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadValidateAndOptimize(reader, conf)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	fonts := make(fontObjects)
	for _, font := range ctx.Optimize.FontObjects {
		fonts[font.ResourceNamesString()] = font
	}

	return &PDFReader{
		filename:   strings.TrimSuffix(filename, ".pdf"),
		outdirpath: outdirpath,
		ctx:        ctx,
		fonts:      fonts,
	}, nil
}

func (p *PDFReader) Page(page int) (*Page, error) {
	r, err := pdfcpu.ExtractPageContent(p.ctx, page)
	if err != nil {
		return nil, err
	}

	content, err := p.parse(r)
	if err != nil {
		return nil, err
	}

	return &Page{content: content}, nil
}

func (p *PDFReader) NumPage() int {
	return p.ctx.PageCount
}

func (p *PDFReader) Close() error {
	err := os.RemoveAll(p.outdirpath)
	if err != nil {
		return err
	}

	return nil
}
