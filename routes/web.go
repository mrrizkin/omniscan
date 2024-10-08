package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/mrrizkin/omniscan/app/handlers"
	"github.com/mrrizkin/omniscan/system/stypes"
)

func WebRoutes(app *stypes.App, handler *handlers.Handlers) {
	ui := app.Group("/", cors.New())
	ui.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(
			"OmniScan: AI-powered OCR solution for swift, accurate data extraction from diverse documents. Simplify your document processing with intelligent recognition technology.",
		)
	})
}
