package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/omniscan/app/providers/app"
)

type DashboardController struct {
	*app.App
}

func (*DashboardController) Construct() interface{} {
	return func(app *app.App) (*DashboardController, error) {
		return &DashboardController{
			App: app,
		}, nil
	}
}

func (c *DashboardController) Index(ctx *fiber.Ctx) error {
	return c.Render(ctx, "pages/dashboard", nil)
}
