package ocr

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	mutasi_scanner "github.com/mrrizkin/omniscan/third-party/mutasi-scanner"
)

func NewService(repo *Repo, scanner *mutasi_scanner.MutasiScanner) *Service {
	return &Service{
		scanner: scanner,
		repo:    repo,
	}
}

func (s *Service) ScanMutasi(
	payload *ScanMutasiPayload,
	fileHeader *multipart.FileHeader,
) (*ScanMutasiResponse, error) {
	if fileHeader == nil {
		return nil, errors.New("file can't be nil")
	}

	file, err := s.convertMultipartFileToBytes(fileHeader)
	if err != nil {
		return nil, err
	}

	scanResult, err := s.scanner.Scan(payload.Provider, file)
	if err != nil {
		return nil, err
	}

	var mutasiExpired time.Time
	if payload.TimeBomb == "" {
		mutasiExpired = time.Now().Add(24 * time.Hour)
	} else {
		mutasiExpired, err = time.Parse("2006-01-02 15:04:05", payload.TimeBomb)
		if err != nil {
			return nil, err
		}
	}

	mutasi := &models.Mutasi{
		Bank:     scanResult.Info.Bank,
		Produk:   scanResult.Info.Produk,
		Rekening: scanResult.Info.Rekening,
		Periode:  scanResult.Info.Periode,
		Expired:  &mutasiExpired,
	}

	err = s.repo.Aggregate(mutasi)
	if err != nil {
		return nil, err
	}

	mutasiDetail := make([]models.MutasiDetail, len(scanResult.Transactions))
	for i, detail := range scanResult.Transactions {
		mutasiDetail[i].Date = detail.Date
		mutasiDetail[i].MutasiID = mutasi.ID
		mutasiDetail[i].Description1 = detail.Description1
		mutasiDetail[i].Description2 = detail.Description2
		mutasiDetail[i].Branch = detail.Branch
		mutasiDetail[i].Change = detail.Change
		mutasiDetail[i].TransactionType = detail.TransactionType
		mutasiDetail[i].Balance = detail.Balance
	}

	err = s.repo.AggregateDetail(mutasiDetail)
	if err != nil {
		return nil, err
	}

	summary, err := s.getSummary(mutasi.ID)
	if err != nil {
		return nil, err
	}

	response := ScanMutasiResponse{
		MutasiID:   mutasi.ID,
		ScanResult: scanResult,
		Meta: Meta{
			FileName: fileHeader.Filename,
			FileSize: fileHeader.Size,
			FileMime: http.DetectContentType(file),
		},
		Summary: *summary,
	}

	return &response, nil
}

func (s *Service) GetSummary(mutasiID uint) (*OverallSummary, error) {
	return s.getSummary(mutasiID)
}

func (s *Service) getSummary(mutasiID uint) (*OverallSummary, error) {
	startBalance, err := s.repo.GetBalance(mutasiID, "start")
	if err != nil {
		return nil, err
	}
	avgBalance, err := s.repo.GetBalance(mutasiID, "avg")
	if err != nil {
		return nil, err
	}
	endBalance, err := s.repo.GetBalance(mutasiID, "end")
	if err != nil {
		return nil, err
	}

	totalIncome, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "total")
	if err != nil {
		return nil, err
	}
	totalExpenses, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "total")
	if err != nil {
		return nil, err
	}

	avgDebit, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "avg")
	if err != nil {
		return nil, err
	}
	avgCredit, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "avg")
	if err != nil {
		return nil, err
	}

	freqDebit, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "debit", "count")
	if err != nil {
		return nil, err
	}
	freqCredit, err := s.repo.GetTransactionStatsByTransactionType(mutasiID, "credit", "count")
	if err != nil {
		return nil, err
	}

	top10Debit, err := s.repo.GetTopChangeByTransactionType(mutasiID, "debit", 10)
	if err != nil {
		return nil, err
	}
	top10Credit, err := s.repo.GetTopChangeByTransactionType(mutasiID, "credit", 10)
	if err != nil {
		return nil, err
	}

	startMonthlyBalance, err := s.repo.GetMonthlyBalances(mutasiID, "start")
	if err != nil {
		return nil, err
	}
	avgMonthlyBalance, err := s.repo.GetMonthlyBalances(mutasiID, "avg")
	if err != nil {
		return nil, err
	}
	endMonthlyBalance, err := s.repo.GetMonthlyBalances(mutasiID, "end")
	if err != nil {
		return nil, err
	}

	totalMonthlyIncome, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"credit",
		"total",
	)
	if err != nil {
		return nil, err
	}
	totalMonthlyExpenses, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"debit",
		"total",
	)
	if err != nil {
		return nil, err
	}

	avgMonthlyDebit, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"debit",
		"avg",
	)
	if err != nil {
		return nil, err
	}
	avgMonthlyCredit, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"credit",
		"avg",
	)
	if err != nil {
		return nil, err
	}

	freqMonthlyDebit, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"debit",
		"count",
	)
	if err != nil {
		return nil, err
	}
	freqMonthlyCredit, err := s.repo.GetMonthlyTransactionStatsByTransactionType(
		mutasiID,
		"credit",
		"count",
	)
	if err != nil {
		return nil, err
	}

	topMonthly10Debit, err := s.repo.GetMonthlyTopChangeByTransactionType(mutasiID, "debit", 10)
	if err != nil {
		return nil, err
	}
	topMonthly10Credit, err := s.repo.GetMonthlyTopChangeByTransactionType(mutasiID, "credit", 10)
	if err != nil {
		return nil, err
	}

	return &OverallSummary{
		Monthly: MonthlySummary{
			StartBalance:   startMonthlyBalance,
			AverageBalance: avgMonthlyBalance,
			EndBalance:     endMonthlyBalance,

			TotalIncome:  totalMonthlyIncome,
			TotalExpense: totalMonthlyExpenses,

			AverageDebit:  avgMonthlyDebit,
			AverageCredit: avgMonthlyCredit,

			FrequencyDebit:  freqMonthlyDebit,
			FrequencyCredit: freqMonthlyCredit,

			TopDebits:  topMonthly10Debit,
			TopCredits: topMonthly10Credit,
		},
		AllTime: Summary{
			StartBalance:   startBalance,
			AverageBalance: avgBalance,
			EndBalance:     endBalance,

			TotalIncome:  totalIncome,
			TotalExpense: totalExpenses,

			AverageDebit:  avgDebit,
			AverageCredit: avgCredit,

			FrequencyDebit:  freqDebit,
			FrequencyCredit: freqCredit,

			TopDebits:  top10Debit,
			TopCredits: top10Credit,
		},
	}, nil
}

func (s *Service) convertMultipartFileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
