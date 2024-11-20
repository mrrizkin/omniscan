package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/omniscan/app/providers/app"
)

type WelcomeController struct {
	*app.App
}

func (*WelcomeController) Construct() interface{} {
	return func(app *app.App) (*WelcomeController, error) {
		return &WelcomeController{
			App: app,
		}, nil
	}
}

func (c *WelcomeController) Index(ctx *fiber.Ctx) error {
	return c.Render(ctx, "pages/welcome", nil)
}
