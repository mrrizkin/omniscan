package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/providers/database"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"gorm.io/gorm"
)

type EStatementRepository struct {
	db  *database.Database
	log *logger.Logger
}

func (r *EStatementRepository) Construct() interface{} {
	return func(db *database.Database, log *logger.Logger) *EStatementRepository {
		return &EStatementRepository{db, log}
	}
}

func (r *EStatementRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

func (r *EStatementRepository) FindAll(page, perPage int) ([]models.EStatement, error) {
	eStatement := make([]models.EStatement, 0)
	err := r.db.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&eStatement).Error
	return eStatement, err
}

func (r *EStatementRepository) FindAllCount() (int64, error) {
	var count int64 = 0
	err := r.db.Model(&models.EStatement{}).Count(&count).Error
	return count, err
}

func (r *EStatementRepository) Aggregate(eStatement *models.EStatement, db *gorm.DB) error {
	return db.Create(eStatement).Error
}

func (r *EStatementRepository) AggregateMetadata(eStatementMetadata *models.EStatementMetadata, db *gorm.DB) error {
	return db.Create(eStatementMetadata).Error
}

func (r *EStatementRepository) AggregateDetail(eStatementDetail []models.EStatementDetail, db *gorm.DB) error {
	return db.Create(eStatementDetail).Error
}

func (r *EStatementRepository) IsFileAlreadyScanned(filename string) bool {
	eStatement := new(models.EStatement)
	err := r.db.Preload("EStatementDetail").Where("filename = ?", filename).
		First(eStatement).
		Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *EStatementRepository) GetEStatementByFilename(filename string) (*models.EStatement, error) {
	eStatement := new(models.EStatement)
	err := r.db.Where("filename = ?", filename).
		First(eStatement).
		Error
	if err != nil {
		return nil, err
	}

	return eStatement, nil
}

func (r *EStatementRepository) GetHeader(eStatementID uint) (*models.EStatement, error) {
	eStatement := new(models.EStatement)
	err := r.db.Where("id = ?", eStatementID).
		First(eStatement).
		Error
	if err != nil {
		return nil, err
	}

	return eStatement, nil
}

func (r *EStatementRepository) GetDetail(eStatementID uint) ([]models.EStatementDetail, error) {
	eStatement := make([]models.EStatementDetail, 0)
	err := r.db.Where("e_statement_id = ?", eStatementID).
		Order("date ASC").
		Find(eStatement).
		Error
	if err != nil {
		return nil, err
	}

	return eStatement, nil
}

func (r *EStatementRepository) GetBalance(idEstatement uint, balanceType string) (float64, error) {
	var err error
	var result struct {
		Balance float64
	}

	gormDB := r.db.Model(&models.EStatementDetail{}).
		Where("e_statement_id = ?", idEstatement)

	switch balanceType {
	case "start":
		err = gormDB.Select("balance").
			Order("date ASC").
			Limit(1).
			Scan(&result).Error
	case "end":
		err = gormDB.Select("balance").
			Order("date DESC").
			Limit(1).
			Scan(&result).Error
	case "avg":
		err = gormDB.Select("AVG(balance) as balance").
			Scan(&result).Error
	default:
		return 0, fmt.Errorf("balanceType: %s is not supportetd", balanceType)
	}

	return result.Balance, err
}

func (r *EStatementRepository) GetTransactionStatsByTransactionType(
	eStatementID uint,
	transactionType, statType string,
) (float64, error) {
	var err error
	var result struct {
		Total float64
	}
	gormDB := r.db.Model(&models.EStatementDetail{}).
		Where("e_statement_id = ? AND transaction_type = ?", eStatementID, transactionType)

	switch statType {
	case "total":
		err = gormDB.Select("SUM(change) as total").
			Scan(&result).Error
	case "avg":
		err = gormDB.Select("AVG(change) as total").
			Scan(&result).Error
	case "count":
		err = gormDB.Select("COUNT(*) as total").
			Scan(&result).Error
	default:
		return 0, fmt.Errorf("statType: %s is not supportetd", statType)
	}

	return result.Total, err
}

func (r *EStatementRepository) GetTopChangeByTransactionType(
	eStatementID uint,
	transactionType string,
	limit int,
) ([]models.EStatementDetail, error) {
	var eStatementDetail []models.EStatementDetail
	err := r.db.Where("e_statement_id = ? AND transaction_type = ?", eStatementID, transactionType).
		Order("change DESC").
		Limit(limit).
		Find(&eStatementDetail).Error
	return eStatementDetail, err
}

func (r *EStatementRepository) GetAnomalyTransactions(
	eStatementID uint,
) ([]models.EStatementDetail, error) {
	var eStatementDetail []models.EStatementDetail
	err := r.db.Where("e_statement_id = ? AND change%1000000 = 0 AND description1 != 'SALDO AWAL'", eStatementID).
		Order("change DESC").
		Find(&eStatementDetail).Error
	return eStatementDetail, err
}

func (r *EStatementRepository) GetTotalChangeByCategory(
	eStatementID uint,
	categoryType string,
) (float64, error) {
	var err error
	var result struct {
		Total float64
	}

	wb := whereBuilder()
	wb.And("e_statement_id = ?", eStatementID)
	switch categoryType {
	case "bank_fee":
		wb.And("((description1 = 'SWITCHING DB' AND balance != 0) OR (description1 = 'BI-FAST DB' AND description2 ILIKE '%BIAYA%') OR description1 = 'BIAYA ADM')")
	case "interest":
		wb.And("description1 = 'BUNGA'")
	case "tax":
		wb.And("description1 ILIKE '%PAJAK%'")
	case "digital_revenue":
		wb.And("description2 ILIKE '%ESPAY DEBIT%'")
	case "transfer_in":
		wb.And("transaction_type = 'credit' AND (description1 LIKE '%TRSF%' OR description1 = 'SWITCHING CR')")
		wb.And("description2 NOT LIKE '%ESPAY DEBIT%'")
	case "transfer_out":
		wb.And("transaction_type = 'debit' AND (description1 LIKE '%TRSF%' OR (description1 = 'SWITCHING DB' AND balance = 0) OR (description1 = 'BI-FAST DB' AND description2 ILIKE '%BIF TRANSFER%') OR description1 ILIKE '%BYR%')")
	case "cash_withdrawal":
		wb.And("description1 ILIKE '%TARIKAN ATM%'")
	default:
		return 0, fmt.Errorf("categoryType: %s is not supportetd", categoryType)
	}

	where, whereArgs := wb.Get()
	err = r.db.Where(where, whereArgs...).
		Model(models.EStatementDetail{}).
		Select("SUM(change) as total").Scan(&result).Error
	return result.Total, err
}

type MonthlyAmount struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type MonthlyEStatementDetails struct {
	Date   time.Time                 `json:"date"`
	Detail []models.EStatementDetail `json:"e_statement_details"`
}

func (r *EStatementRepository) GetMonthlyBalances(eStatementID uint, balanceType string) ([]MonthlyAmount, error) {
	var results []MonthlyAmount
	var sql string

	switch balanceType {
	case "start":
		sql = `SELECT DISTINCT ON (DATE_TRUNC('month', date))
    date,
    balance as amount
FROM
    e_statement_details
WHERE e_statement_id = ?
ORDER BY
    DATE_TRUNC('month', date),
    date ASC`
	case "end":
		sql = `SELECT DISTINCT ON (DATE_TRUNC('month', date))
    date,
    balance as amount
FROM
    e_statement_details
WHERE e_statement_id = ?
ORDER BY
    DATE_TRUNC('month', date),
    date DESC`
	case "avg":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    AVG(balance) AS amount
FROM
    e_statement_details
WHERE e_statement_id = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	default:
		return nil, fmt.Errorf("balanceType: %s is not supportetd", balanceType)
	}

	err := r.db.Raw(sql, eStatementID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly balances: %w", err)
	}

	return results, nil
}

func (r *EStatementRepository) GetMonthlyTransactionStatsByTransactionType(
	eStatementID uint,
	transactionType, statType string,
) ([]MonthlyAmount, error) {
	var results []MonthlyAmount
	var sql string

	switch statType {
	case "total":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    SUM(balance) AS amount
FROM
    e_statement_details
WHERE e_statement_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	case "avg":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    AVG(balance) AS amount
FROM
    e_statement_details
WHERE e_statement_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	case "count":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    COUNT(*) AS amount
FROM
    e_statement_details
WHERE e_statement_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	default:
		return nil, fmt.Errorf("statType: %s is not supportetd", statType)
	}

	err := r.db.Raw(sql, eStatementID, transactionType).Scan(&results).Error
	return results, err
}

func (r *EStatementRepository) GetMonthlyTopChangeByTransactionType(
	eStatementID uint,
	transactionType string,
	limit int,
) ([]MonthlyEStatementDetails, error) {
	var months []struct{ Date time.Time }

	err := r.db.Raw(`SELECT
    DATE_TRUNC('month', date) AS date
FROM
    e_statement_details
WHERE e_statement_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`, eStatementID, transactionType).
		Scan(&months).Error
	if err != nil {
		return nil, err
	}

	eStatementDetails := make([]MonthlyEStatementDetails, 0)
	for _, month := range months {
		start, end := getMonthStartAndEnd(month.Date)

		var detail []models.EStatementDetail

		err := r.db.Where(`e_statement_id = ?
AND transaction_type = ?
AND date >= ?
AND date <= ?
      `, eStatementID, transactionType, start, end).
			Order("change DESC").
			Limit(limit).
			Find(&detail).Error

		if err != nil {
			return nil, err
		}

		eStatementDetails = append(eStatementDetails, MonthlyEStatementDetails{
			Date:   start,
			Detail: detail,
		})
	}

	return eStatementDetails, nil
}

func (r *EStatementRepository) Bomb() error {
	var expiredEStatements []models.EStatement
	err := r.db.Where("expired IS NOT NULL AND expired < ?", time.Now()).
		Order("date ASC").
		Find(&expiredEStatements).
		Error
	if err != nil {
		return err
	}

	if len(expiredEStatements) == 0 {
		return nil
	}

	ids := make([]uint, len(expiredEStatements))
	for i, expiredEStatement := range expiredEStatements {
		ids[i] = expiredEStatement.ID
	}

	chunkIds := chunk(ids, 1000)
	for _, chunk := range chunkIds {
		log.Info("deleting expired e-statements", "ids", chunk)
		if err := r.db.Unscoped().
			Where("e_statement_id IN ?", chunk).
			Delete(&models.EStatementDetail{}).
			Error; err != nil {
			return err
		}

		if err := r.db.Unscoped().
			Where("id IN ?", chunk).
			Delete(&models.EStatement{}).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func getMonthStartAndEnd(t time.Time) (time.Time, time.Time) {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return start, end
}
