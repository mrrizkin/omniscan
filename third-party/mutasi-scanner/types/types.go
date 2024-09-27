package types

import (
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
