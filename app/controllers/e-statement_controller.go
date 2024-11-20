package controllers

import (
	"fmt"
	"math"
	"os"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/mrrizkin/omniscan/app/controllers/types"
	"github.com/mrrizkin/omniscan/app/providers/app"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/repositories"
	estatement "github.com/mrrizkin/omniscan/app/services/e-statement"
)

type EStatementController struct {
	*app.App

	log *logger.Logger

	eStatementRepo    *repositories.EStatementRepository
	eStatementService *estatement.EStatementService
}

func (*EStatementController) Construct() interface{} {
	return func(
		app *app.App,
		log *logger.Logger,

		eStatementRepo *repositories.EStatementRepository,
		eStatementService *estatement.EStatementService,
	) (*EStatementController, error) {
		return &EStatementController{
			App: app,
			log: log,

			eStatementRepo:    eStatementRepo,
			eStatementService: eStatementService,
		}, nil
	}
}

// EStatementFindAll godoc
//
//	@Summary		Get all e-statements
//	@Description	Retrieve a list of all e-statements with pagination
//	@Tags			E-Statements
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"
//	@Param			per_page	query		int																false	"Number of items per page"
//	@Success		200			{object}	types.Response{data=[]models.EStatement,meta=types.PaginationMeta}	"Successfully retrieved e-statements"
//	@Failure		500			{object}	validator.GlobalErrorResponse									"Internal server error"
//	@Router			/e-statements [get]
func (c *EStatementController) EStatementFindAll(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	perPage := ctx.QueryInt("per_page", 10)

	estatements, err := c.eStatementService.FindAll(page, perPage)
	if err != nil {
		c.log.Error("failed get e-statements", "err", err)
		return &fiber.Error{
			Code:    500,
			Message: "failed get e-statements",
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success get e-statements",
		Data:    estatements.Result,
		Meta: &types.PaginationMeta{
			Page:      page,
			PerPage:   perPage,
			Total:     estatements.Total,
			PageCount: int(math.Ceil(float64(estatements.Total) / float64(perPage))),
		},
	})
}

// EStatementScan godoc
//
//	@Summary		Scan an e-statement
//	@Description	Scan an e-statement and return the result
//	@Tags			E-Statements
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file		true	"PDF file to scan"
//	@Success		200		{object}	types.Response{data=models.EStatement}	"Successfully scanned e-statement"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/e-statements [post]
func (c *EStatementController) EStatementScan(ctx *fiber.Ctx) error {
	var payload estatement.ScanEStatementPayload

	err := c.ParseBodyAndValidate(ctx, &payload)
	if err != nil {
		return err
	}

	filePayload, err := ctx.FormFile("file")
	if err != nil {
		c.log.Error("failed to get file", "err", err)
		return &fiber.Error{
			Code:    400,
			Message: "failed to get file",
		}
	}
	if filePayload == nil {
		c.log.Error("failed to get file, file null")
		return &fiber.Error{
			Code:    400,
			Message: "failed to get file, file null",
		}
	}

	filePath := "storage/" + randomStr(10) + filePayload.Filename
	err = ctx.SaveFile(filePayload, filePath)
	if err != nil {
		return &fiber.Error{
			Code:    500,
			Message: "Failed to save the pdf",
		}
	}

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			c.log.Error(fmt.Sprintf("failed to remove the uploaded file: %s", filePath), "err", err)
		}
	}()

	scanResult, err := c.eStatementService.ScanEStatement(&payload, filePayload, filePath)
	if err != nil {
		c.log.Error("failed to scan e-statement", "err", err)
		return &fiber.Error{
			Code:    400,
			Message: fmt.Sprintf("failed to scan e-statement: %s", err),
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success scan e-statement",
		Data:    scanResult,
	})
}

// EStatementGetSumary godoc
//
//	@Summary		Get an e-statement summary
//	@Description	Get an e-statement summary
//	@Tags			E-Statements
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"E-Statement ID"
//	@Success		200	{object}	types.Response{data=estatement.OverallSummary}	"Successfully retrieved e-statement summary"
//	@Failure		400	{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500	{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/e-statements/{id}/summary [get]
func (c *EStatementController) EStatementGetSumary(ctx *fiber.Ctx) error {
	eStatementID, err := ctx.ParamsInt("id")
	if err != nil {
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	summary, err := c.eStatementService.GetSummary(uint(eStatementID))
	if err != nil {
		c.log.Error("failed to get summary", "err", err)
		return &fiber.Error{
			Code:    400,
			Message: "failed to get summary",
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success get e-statement summary",
		Data:    summary,
	})
}

func randomStr(length int) string {
	str, _ := gonanoid.Generate(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
		length,
	)
	return str
}
