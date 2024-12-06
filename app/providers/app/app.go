package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/providers/session"
	"github.com/mrrizkin/omniscan/app/providers/validator"
	"github.com/mrrizkin/omniscan/app/providers/view"
	"github.com/mrrizkin/omniscan/config"
)

// @title						OmniScan API
// @version					    1.0
// @description				    OmniScan API provides a comprehensive set of endpoints for scanning and analyzing e-statement data.
// @termsOfService				https://dak.id/terms/
// @contact.name				API Support Team
// @contact.url				    https://dak.id/support
// @contact.email				support@dak.id
// @license.name				MIT License
// @license.url				    https://opensource.org/licenses/MIT
// @host						localhost:3000
// @BasePath					/api/v1
// @externalDocs.description	OpenAPI Specification
// @externalDocs.url			https://swagger.io/specification/
type App struct {
	*fiber.App

	validator *validator.Validator
	session   *session.Session
	log       *logger.Logger
	view      *view.View
	config    *config.App
}

func (*App) Construct() interface{} {
	return func(
		lc fx.Lifecycle,
		config *config.App,
		session *session.Session,
		validator *validator.Validator,
		log *logger.Logger,
		view *view.View,
	) *App {
		app := createServer(config, log)

		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return app.Shutdown()
			},
		})

		return &App{
			App:       app,
			session:   session,
			validator: validator,
			log:       log,
			view:      view,
			config:    config,
		}
	}
}

// bodyParseValidate parses the request body into the given struct and validates it.
// It returns an error if parsing fails or if the data doesn't pass validation.
func (a *App) ParseBodyAndValidate(c *fiber.Ctx, out interface{}) error {
	err := c.BodyParser(out)
	if err != nil {
		a.log.Error("failed to parse payload", "err", err)
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	err = a.validator.MustValidate(out)
	if err != nil {
		a.log.Error("validation error", "err", err)
		return err
	}

	return nil
}

func (a *App) ApiRoutes() fiber.Router {
	a.App.Get("/api/v1/docs/swagger", swagger(a.config))
	return a.App.Group("/api", cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, boot-api-token",
	}))
}

func (a *App) WebRoutes() fiber.Router {
	return a.App.Group("/",
		csrf.New(csrf.Config{
			KeyLookup:         fmt.Sprintf("cookie:%s", a.config.CSRF_KEY),
			CookieName:        a.config.CSRF_COOKIE_NAME,
			CookieSameSite:    a.config.CSRF_SAME_SITE,
			CookieSecure:      a.config.CSRF_SECURE,
			CookieSessionOnly: true,
			CookieHTTPOnly:    a.config.CSRF_HTTP_ONLY,
			SingleUseToken:    true,
			Expiration:        time.Duration(a.config.CSRF_EXPIRATION) * time.Second,
			KeyGenerator:      utils.UUIDv4,
			ErrorHandler:      csrf.ConfigDefault.ErrorHandler,
			Extractor:         csrf.CsrfFromCookie(a.config.CSRF_KEY),
			Session:           a.session.Store,
			SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
		}),
		cors.New(),
		helmet.New(),
	)
}

func (a *App) Render(c *fiber.Ctx, template string, data map[string]interface{}) error {
	html, err := a.view.Render(template, data)
	if err != nil {
		return err
	}

	return c.Type("html").Send(html)
}

func createServer(
	appcfg *config.App,
	log *logger.Logger,
) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:               appcfg.PREFORK,
		AppName:               appcfg.NAME,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(validator.GlobalErrorResponse{
				Status: "error",
				Detail: err.Error(),
			})
		},
	})

	app.Static("/", "public")
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: log.GetLogger().(*zerolog.Logger),
	}))
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
			log.Error(fmt.Sprintf("panic: %v\n", e))
		},
	}))
	app.Use(idempotency.New())

	return app
}

func swagger(config *config.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		if config.IsProduction() {
			return c.Status(fiber.StatusNotFound).Send(nil)
		}

		html := fmt.Sprintf(`<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Swagger API Reference - Scalar</title>
        <link rel="icon" type="image/svg+xml" href="https://docs.scalar.com/favicon.svg">
        <link rel="icon" type="image/png" href="https://docs.scalar.com/favicon.png">
    </head>
    <body>
        <script id="api-reference" data-url="%s"></script>
        <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    </body>
</html>`, config.SWAGGER_PATH)

		return c.Type("html").Send([]byte(html))
	}
}
