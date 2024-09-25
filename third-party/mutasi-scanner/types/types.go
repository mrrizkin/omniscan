package types

import (
	"io"
	"time"
)

type MutasiScanner interface {
	ProcessFromPath(path string) (*Transactions, error)
	ProcessFromReader(r io.Reader) (*Transactions, error)
	ProcessFromBytes(b []byte) (*Transactions, error)
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

type Transactions struct {
	StartBalance           float64        `json:"start_balance"`
	EndBalance             float64        `json:"end_balance"`
	TransactionDebitTotal  float64        `json:"transaction_debit_total"`
	TransactionCreditTotal float64        `json:"transaction_credit_total"`
	TransactionDebitCount  float64        `json:"transaction_debit_count"`
	TransactionCreditCount float64        `json:"transaction_credit_count"`
	Transactions           []*Transaction `json:"transactions"`
}
