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
	filePath string,
) (*ScanEStatementResponse, error) {
	if fileHeader == nil {
		return nil, errors.New("file can't be nil")
	}

	file, err := s.convertMultipartFileToBytes(fileHeader)
	if err != nil {
		s.log.Error("failed to convert multipart file to bytes", "err", err)
		return nil, err
	}

	metadata, err := estatementscanner.ExtractEStatementMetadataFile(filePath)
	if err != nil {
		s.log.Error("failed to extract e-statement metadata", "err", err)
		return nil, err
	}

	summaryField := make([]string, 0)
	if payload.Summary != "" {
		fields := strings.Split(payload.Summary, ",")
		for _, field := range fields {
			summaryField = append(summaryField, strings.TrimSpace(field))
		}
	}

	if s.repo.IsFileAlreadyScanned(fileHeader.Filename) {
		return s.getExistingEStatementResponse(fileHeader, file, metadata, summaryField)
	}

	return s.processNewEStatement(payload, fileHeader, file, metadata, summaryField)
}

func (s *EStatementService) GetSummary(eStatementID uint) (*OverallSummary, error) {
	return s.getSummary(eStatementID)
}

func (s *EStatementService) getExistingEStatementResponse(
	fileHeader *multipart.FileHeader,
	file []byte,
	metadata *types.PDFMetadata,
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

	scanResult := types.ScanResult{
		Transactions: transactions,
		Info: types.ScanInfo{
			Bank:     eStatement.Bank,
			Produk:   eStatement.Produk,
			Rekening: eStatement.Rekening,
			Periode:  eStatement.Periode,
		},
	}

	return s.createScanEStatementResponse(
		eStatement.ID,
		&scanResult,
		fileHeader,
		file,
		metadata,
		summary,
	), nil
}

func (s *EStatementService) processNewEStatement(
	payload *ScanEStatementPayload,
	fileHeader *multipart.FileHeader,
	file []byte,
	metadata *types.PDFMetadata,
	summaryField []string,
) (*ScanEStatementResponse, error) {
	scanResult, err := s.scanner.Scan(payload.Provider, file)
	if err != nil {
		s.log.Error("failed to scan e-statement when scanning", "err", err)
		return nil, err
	}

	expiredEStatement, err := calculateEStatementExpiry(payload.TimeBomb)
	if err != nil {
		s.log.Error("failed to calculate e-statement expiry", "err", err)
		return nil, err
	}

	eStatement := s.createEStatementModel(scanResult, fileHeader.Filename, expiredEStatement)

	if err := s.repo.Aggregate(eStatement); err != nil {
		s.log.Error("failed to aggregate e-statement", "err", err)
		return nil, err
	}

	eStatementDetail := s.createEStatementDetailModels(scanResult.Transactions, eStatement.ID)

	if err := s.repo.AggregateDetail(eStatementDetail); err != nil {
		s.log.Error("failed to aggregate e-statement detail", "err", err)
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
		fileHeader,
		file,
		metadata,
		summary,
	), nil
}
