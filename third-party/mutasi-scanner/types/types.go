package types

import (
	"encoding/json"
	"io"
	"time"
)

type MutasiScanner interface {
	ProcessFromPath(path string) (*ScanResult, error)
	ProcessFromReader(r io.Reader) (*ScanResult, error)
	ProcessFromBytes(b []byte) (*ScanResult, error)
}

type Transaction struct {
	Date            time.Time `json:"date,omitempty"`
	Description1    string    `json:"description1,omitempty"`
	Description2    string    `json:"description2,omitempty"`
	Branch          string    `json:"branch,omitempty"`
	Change          float64   `json:"change,omitempty"`
	TransactionType string    `json:"transaction_type,omitempty"`
	Balance         float64   `json:"balance,omitempty"`
}

type ScanInfo struct {
	Bank     string `json:"bank"`
	Produk   string `json:"produk"`
	Rekening string `json:"rekening"`
	Periode  string `json:"periode"`
}

type ScanResult struct {
	Info         ScanInfo       `json:"info"`
	Transactions []*Transaction `json:"transactions"`
}

type PDFMetadata struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Subject      string `json:"subject"`
	Keywords     string `json:"keywords"`
	Creator      string `json:"creator"`
	Producer     string `json:"producer"`
	CreationDate string `json:"creation_date"`
	ModDate      string `json:"mod_date"`
	PageCount    int    `json:"page_count"`
	PDFVersion   string `json:"pdf_version"`
}

func (p *PDFMetadata) String() string {
	encode, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(encode)
}
