package routes

import (
	"github.com/mrrizkin/omniscan/app/controllers"
	"github.com/mrrizkin/omniscan/app/providers/app"
)

func WebRoutes(
	app *app.App,

	welcomeController *controllers.WelcomeController,
) {
	router := app.WebRoutes()
	router.Get("/", welcomeController.Index)
}
