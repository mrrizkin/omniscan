package estatement

import (
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

func (s *EStatementService) createEStatementModel(
	scanResult *types.ScanResult,
	filename string,
	expiry time.Time,
) *models.EStatement {
	return &models.EStatement{
		Bank:     scanResult.Info.Bank,
		Filename: filename,
		Produk:   scanResult.Info.Produk,
		Rekening: scanResult.Info.Rekening,
		Periode:  scanResult.Info.Periode,
		Expired:  &expiry,
	}
}

func (s *EStatementService) createEStatementMetadataModel(
	scanResult *types.ScanResult,
	eStatementID uint,
) (*models.EStatementMetadata, error) {
	return models.IntoEstatementMetadata(eStatementID, scanResult.Metadata)
}

func (s *EStatementService) createEStatementDetailModels(
	transactions []*types.Transaction,
	eStatementID uint,
) []models.EStatementDetail {
	eStatementDetails := make([]models.EStatementDetail, len(transactions))
	for i, detail := range transactions {
		eStatementDetails[i] = models.EStatementDetail{
			Date:            detail.Date,
			EStatementID:    eStatementID,
			Description1:    detail.Description1,
			Description2:    detail.Description2,
			Branch:          detail.Branch,
			Change:          detail.Change,
			TransactionType: detail.TransactionType,
			Balance:         detail.Balance,
		}
	}
	return eStatementDetails
}

func (s *EStatementService) createScanEStatementResponse(
	eStatementID uint,
	scanResult *types.ScanResult,
	summary *OverallSummary,
) *ScanEStatementResponse {
	return &ScanEStatementResponse{
		EStatementID: eStatementID,
		ScanResult:   scanResult,
		Summary:      *summary,
	}
}
