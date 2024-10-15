package ocr

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/mrrizkin/omniscan/system/stypes"
	mutasi_scanner "github.com/mrrizkin/omniscan/third-party/mutasi-scanner"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

func NewService(repo *Repo, scanner *mutasi_scanner.MutasiScanner) *Service {
	return &Service{
		scanner: scanner,
		repo:    repo,
	}
}

func (s *Service) FindAll(pagination stypes.Pagination) (*PaginatedMutasi, error) {
	mutasis, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	mutasisCount, err := s.repo.FindAllCount()
	if err != nil {
		return nil, err
	}

	return &PaginatedMutasi{
		Result: mutasis,
		Total:  int(mutasisCount),
	}, nil
}

func (s *Service) ScanMutasi(
	payload *ScanMutasiPayload,
	fileHeader *multipart.FileHeader,
	filePath string,
) (*ScanMutasiResponse, error) {
	if fileHeader == nil {
		return nil, errors.New("file can't be nil")
	}

	file, err := s.convertMultipartFileToBytes(fileHeader)
	if err != nil {
		return nil, err
	}

	metadata, err := mutasi_scanner.ExtractMutasiMetadataFile(filePath)
	if err != nil {
		return nil, err
	}

	summaryField := make([]string, 0)
	if payload.Summary != "" {
		fmt.Println(payload.Summary)
		fields := strings.Split(payload.Summary, ",")
		for _, field := range fields {
			summaryField = append(summaryField, strings.TrimSpace(field))
		}
	}

	if s.repo.IsFileAlreadyScanned(fileHeader.Filename) {
		return s.getExistingMutasiResponse(fileHeader, file, metadata, summaryField)
	}

	return s.processNewMutasi(payload, fileHeader, file, metadata, summaryField)
}

func (s *Service) GetSummary(mutasiID uint) (*OverallSummary, error) {
	return s.getSummary(mutasiID)
}

func (s *Service) getExistingMutasiResponse(
	fileHeader *multipart.FileHeader,
	file []byte,
	metadata *types.PDFMetadata,
	summaryField []string,
) (*ScanMutasiResponse, error) {
	mutasi, err := s.repo.GetMutasiByFilename(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	summary, err := s.getSummary(mutasi.ID, summaryField...)
	if err != nil {
		return nil, err
	}

	transactions := s.convertMutasiDetailToTransactions(mutasi.MutasiDetail)

	scanResult := types.ScanResult{
		Transactions: transactions,
		Info: types.ScanInfo{
			Bank:     mutasi.Bank,
			Produk:   mutasi.Produk,
			Rekening: mutasi.Rekening,
			Periode:  mutasi.Periode,
		},
	}

	return s.createScanMutasiResponse(
		mutasi.ID,
		&scanResult,
		fileHeader,
		file,
		metadata,
		summary,
	), nil
}

func (s *Service) processNewMutasi(
	payload *ScanMutasiPayload,
	fileHeader *multipart.FileHeader,
	file []byte,
	metadata *types.PDFMetadata,
	summaryField []string,
) (*ScanMutasiResponse, error) {
	scanResult, err := s.scanner.Scan(payload.Provider, file)
	if err != nil {
		return nil, err
	}

	mutasiExpired, err := calculateMutasiExpiry(payload.TimeBomb)
	if err != nil {
		return nil, err
	}

	mutasi := s.createMutasiModel(scanResult, fileHeader.Filename, mutasiExpired)

	if err := s.repo.Aggregate(mutasi); err != nil {
		return nil, err
	}

	mutasiDetail := s.createMutasiDetailModels(scanResult.Transactions, mutasi.ID)

	if err := s.repo.AggregateDetail(mutasiDetail); err != nil {
		return nil, err
	}

	summary, err := s.getSummary(mutasi.ID, summaryField...)
	if err != nil {
		return nil, err
	}

	return s.createScanMutasiResponse(
		mutasi.ID,
		scanResult,
		fileHeader,
		file,
		metadata,
		summary,
	), nil
}
