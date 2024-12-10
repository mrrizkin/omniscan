package bca

import (
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
	"github.com/mrrizkin/omniscan/pkg/pdf"
)

func ScanFromBytes(filename string, pdfReader pdf.PDFReader) (*types.ScanResult, error) {
	trx, header, err := processPdf(pdfReader)
	if err != nil {
		return nil, err
	}

	return maptrx(header, trx), nil
}

func maptrx(header Header, trxs Transactions) *types.ScanResult {
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
