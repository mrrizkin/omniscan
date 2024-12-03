package scan

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/omniscan/app/providers/app"
)

type EStatementController struct {
	*app.App
}

func (*EStatementController) Construct() interface{} {
	return func(app *app.App) (*EStatementController, error) {
		return &EStatementController{
			App: app,
		}, nil
	}
}

func (c *EStatementController) Index(ctx *fiber.Ctx) error {
	return c.Render(ctx, "pages/scan/e-statement/index", nil)
}
