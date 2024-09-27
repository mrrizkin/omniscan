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

func (bca *BCA) ProcessFromPath(path string) (*types.ScanResult, error) {
	f, pdfR, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	trx, header, err := processPdf(pdfR)
	if err != nil {
		return nil, err
	}

	return bca.maptrx(header, trx), nil
}
func (bca *BCA) ProcessFromReader(r io.Reader) (*types.ScanResult, error) {
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
func (bca *BCA) ProcessFromBytes(b []byte) (*types.ScanResult, error) {
	bytesR := bytes.NewReader(b)
	pdfR, err := pdf.NewReader(bytesR, bytesR.Size())
	if err != nil {
		return nil, err
	}
	trx, header, err := processPdf(pdfR)
	if err != nil {
		return nil, err
	}

	return bca.maptrx(header, trx), nil
}

func (*BCA) maptrx(header Header, trxs Transactions) *types.ScanResult {
	res := new(types.ScanResult)
	res.Transactions = make([]*types.Transaction, len(trxs))
	totalBalance := 0.0
	for i, t := range trxs {
		totalBalance = totalBalance + t.Balance

		trxType := ""
		if t.DirectionCr != nil {
			if *t.DirectionCr {
				trxType = "credit"
			} else {
				trxType = "debit"
			}
		}

		res.Transactions[i] = &types.Transaction{
			Date:            t.Date,
			Description1:    t.Description1,
			Description2:    t.Description2,
			Branch:          t.Branch,
			Change:          t.Change,
			TransactionType: trxType,
			Balance:         t.Balance,
		}
	}

	countTransaction := len(trxs)
	if countTransaction < 1 {
		countTransaction = 1
	}

	res.Info.Bank = "BCA"
	res.Info.Produk = header.Product
	res.Info.Rekening = header.Rekening
	res.Info.Periode = header.Periode

	return res
}
