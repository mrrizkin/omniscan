package estatement

import (
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

type PaginatedEStatement struct {
	Result []models.EStatement
	Total  int
}

type ScanEStatementPayload struct {
	PDFLib   string `form:"pdf_library"  validate:"required"`
	Bank     string `form:"bank"`
	TimeBomb string `form:"time_bomb"`
	Summary  string `form:"summary"`
}

type ScanEStatementResponse struct {
	*types.ScanResult

	EStatementID uint           `json:"e_statement_id"`
	Summary      OverallSummary `json:"summary"`
}

type Summary struct {
	StartBalance   float64 `json:"start_balance,omitempty"`
	AverageBalance float64 `json:"average_balance,omitempty"`
	EndBalance     float64 `json:"end_balance,omitempty"`

	TotalIncome  float64 `json:"total_income,omitempty"`
	TotalExpense float64 `json:"total_expense,omitempty"`

	TopDebits  []models.EStatementDetail `json:"top_debits,omitempty"`
	TopCredits []models.EStatementDetail `json:"top_credits,omitempty"`

	AnomalyTransactions []models.EStatementDetail `json:"anomaly_transactions,omitempty"`

	TotalBankFee        float64 `json:"total_bank_fee,omitempty"`
	TotalInterest       float64 `json:"total_interest,omitempty"`
	TotalTax            float64 `json:"total_tax,omitempty"`
	TotalDigitalRevenue float64 `json:"total_digital_revenue,omitempty"`
	TotalTransferIn     float64 `json:"total_transfer_in,omitempty"`
	TotalTransferOut    float64 `json:"total_transfer_out,omitempty"`
	TotalCashWithdrawal float64 `json:"total_cash_withdrawal,omitempty"`

	AverageCredit float64 `json:"average_credit,omitempty"`
	AverageDebit  float64 `json:"average_debit,omitempty"`

	FrequencyDebit  float64 `json:"frequency_debit,omitempty"`
	FrequencyCredit float64 `json:"frequency_credit,omitempty"`
}

type MonthlySummary struct {
	StartBalance   []MonthlyAmount `json:"start_balance,omitempty"`
	AverageBalance []MonthlyAmount `json:"average_balance,omitempty"`
	EndBalance     []MonthlyAmount `json:"end_balance,omitempty"`

	TotalIncome  []MonthlyAmount `json:"total_income,omitempty"`
	TotalExpense []MonthlyAmount `json:"total_expense,omitempty"`

	TopDebits  []MonthlyEStatementDetails `json:"top_debits,omitempty"`
	TopCredits []MonthlyEStatementDetails `json:"top_credits,omitempty"`

	AverageCredit []MonthlyAmount `json:"average_credit,omitempty"`
	AverageDebit  []MonthlyAmount `json:"average_debit,omitempty"`

	FrequencyDebit  []MonthlyAmount `json:"frequency_debit,omitempty"`
	FrequencyCredit []MonthlyAmount `json:"frequency_credit,omitempty"`
}

type OverallSummary struct {
	AllTime Summary        `json:"all_time,omitempty"`
	Monthly MonthlySummary `json:"monthly,omitempty"`
}

type Meta struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileMime string `json:"file_mime"`

	Title        string `json:"title"`
	Author       string `json:"author"`
	Subject      string `json:"subject"`
	Keywords     string `json:"keywords"`
	Creator      string `json:"creator"`
	Producer     string `json:"producer"`
	CreationDate string `json:"creation_date"`
	ModDate      string `json:"mod_date"`
	PageCount    int    `json:"page_count"`
	PDFVersion   string `json:"pdf_version"`
}

type MonthlyAmount struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type MonthlyEStatementDetails struct {
	Date   time.Time                 `json:"date"`
	Detail []models.EStatementDetail `json:"e_statement_details"`
}
