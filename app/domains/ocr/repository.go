package ocr

import (
	"fmt"
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/utils"
	"github.com/mrrizkin/omniscan/system/database"
)

func NewRepo(db *database.Database) *Repo {
	return &Repo{db}
}

func (r *Repo) Aggregate(mutasi *models.Mutasi) error {
	return r.db.Create(mutasi).Error
}

func (r *Repo) AggregateDetail(mutasiDetail []models.MutasiDetail) error {
	return r.db.Create(mutasiDetail).Error
}

func (r *Repo) GetHeader(idMutasi uint) (*models.Mutasi, error) {
	mutasi := new(models.Mutasi)
	err := r.db.Where("id = ?", idMutasi).
		First(mutasi).
		Error
	if err != nil {
		return nil, err
	}

	return mutasi, nil
}

func (r *Repo) GetDetail(idMutasi uint) ([]models.MutasiDetail, error) {
	mutasi := make([]models.MutasiDetail, 0)
	err := r.db.Where("mutasi_id = ?", idMutasi).
		Order("date ASC").
		Find(mutasi).
		Error
	if err != nil {
		return nil, err
	}

	return mutasi, nil
}

func (r *Repo) GetBalance(mutasiID uint, balanceType string) (float64, error) {
	var err error
	var result struct {
		Balance float64
	}

	gormDB := r.db.Model(&models.MutasiDetail{}).
		Where("mutasi_id = ?", mutasiID)

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

func (r *Repo) GetTransactionStatsByTransactionType(
	mutasiID uint,
	transactionType, statType string,
) (float64, error) {
	var err error
	var result struct {
		Total float64
	}
	gormDB := r.db.Model(&models.MutasiDetail{}).
		Where("mutasi_id = ? AND transaction_type = ?", mutasiID, transactionType)

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

func (r *Repo) GetTopChangeByTransactionType(
	mutasiID uint,
	transactionType string,
	limit int,
) ([]models.MutasiDetail, error) {
	var mutasiDetail []models.MutasiDetail
	err := r.db.Where("mutasi_id = ? AND transaction_type = ?", mutasiID, transactionType).
		Order("change DESC").
		Limit(limit).
		Find(&mutasiDetail).Error
	return mutasiDetail, err
}

func (r *Repo) GetMonthlyBalances(mutasiID uint, balanceType string) ([]MonthlyAmount, error) {
	var results []MonthlyAmount
	var sql string

	switch balanceType {
	case "start":
		sql = `SELECT DISTINCT ON (DATE_TRUNC('month', date))
    date,
    balance as amount
FROM
    mutasi_details
WHERE mutasi_id = ?
ORDER BY
    DATE_TRUNC('month', date),
    date ASC`
	case "end":
		sql = `SELECT DISTINCT ON (DATE_TRUNC('month', date))
    date,
    balance as amount
FROM
    mutasi_details
WHERE mutasi_id = ?
ORDER BY
    DATE_TRUNC('month', date),
    date DESC`
	case "avg":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    AVG(balance) AS amount
FROM
    mutasi_details
WHERE mutasi_id = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	default:
		return nil, fmt.Errorf("balanceType: %s is not supportetd", balanceType)
	}

	err := r.db.Raw(sql, mutasiID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly balances: %w", err)
	}

	return results, nil
}

func (r *Repo) GetMonthlyTransactionStatsByTransactionType(
	mutasiID uint,
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
    mutasi_details
WHERE mutasi_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	case "avg":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    AVG(balance) AS amount
FROM
    mutasi_details
WHERE mutasi_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	case "count":
		sql = `SELECT
    DATE_TRUNC('month', date) AS date,
    COUNT(*) AS amount
FROM
    mutasi_details
WHERE mutasi_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`
	default:
		return nil, fmt.Errorf("statType: %s is not supportetd", statType)
	}

	err := r.db.Raw(sql, mutasiID, transactionType).Scan(&results).Error
	return results, err
}

func (r *Repo) GetMonthlyTopChangeByTransactionType(
	mutasiID uint,
	transactionType string,
	limit int,
) ([]MonthlyMutasiDetails, error) {
	var months []struct{ Date time.Time }

	err := r.db.Raw(`SELECT
    DATE_TRUNC('month', date) AS date
FROM
    mutasi_details
WHERE mutasi_id = ? AND transaction_type = ?
GROUP BY
    DATE_TRUNC('month', date)
ORDER BY
    date`, mutasiID, transactionType).
		Scan(&months).Error
	if err != nil {
		return nil, err
	}

	mutasiDetail := make([]MonthlyMutasiDetails, 0)
	for _, month := range months {
		start, end := getMonthStartAndEnd(month.Date)

		var detail []models.MutasiDetail

		err := r.db.Where(`mutasi_id = ?
AND transaction_type = ?
AND date >= ?
AND date <= ?
      `, mutasiID, transactionType, start, end).
			Order("change DESC").
			Limit(limit).
			Find(&detail).Error

		if err != nil {
			return nil, err
		}

		mutasiDetail = append(mutasiDetail, MonthlyMutasiDetails{
			Date:   start,
			Detail: detail,
		})
	}

	return mutasiDetail, nil
}

func (r *Repo) Bomb() error {
	var expiredMutasis []models.Mutasi
	err := r.db.Where("expired IS NOT NULL AND expired < ?", time.Now()).
		Order("date ASC").
		Find(&expiredMutasis).
		Error
	if err != nil {
		return err
	}

	if len(expiredMutasis) == 0 {
		return nil
	}

	ids := make([]uint, len(expiredMutasis))
	for i, expiredMutasi := range expiredMutasis {
		ids[i] = expiredMutasi.ID
	}

	chunkIds := utils.Chunk(ids, 1000)
	for _, chunk := range chunkIds {
		if err := r.db.Unscoped().
			Where("mutasi_id IN ?", chunk).
			Delete(&models.MutasiDetail{}).
			Error; err != nil {
			return err
		}

		if err := r.db.Unscoped().
			Where("id IN ?", chunk).
			Delete(&models.Mutasi{}).
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
