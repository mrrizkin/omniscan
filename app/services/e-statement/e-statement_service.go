package estatement

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/repositories"
	estatementscanner "github.com/mrrizkin/omniscan/pkg/e-statement-scanner"

	"github.com/mrrizkin/omniscan/pkg/e-statement-scanner/types"
)

type EStatementService struct {
	log *logger.Logger

	repo    *repositories.EStatementRepository
	scanner *estatementscanner.EStatementScanner
}

func (*EStatementService) Construct() interface{} {
	return func(
		log *logger.Logger,

		repo *repositories.EStatementRepository,
		scanner *estatementscanner.EStatementScanner,
	) *EStatementService {
		return &EStatementService{log, repo, scanner}
	}
}

func (s *EStatementService) FindAll(page, perPage int) (*PaginatedEStatement, error) {
	eStatements, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	eStatementsCount, err := s.repo.FindAllCount()
	if err != nil {
		return nil, err
	}

	return &PaginatedEStatement{
		Result: eStatements,
		Total:  int(eStatementsCount),
	}, nil
}

func (s *EStatementService) ScanEStatement(
	payload *ScanEStatementPayload,
	fileHeader *multipart.FileHeader,
) (*ScanEStatementResponse, error) {
	if fileHeader == nil {
		return nil, errors.New("file can't be nil")
	}

	file, err := s.convertMultipartFileToBytes(fileHeader)
	if err != nil {
		s.log.Error("failed to convert multipart file to bytes", "err", err)
		return nil, err
	}

	summaryField := make([]string, 0)
	if payload.Summary != "" {
		fields := strings.Split(payload.Summary, ",")
		for _, field := range fields {
			summaryField = append(summaryField, strings.TrimSpace(field))
		}
	}

	if !payload.IsScanOnly() {
		if s.repo.IsFileAlreadyScanned(fileHeader.Filename) {
			return s.getExistingEStatementResponse(fileHeader, summaryField)
		}
	}

	return s.processNewEStatement(payload, fileHeader, file, summaryField)
}

func (s *EStatementService) GetSummary(eStatementID uint) (*OverallSummary, error) {
	return s.getSummary(eStatementID)
}

func (s *EStatementService) getExistingEStatementResponse(
	fileHeader *multipart.FileHeader,
	summaryField []string,
) (*ScanEStatementResponse, error) {
	eStatement, err := s.repo.GetEStatementByFilename(fileHeader.Filename)
	if err != nil {
		s.log.Error("failed to get e-statement by filename", "err", err)
		return nil, err
	}

	summary, err := s.getSummary(eStatement.ID, summaryField...)
	if err != nil {
		s.log.Error("failed to get summary", "err", err)
		return nil, err
	}

	transactions := s.convertEStatementDetailToTransactions(eStatement.EStatementDetail)
	metadata, err := eStatement.EStatementMetadata.ToMetadata()
	if err != nil {
		s.log.Error("failed to convert e-statement metadata to metadata", "err", err)
		return nil, err
	}

	scanResult := types.ScanResult{
		Transactions: transactions,
		Info: types.ScanInfo{
			Bank:     eStatement.Bank,
			Produk:   eStatement.Produk,
			Rekening: eStatement.Rekening,
			Periode:  eStatement.Periode,
		},
		Metadata: metadata,
	}

	return s.createScanEStatementResponse(
		eStatement.ID,
		&scanResult,
		summary,
	), nil
}

func (s *EStatementService) processNewEStatement(
	payload *ScanEStatementPayload,
	fileHeader *multipart.FileHeader,
	file []byte,
	summaryField []string,
) (*ScanEStatementResponse, error) {
	scanResult, err := s.scanner.Scan(payload.Bank, payload.PDFLib, fileHeader.Filename, file)
	if err != nil {
		s.log.Error("failed to scan e-statement when scanning", "err", err)
		return nil, err
	}

	if scanResult == nil {
		s.log.Error("scan result is nil")
		return nil, errors.New("scan result is nil")
	}

	if len(scanResult.Transactions) == 0 {
		s.log.Error("transactions is empty")
		return nil, errors.New("transactions is empty")
	}

	if payload.IsScanOnly() {
		return s.createScanEStatementResponse(0, scanResult, &OverallSummary{}), nil
	}

	expiredEStatement, err := calculateEStatementExpiry(payload.TimeBomb)
	if err != nil {
		s.log.Error("failed to calculate e-statement expiry", "err", err)
		return nil, err
	}

	eStatement := s.createEStatementModel(scanResult, fileHeader.Filename, expiredEStatement)

	transaction := s.repo.Begin()
	err = transaction.Error
	if err != nil {
		transaction.Rollback()
		s.log.Error("failed to begin transaction", "err", err)
		return nil, err
	}

	if err := s.repo.Aggregate(eStatement, transaction); err != nil {
		transaction.Rollback()
		s.log.Error("failed to aggregate e-statement", "err", err)
		return nil, err
	}

	eStatementMetadata, err := s.createEStatementMetadataModel(scanResult, eStatement.ID)
	if err != nil {
		transaction.Rollback()
		s.log.Error("failed to create e-statement metadata model", "err", err)
		return nil, err
	}

	if err := s.repo.AggregateMetadata(eStatementMetadata, transaction); err != nil {
		transaction.Rollback()
		s.log.Error("failed to aggregate e-statement metadata", "err", err)
		return nil, err
	}

	eStatementDetail := s.createEStatementDetailModels(scanResult.Transactions, eStatement.ID)
	if err := s.repo.AggregateDetail(eStatementDetail, transaction); err != nil {
		transaction.Rollback()
		s.log.Error("failed to aggregate e-statement detail", "err", err)
		return nil, err
	}

	if err := transaction.Commit().Error; err != nil {
		s.log.Error("failed to commit transaction", "err", err)
		return nil, err
	}

	summary, err := s.getSummary(eStatement.ID, summaryField...)
	if err != nil {
		s.log.Error("failed to get summary", "err", err)
		return nil, err
	}

	return s.createScanEStatementResponse(
		eStatement.ID,
		scanResult,
		summary,
	), nil
}
