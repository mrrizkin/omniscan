package bcaledongthuc

import (
	"bytes"

	"github.com/ledongthuc/pdf"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

type BCAOld struct{}

func New() *BCAOld {
	return &BCAOld{}
}

func (bca *BCAOld) ProcessFromBytes(filename string, b []byte) (*types.ScanResult, error) {
	r := bytes.NewReader(b)
	pdfR, err := pdf.NewReader(r, r.Size())
	if err != nil {
		return nil, err
	}
	trx, header, err := processPdf(pdfR)
	if err != nil {
		return nil, err
	}

	return bca.maptrx(header, trx), nil
}

func (*BCAOld) maptrx(header Header, trxs Transactions) *types.ScanResult {
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
