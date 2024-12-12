package mandiri

import (
	"strconv"
	"strings"
	"time"

	"github.com/mrrizkin/omniscan/pkg/pdf/types"
)

type Transaction struct {
	Date         time.Time `json:"date,omitempty"`
	Description1 string    `json:"description1,omitempty"`
	Description2 string    `json:"description2,omitempty"`
	Branch       string    `json:"branch,omitempty"`
	Change       float64   `json:"change,omitempty"`
	DirectionCr  *bool     `json:"directionCr,omitempty"`
	Balance      float64   `json:"balance,omitempty"`
}

type Transactions []*Transaction

type Header struct {
	Product  string
	Rekening string
	Periode  string
}

type leftCol float64
type rightCol float64

func (c leftCol) Is(x float64) bool {
	diff := c - leftCol(x)
	if diff < 0 {
		diff = diff * -1
	}
	return diff < 5.0
}
func (c rightCol) Is(x float64) bool {
	return rightCol(x) > c
}

const dateCol leftCol = 46.04
const changeAmountDebitCol rightCol = 400.0
const changeAmountCreditCol rightCol = 500.0
const descriptionCol leftCol = 99.61
const summaryFirstCol leftCol = 99.61

// a row with a new date signifies a new transaction
//
// as a PDF text row might not be a new transaction,
// but adds detail to the previous transaction, prevT
// is added as argument to add transaction detail
//
// the returned isNew tells whether the returned
// *transaction is a new transaction that should
// be added to a transaction slice
func IngestRow(
	prevT *Transaction,
	row *types.Row,
	year string,
) (isNew bool, t *Transaction, shouldStopProcessing bool) {
	words := make(types.TextHorizontal, len(row.Content))
	copy(words, row.Content)
	if len(words) < 2 {
		if prevT == nil {
			return
		}
		t = prevT
		shouldStopProcessing = readSupplementary(t, words)
		return
	}
	firstWord := words[0]
	date, dateErr := time.Parse("02/01/2006", firstWord.S+"/"+year)
	hasDate := dateErr == nil && dateCol.Is(firstWord.X)
	if !hasDate {
		if prevT == nil {
			return
		}
		t = prevT
	} else {
		isNew = true
		words = words[1:]
		t = &Transaction{
			Date: date,
		}
	}

	shouldStopProcessing = readSupplementary(t, words)
	return
}

// readSupplementary try to read words in a row that are apart of
// date and balance information
func readSupplementary(t *Transaction, words types.TextHorizontal) (stopProcessingNext bool) {
	for i, word := range words {
		if i == 0 && word.S == "Saldo Awal" && summaryFirstCol.Is(word.X) {
			return true
		}
		if changeAmountDebitCol.Is(word.X) {
			amount, amountErr := strconv.ParseFloat(
				strings.ReplaceAll(word.S, ",", ""),
				32)
			if amountErr == nil {
				if t.DirectionCr == nil {
					isCr := false
					t.DirectionCr = &isCr
				}
				t.Change = amount
			}
		}
		if changeAmountCreditCol.Is(word.X) {
			amount, amountErr := strconv.ParseFloat(
				strings.ReplaceAll(word.S, ",", ""),
				32)
			if amountErr == nil {
				if t.DirectionCr == nil {
					isCr := true
					t.DirectionCr = &isCr
				}
				t.Change = amount
			}
		}
		if descriptionCol.Is(word.X) {
			if t.Description1 == "" {
				t.Description1 = word.S
			} else {
				t.Description1 = t.Description1 + "\n" + word.S
			}
		}
	}
	return
}
