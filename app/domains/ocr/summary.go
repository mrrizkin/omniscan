package ocr

import (
	"sync"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/utils"
)

type summaryField struct {
	name     string
	fetchFn  func() (interface{}, error)
	assignFn func(*OverallSummary, interface{})
}

var monthlyFields = []string{
	"MonthlyStartBalance", "MonthlyAverageBalance", "MonthlyEndBalance",
	"MonthlyTotalIncome", "MonthlyTotalExpenses", "MonthlyAverageDebit",
	"MonthlyAverageCredit", "MonthlyFrequencyDebit",
	"MonthlyFrequencyCredit", "MonthlyTopDebits", "MonthlyTopCredits",
	"MonthlyAnomalyTransactions",
}

var allTimeFields = []string{
	"StartBalance", "AverageBalance", "EndBalance", "TotalIncome",
	"TotalExpenses", "AverageDebit", "AverageCredit", "FrequencyDebit",
	"FrequencyCredit", "TopDebits", "TopCredits", "AnomalyTransactions",
	"TotalBankFee", "TotalInterest", "TotalTax", "TotalDigitalRevenue",
	"TotalTransferIn", "TotalTransferOut", "TotalCashWithdrawal",
}

func skipFetch(field string, names []string) bool {
	if len(names) != 0 {
		if !utils.InArrayString("All", names) {
			if utils.InArrayString("AllTime", names) &&
				!utils.InArrayString(field, allTimeFields) {
				return true
			}

			if utils.InArrayString("AllMonthly", names) &&
				!utils.InArrayString(field, monthlyFields) {
				return true
			}

			if !utils.InArrayString("AllTime", names) &&
				!utils.InArrayString("AllMonthly", names) &&
				!utils.InArrayString(field, names) {
				return true
			}
		}
	}

	return false
}

func (s *Service) getSummary(mutasiID uint, names ...string) (*OverallSummary, error) {
	summary := &OverallSummary{
		Monthly: MonthlySummary{},
		AllTime: Summary{},
	}

	fields := []summaryField{
		// AllTime fields
		{
			"StartBalance",
			func() (interface{}, error) { return s.repo.GetBalance(mutasiID, "start") },
			func(s *OverallSummary, v interface{}) {
				startBalance, ok := v.(float64)
				if ok {
					s.AllTime.StartBalance = startBalance
				}
			},
		},
		{
			"AverageBalance",
			func() (interface{}, error) { return s.repo.GetBalance(mutasiID, "avg") },
			func(s *OverallSummary, v interface{}) {
				avgBalance, ok := v.(float64)
				if ok {
					s.AllTime.AverageBalance = avgBalance
				}
			},
		},
		{
			"EndBalance",
			func() (interface{}, error) { return s.repo.GetBalance(mutasiID, "end") },
			func(s *OverallSummary, v interface{}) {
				endBalance, ok := v.(float64)
				if ok {
					s.AllTime.EndBalance = endBalance
				}
			},
		},
		{
			"TotalIncome",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "total")
			},
			func(s *OverallSummary, v interface{}) {
				totalIncome, ok := v.(float64)
				if ok {
					s.AllTime.TotalIncome = totalIncome
				}
			},
		},
		{
			"TotalExpenses",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "total")
			},
			func(s *OverallSummary, v interface{}) {
				totalExpenses, ok := v.(float64)
				if ok {
					s.AllTime.TotalExpense = totalExpenses
				}
			},
		},
		{
			"AverageDebit",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "avg")
			},
			func(s *OverallSummary, v interface{}) {
				avgDebit, ok := v.(float64)
				if ok {
					s.AllTime.AverageDebit = avgDebit
				}
			},
		},
		{
			"AverageCredit",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "avg")
			},
			func(s *OverallSummary, v interface{}) {
				avgCredit, ok := v.(float64)
				if ok {
					s.AllTime.AverageCredit = avgCredit
				}
			},
		},
		{
			"FrequencyDebit",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "count")
			},
			func(s *OverallSummary, v interface{}) {
				freqDebit, ok := v.(float64)
				if ok {
					s.AllTime.FrequencyDebit = freqDebit
				}
			},
		},
		{
			"FrequencyCredit",
			func() (interface{}, error) {
				return s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "count")
			},
			func(s *OverallSummary, v interface{}) {
				freqCredit, ok := v.(float64)
				if ok {
					s.AllTime.FrequencyCredit = freqCredit
				}
			},
		},
		{
			"TopDebits",
			func() (interface{}, error) {
				return s.repo.GetTopChangeByTransactionType(mutasiID, "debit", 10)
			},
			func(s *OverallSummary, v interface{}) {
				topDebits, ok := v.([]models.MutasiDetail)
				if ok {
					s.AllTime.TopDebits = topDebits
				}
			},
		},
		{
			"TopCredits",
			func() (interface{}, error) {
				return s.repo.GetTopChangeByTransactionType(mutasiID, "credit", 10)
			},
			func(s *OverallSummary, v interface{}) {
				topCredits, ok := v.([]models.MutasiDetail)
				if ok {
					s.AllTime.TopCredits = topCredits
				}
			},
		},
		{
			"AnomalyTransactions",
			func() (interface{}, error) {
				return s.repo.GetAnomalyTransactions(mutasiID)
			},
			func(s *OverallSummary, v interface{}) {
				anomalyTransactions, ok := v.([]models.MutasiDetail)
				if ok {
					s.AllTime.AnomalyTransactions = anomalyTransactions
				}
			},
		},
		{
			"TotalBankFee",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "bank_fee")
			},
			func(s *OverallSummary, v interface{}) {
				totalBankFee, ok := v.(float64)
				if ok {
					s.AllTime.TotalBankFee = totalBankFee
				}
			},
		},
		{
			"TotalInterest",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "interest")
			},
			func(s *OverallSummary, v interface{}) {
				totalInterest, ok := v.(float64)
				if ok {
					s.AllTime.TotalInterest = totalInterest
				}
			},
		},
		{
			"TotalTax",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "tax")
			},
			func(s *OverallSummary, v interface{}) {
				totalTax, ok := v.(float64)
				if ok {
					s.AllTime.TotalTax = totalTax
				}
			},
		},
		{
			"TotalDigitalRevenue",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "digital_revenue")
			},
			func(s *OverallSummary, v interface{}) {
				totalDigitalRevenue, ok := v.(float64)
				if ok {
					s.AllTime.TotalDigitalRevenue = totalDigitalRevenue
				}
			},
		},
		{
			"TotalTransferIn",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "transfer_in")
			},
			func(s *OverallSummary, v interface{}) {
				totalTransferIn, ok := v.(float64)
				if ok {
					s.AllTime.TotalTransferIn = totalTransferIn
				}
			},
		},
		{
			"TotalTransferOut",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "transfer_out")
			},
			func(s *OverallSummary, v interface{}) {
				totalTransferOut, ok := v.(float64)
				if ok {
					s.AllTime.TotalTransferOut = totalTransferOut
				}
			},
		},
		{
			"TotalCashWithdrawal",
			func() (interface{}, error) {
				return s.repo.GetTotalChangeByCategory(mutasiID, "cash_withdrawal")
			},
			func(s *OverallSummary, v interface{}) {
				totalCashWithdrawal, ok := v.(float64)
				if ok {
					s.AllTime.TotalCashWithdrawal = totalCashWithdrawal
				}
			},
		},

		// Monthly fields
		{
			"MonthlyStartBalance",
			func() (interface{}, error) { return s.repo.GetMonthlyBalances(mutasiID, "start") },
			func(s *OverallSummary, v interface{}) {
				monthlyStartBalance, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.StartBalance = monthlyStartBalance
				}
			},
		},
		{
			"MonthlyAverageBalance",
			func() (interface{}, error) { return s.repo.GetMonthlyBalances(mutasiID, "avg") },
			func(s *OverallSummary, v interface{}) {
				monthlyAvgBalance, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.AverageBalance = monthlyAvgBalance
				}
			},
		},
		{
			"MonthlyEndBalance",
			func() (interface{}, error) { return s.repo.GetMonthlyBalances(mutasiID, "end") },
			func(s *OverallSummary, v interface{}) {
				monthlyEndBalance, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.EndBalance = monthlyEndBalance
				}
			},
		},
		{
			"MonthlyTotalIncome",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(
					mutasiID,
					"credit",
					"total",
				)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyTotalIncome, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.TotalIncome = monthlyTotalIncome
				}
			},
		},
		{
			"MonthlyTotalExpenses",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(
					mutasiID,
					"debit",
					"total",
				)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyTotalExpenses, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.TotalExpense = monthlyTotalExpenses
				}
			},
		},
		{
			"MonthlyAverageDebit",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(mutasiID, "debit", "avg")
			},
			func(s *OverallSummary, v interface{}) {
				monthlyAvgDebit, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.AverageDebit = monthlyAvgDebit
				}
			},
		},
		{
			"MonthlyAverageCredit",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(mutasiID, "credit", "avg")
			},
			func(s *OverallSummary, v interface{}) {
				monthlyAvgCredit, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.AverageCredit = monthlyAvgCredit
				}
			},
		},
		{
			"MonthlyFrequencyDebit",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(
					mutasiID,
					"debit",
					"count",
				)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyFreqDebit, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.FrequencyDebit = monthlyFreqDebit
				}
			},
		},
		{
			"MonthlyFrequencyCredit",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTransactionStatsByTransactionType(
					mutasiID,
					"credit",
					"count",
				)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyFreqCredit, ok := v.([]MonthlyAmount)
				if ok {
					s.Monthly.FrequencyCredit = monthlyFreqCredit
				}
			},
		},
		{
			"MonthlyTopDebits",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTopChangeByTransactionType(mutasiID, "debit", 10)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyTopDebits, ok := v.([]MonthlyMutasiDetails)
				if ok {
					s.Monthly.TopDebits = monthlyTopDebits
				}
			},
		},
		{
			"MonthlyTopCredits",
			func() (interface{}, error) {
				return s.repo.GetMonthlyTopChangeByTransactionType(mutasiID, "credit", 10)
			},
			func(s *OverallSummary, v interface{}) {
				monthlyTopCredits, ok := v.([]MonthlyMutasiDetails)
				if ok {
					s.Monthly.TopCredits = monthlyTopCredits
				}
			},
		},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(fields))

	for _, field := range fields {
		if skipFetch(field.name, names) {
			continue
		}

		wg.Add(1)
		go func(f summaryField) {
			defer wg.Done()
			value, err := f.fetchFn()
			if err != nil {
				errCh <- err
				return
			}
			f.assignFn(summary, value)
		}(field)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return summary, nil
}
