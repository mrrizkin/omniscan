package ocr

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

func (s *Service) createMutasiModel(
	scanResult *types.ScanResult,
	filename string,
	expiry time.Time,
) *models.Mutasi {
	return &models.Mutasi{
		Bank:     scanResult.Info.Bank,
		Filename: filename,
		Produk:   scanResult.Info.Produk,
		Rekening: scanResult.Info.Rekening,
		Periode:  scanResult.Info.Periode,
		Expired:  &expiry,
	}
}

func (s *Service) createMutasiDetailModels(
	transactions []*types.Transaction,
	mutasiID uint,
) []models.MutasiDetail {
	mutasiDetail := make([]models.MutasiDetail, len(transactions))
	for i, detail := range transactions {
		mutasiDetail[i] = models.MutasiDetail{
			Date:            detail.Date,
			MutasiID:        mutasiID,
			Description1:    detail.Description1,
			Description2:    detail.Description2,
			Branch:          detail.Branch,
			Change:          detail.Change,
			TransactionType: detail.TransactionType,
			Balance:         detail.Balance,
		}
	}
	return mutasiDetail
}

func (s *Service) createScanMutasiResponse(
	mutasiID uint,
	scanResult *types.ScanResult,
	fileHeader *multipart.FileHeader,
	file []byte,
	metadata *types.PDFMetadata,
	summary *OverallSummary,
) *ScanMutasiResponse {
	return &ScanMutasiResponse{
		MutasiID:   mutasiID,
		ScanResult: scanResult,
		Meta: Meta{
			FileName:     fileHeader.Filename,
			FileSize:     fileHeader.Size,
			FileMime:     http.DetectContentType(file),
			Title:        metadata.Title,
			Author:       metadata.Author,
			Subject:      metadata.Subject,
			Keywords:     metadata.Keywords,
			Creator:      metadata.Creator,
			Producer:     metadata.Producer,
			CreationDate: metadata.CreationDate,
			ModDate:      metadata.ModDate,
			PageCount:    metadata.PageCount,
			PDFVersion:   metadata.PDFVersion,
		},
		Summary: *summary,
	}
}
