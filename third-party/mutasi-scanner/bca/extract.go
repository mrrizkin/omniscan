package bca

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
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
func processPdf(pdfR *pdf.Reader) (Transactions, error) {
	totalPage := pdfR.NumPage()
	transactions := make([]*Transaction, 0)
	var currentTransaction *Transaction = nil
	var isNew = false
	year := "1900"
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := pdfR.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		unsortedRows, _ := p.GetTextByRow()
		sortedRows := make(rowSortable, len(unsortedRows))
		copy(sortedRows, unsortedRows)
		sort.Sort(sortedRows)
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
					if year == "1900" {
						for _, m := range months {
							if strings.Contains(word.S, m) {
								yearStr := strings.TrimPrefix(word.S, m+" ")
								if yearRegex.MatchString(yearStr) {
									_, err := strconv.ParseInt(yearStr, 10, 32)
									if err != nil {
										year = "1900"
									} else {
										year = yearStr
									}
								}
							}
						}
					}
					if strings.Contains("TANGGAL", word.S) && wordIndex == 0 {
						m++
					}
					if strings.Contains("KETERANGAN", word.S) && wordIndex == 1 {
						m++
					}
					if strings.Contains("CBG", word.S) && wordIndex == 2 {
						m++
					}
					if strings.Contains("MUTASI", word.S) && wordIndex == 3 {
						m++
					}
					if strings.Contains("SALDO", word.S) && wordIndex == 4 {
						m++
					}
					if m == 5 {
						aftTanggal = true
					}
				}
			}
		}
	}
	return transactions, nil
}
