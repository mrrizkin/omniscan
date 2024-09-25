package bca

import (
	"bytes"
	"io"

	"github.com/ledongthuc/pdf"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

type BCA struct{}

func New() types.MutasiScanner {
	return &BCA{}
}

func (bca *BCA) ProcessFromPath(path string) (*types.Transactions, error) {
	f, pdfR, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	trx, err := processPdf(pdfR)
	if err != nil {
		return nil, err
	}

	return bca.maptrx(trx), nil
}
func (bca *BCA) ProcessFromReader(r io.Reader) (*types.Transactions, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	trx, err := bca.ProcessFromBytes(b)
	if err != nil {
		return nil, err
	}

	return trx, nil
}
func (bca *BCA) ProcessFromBytes(b []byte) (*types.Transactions, error) {
	bytesR := bytes.NewReader(b)
	pdfR, err := pdf.NewReader(bytesR, bytesR.Size())
	if err != nil {
		return nil, err
	}
	trx, err := processPdf(pdfR)
	if err != nil {
		return nil, err
	}

	return bca.maptrx(trx), nil
}

func (*BCA) maptrx(trxs Transactions) *types.Transactions {
	trx := new(types.Transactions)
	trx.Transactions = make([]*types.Transaction, len(trxs))
	startIndex := 0
	lastIndex := len(trxs) - 1
	for i, t := range trxs {
		if i == startIndex {
			trx.StartBalance = t.Balance
		}

		if i == lastIndex {
			trx.EndBalance = t.Balance
		}

		trxType := ""
		if t.DirectionCr != nil {
			if *t.DirectionCr {
				trxType = "credit"
				trx.TransactionCreditCount = trx.TransactionCreditCount + 1
				trx.TransactionCreditTotal = trx.TransactionCreditTotal + t.Change
			} else {
				trxType = "debit"
				trx.TransactionDebitCount = trx.TransactionDebitCount + 1
				trx.TransactionDebitTotal = trx.TransactionDebitTotal + t.Change
			}
		}

		trx.Transactions[i] = &types.Transaction{
			Date:            t.Date,
			Description1:    t.Description1,
			Description2:    t.Description2,
			Branch:          t.Branch,
			Change:          t.Change,
			TransactionType: trxType,
			Balance:         t.Balance,
		}
	}

	return trx
}
