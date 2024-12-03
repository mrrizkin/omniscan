package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/omniscan/app/controllers"
	"github.com/mrrizkin/omniscan/app/controllers/scan"
	"github.com/mrrizkin/omniscan/app/providers/app"
)

func WebRoutes(
	app *app.App,

	dashboardController *controllers.DashboardController,
	eStatementController *scan.EStatementController,
) {
	router := app.WebRoutes()
	router.Get("/", dashboardController.Index)

	scan := router.Group("/scan")

	eStatementRoute := scan.Group("/e-statement")
	eStatementRoute.Get("/", eStatementController.Index)

	router.All("*", func(c *fiber.Ctx) error {
		return app.Render(c, "pages/404", nil)
	})
}
