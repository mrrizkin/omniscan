package ocr

import (
	"io"
	"mime/multipart"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/third-party/mutasi-scanner/types"
)

func (s *Service) convertMutasiDetailToTransactions(
	details []models.MutasiDetail,
) []*types.Transaction {
	transactions := make([]*types.Transaction, len(details))
	for i, detail := range details {
		transactions[i] = &types.Transaction{
			Date:            detail.Date,
			Description1:    detail.Description1,
			Description2:    detail.Description2,
			Branch:          detail.Branch,
			Change:          detail.Change,
			TransactionType: detail.TransactionType,
			Balance:         detail.Balance,
		}
	}
	return transactions
}

func (s *Service) convertMultipartFileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
