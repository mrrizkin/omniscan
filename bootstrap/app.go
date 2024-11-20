package bootstrap

import (
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/omniscan/app/console"
	"github.com/mrrizkin/omniscan/app/controllers"
	"github.com/mrrizkin/omniscan/app/middleware"
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/providers"
	"github.com/mrrizkin/omniscan/app/providers/app"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/providers/scheduler"
	"github.com/mrrizkin/omniscan/app/repositories"
	"github.com/mrrizkin/omniscan/app/services"
	"github.com/mrrizkin/omniscan/config"
	estatementscanner "github.com/mrrizkin/omniscan/pkg/e-statement-scanner"
	"github.com/mrrizkin/omniscan/routes"
)

func App() *fx.App {
	return fx.New(
		config.New(),
		controllers.New(),
		middleware.New(),
		models.New(),
		providers.New(),
		repositories.New(),
		services.New(),

		// deps | pkg
		fx.Provide(estatementscanner.New),

		fx.Invoke(
			app.Boot,
			console.Schedule,
			models.AutoMigrate,
			routes.ApiRoutes,
			routes.WebRoutes,
			serveHTTP,
			startScheduler,
		),

		fx.WithLogger(useLogger),
	)
}

func serveHTTP(app *app.App, cfg *config.App, log *logger.Logger) error {
	log.Info("starting server", "port", cfg.PORT)
	return app.Listen(fmt.Sprintf(":%d", cfg.PORT))
}

func startScheduler(scheduler *scheduler.Scheduler) {
	scheduler.Start()
}

func useLogger(logger *logger.Logger) fxevent.Logger {
	return logger
}
