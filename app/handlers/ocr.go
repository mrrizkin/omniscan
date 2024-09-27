package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mrrizkin/omniscan/app/domains/ocr"
	"github.com/mrrizkin/omniscan/system/stypes"
)

func (h *Handlers) ScanMutasi(c *fiber.Ctx) error {
	payload := new(ocr.ScanMutasiPayload)
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

	scanResult, err := h.ocrService.ScanMutasi(payload, filePayload)
	if err != nil {
		h.System.Logger.Error(err, "failed to scan mutasi")
		return &fiber.Error{
			Code:    400,
			Message: "failed to scan mutasi",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success scan mutasi",
		Data:    scanResult,
	})
}

func (h *Handlers) GetSumary(c *fiber.Ctx) error {
	mutasiID, err := c.ParamsInt("id")
	if err != nil {
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	summary, err := h.ocrService.GetSummary(uint(mutasiID))
	if err != nil {
		h.System.Logger.Error(err, "failed to get summary")
		return &fiber.Error{
			Code:    400,
			Message: "failed to get summary",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success scan mutasi",
		Data:    summary,
	})
}
