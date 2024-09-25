package handlers

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"

	"github.com/mrrizkin/omniscan/system/stypes"
)

type ScanMutasiPayload struct {
	Provider string `form:"provider" validate:"required"`
}

func (h *Handlers) ScanMutasi(c *fiber.Ctx) error {
	payload := new(ScanMutasiPayload)
	err := c.BodyParser(payload)
	if err != nil {
		h.System.Logger.Error(err, "failed to parse payload")
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	validationError := h.System.Validator.MustValidate(payload)
	if validationError != nil {
		return validationError
	}

	filePayload, err := c.FormFile("file")
	if err != nil {
		h.System.Logger.Error(err, "failed to get file")
		return &fiber.Error{
			Code:    400,
			Message: "failed to get file",
		}
	}
	if filePayload == nil {
		h.System.Logger.Error(err, "failed to get file, file null")
		return &fiber.Error{
			Code:    400,
			Message: "failed to get file, file null",
		}
	}
	file, err := convertMultipartFileToBytes(filePayload)
	if err != nil {
		h.System.Logger.Error(err, "failed to open file, convert to bytes")
		return &fiber.Error{
			Code:    400,
			Message: "failed to open file",
		}
	}

	trx, err := h.Library.MutasiScanner.Scan(payload.Provider, file)
	if err != nil {
		h.System.Logger.Error(err, "failed to scan mutasi")
		return &fiber.Error{
			Code:    500,
			Message: "failed scan mutasi",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success scan mutasi",
		Data:    trx,
	})
}

func convertMultipartFileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
