package ocr

import (
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/database"
	mutasi_scanner "github.com/mrrizkin/omniscan/third-party/mutasi-scanner"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

type Repo struct {
	db *database.Database
}

type Service struct {
	repo    *Repo
	scanner *mutasi_scanner.MutasiScanner
}

type PaginatedMutasi struct {
	Result []models.Mutasi
	Total  int
}

type ScanMutasiPayload struct {
	Provider string `form:"provider"  validate:"required"`
	TimeBomb string `form:"time_bomb"`
}

type ScanMutasiResponse struct {
	*types.ScanResult

	MutasiID uint           `json:"mutasi_id"`
	Meta     Meta           `json:"meta"`
	Summary  OverallSummary `json:"summary"`
}

type Summary struct {
	StartBalance   float64 `json:"start_balance"`
	AverageBalance float64 `json:"average_balance"`
	EndBalance     float64 `json:"end_balance"`

	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`

	TopDebits  []models.MutasiDetail `json:"top_debits"`
	TopCredits []models.MutasiDetail `json:"top_credits"`

	AnomalyTransactions []models.MutasiDetail `json:"anomaly_transactions"`

	TotalBankFee        float64 `json:"total_bank_fee"`
	TotalInterest       float64 `json:"total_interest"`
	TotalTax            float64 `json:"total_tax"`
	TotalDigitalRevenue float64 `json:"total_digital_revenue"`
	TotalTransferIn     float64 `json:"total_transfer_in"`
	TotalTransferOut    float64 `json:"total_transfer_out"`
	TotalCashWithdrawal float64 `json:"total_cash_withdrawal"`

	AverageCredit float64 `json:"average_credit"`
	AverageDebit  float64 `json:"average_debit"`

	FrequencyDebit  float64 `json:"frequency_debit"`
	FrequencyCredit float64 `json:"frequency_credit"`
}

type MonthlySummary struct {
	StartBalance   []MonthlyAmount `json:"start_balance"`
	AverageBalance []MonthlyAmount `json:"average_balance"`
	EndBalance     []MonthlyAmount `json:"end_balance"`

	TotalIncome  []MonthlyAmount `json:"total_income"`
	TotalExpense []MonthlyAmount `json:"total_expense"`

	TopDebits  []MonthlyMutasiDetails `json:"top_debits"`
	TopCredits []MonthlyMutasiDetails `json:"top_credits"`

	AverageCredit []MonthlyAmount `json:"average_credit"`
	AverageDebit  []MonthlyAmount `json:"average_debit"`

	FrequencyDebit  []MonthlyAmount `json:"frequency_debit"`
	FrequencyCredit []MonthlyAmount `json:"frequency_credit"`
}

type OverallSummary struct {
	AllTime Summary        `json:"all_time"`
	Monthly MonthlySummary `json:"monthly"`
}

type Meta struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileMime string `json:"file_mime"`
}

type MonthlyAmount struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type MonthlyMutasiDetails struct {
	Date   time.Time             `json:"date"`
	Detail []models.MutasiDetail `json:"mutasi_details"`
}
