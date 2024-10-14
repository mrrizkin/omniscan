package handlers

import (
	"fmt"
	"math"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/mrrizkin/omniscan/app/domains/ocr"
	"github.com/mrrizkin/omniscan/app/utils"
	"github.com/mrrizkin/omniscan/system/stypes"
)

func (h *Handlers) MutasiFindAll(c *fiber.Ctx) error {
	pagination := h.GetPaginationQuery(c)
	mutasis, err := h.ocrService.FindAll(pagination)
	if err != nil {
		h.System.Logger.Error(err, "failed get mutasis")
		return &fiber.Error{
			Code:    500,
			Message: "failed get mutasis",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success get mutasis",
		Data:    mutasis.Result,
		Meta: &stypes.PaginationMeta{
			Page:      pagination.Page,
			PerPage:   pagination.PerPage,
			Total:     mutasis.Total,
			PageCount: int(math.Ceil(float64(mutasis.Total) / float64(pagination.PerPage))),
		},
	})
}

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

	filePath := "storage/" + utils.RandomStr(10) + filePayload.Filename
	err = c.SaveFile(filePayload, filePath)
	if err != nil {
		return &fiber.Error{
			Code:    500,
			Message: "Failed to save the pdf",
		}
	}

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			h.System.Logger.Error(
				err,
				fmt.Sprintf("failed to remove the uploaded file: %s", filePath),
			)
		}
	}()

	scanResult, err := h.ocrService.ScanMutasi(
		payload,
		filePayload,
		filePath,
	)
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
