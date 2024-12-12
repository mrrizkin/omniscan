package mandiri

import (
	"regexp"
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf"
)

var yearRegex = regexp.MustCompile(`\d\d\d\d`)
var months = []string{
	"JANUARI",
	"FEBRUARI",
	"MARET",
	"APRIL",
	"MEI",
	"JUNI",
	"JULI",
	"AGUSTUS",
	"SEPTEMBER",
	"OKTOBER",
	"NOVEMBER",
	"DESEMBER",
}

// this is the internal function called by the exported
// ProcessPdf*** functions
func processPdf(pdfR pdf.PDFReader) (Transactions, Header, error) {
	totalPage := pdfR.NumPage()
	transactions := make([]*Transaction, 0)
	var currentTransaction *Transaction = nil
	var isNew = false
	year := "1900"
	header := Header{
		Product:  "",
		Rekening: "",
		Periode:  "",
	}
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p, err := pdfR.Page(pageIndex)
		if err != nil {
			return nil, header, err
		}

		sortedRows, err := p.GetTextByRow(4)
		if err != nil {
			return nil, header, err
		}
		aftTanggal := false
		shouldStopProcessing := false
		for _, row := range sortedRows {
			if aftTanggal {
				isNew, currentTransaction, shouldStopProcessing = IngestRow(
					currentTransaction,
					row,
					year,
				)
				if isNew {
					transactions = append(transactions, currentTransaction)
				}
				if shouldStopProcessing {
					break
				}
			} else {
				// here we try to ignore statement end-footer
				m := 0
				for wordIndex, word := range row.Content {
					if pageIndex == 1 {
						if strings.Contains(word.S, "REKENING") && wordIndex == 0 {
							if len(row.Content) == 1 {
								header.Product = word.S
							}
						}
						if strings.Contains(word.S, "Nomor Rekening") && wordIndex == 0 {
							if len(row.Content) > 1 {
								txt := row.Content[1]
								header.Rekening = txt.S
							}
						}
						if strings.Contains(word.S, "Periode") && wordIndex == 0 {
							if len(row.Content) > 3 {
								text := []string{}
								for _, txt := range row.Content[1:] {
									text = append(text, txt.S)
								}
								header.Periode = strings.Join(text, " ")
								toDate := strings.Split(row.Content[3].S, "/")
								year = toDate[len(toDate)-1]
							}
						}
					}
					if strings.Contains("TANGGAL", word.S) && wordIndex == 0 {
						m++
					}
					if strings.Contains("TRANSAKSI", word.S) && wordIndex == 1 {
						m++
					}
					if strings.Contains("DEBIT", word.S) && wordIndex == 2 {
						m++
					}
					if strings.Contains("KREDIT", word.S) && wordIndex == 3 {
						m++
					}
					if m == 4 {
						aftTanggal = true
					}
				}
			}
		}
	}
	return transactions, header, nil
}
