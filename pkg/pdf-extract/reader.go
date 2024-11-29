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

type fontObject struct {
	*model.FontObject

	decoder *ToUnicodeDecoder
}

func (font *fontObject) ToUnicode(ctx *model.Context) error {
	toUnicode, ok := font.FontDict.Find("ToUnicode")
	if !ok {
		return nil
	}

	toUnicodeStream, valid, err := ctx.DereferenceStreamDict(toUnicode)
	if err != nil {
		return err
	}

	if !valid {
		return nil
	}

	err = toUnicodeStream.Decode()
	if err != nil {
		return err
	}

	decoder, err := NewToUnicodeDecoder(toUnicodeStream.Content)
	if err != nil {
		return err
	}

	font.decoder = decoder

	return nil
}

func (fo *fontObject) Decode(raw string) (text string) {
	if fo.decoder != nil {
		return fo.decoder.Decode([]byte(raw))
	}
	return raw
}

type fontObjects map[string]*fontObject

func (fo fontObjects) Get(resourceName string) (*fontObject, bool) {
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
		fo := fontObject{FontObject: font}
		fo.ToUnicode(ctx)
		fonts[font.ResourceNamesString()] = &fo
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

	// TODO: I don't know why pdf of each page is generated so we need to delete them
	if files, err := os.ReadDir("./storage"); err == nil {
		for _, file := range files {
			if strings.Contains(file.Name(), p.filename) {
				err := os.Remove("./storage/" + file.Name())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
