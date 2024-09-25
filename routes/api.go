package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/mrrizkin/omniscan/app/handlers"
)

func ApiRoutes(api fiber.Router, handler *handlers.Handlers) {
	v1 := api.Group("/v1", cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, finteligo-api-token",
	}))

	v1.Post("/scan-mutasi", handler.ScanMutasi)

	v1.Get("/user", handler.UserFindAll)
	v1.Get("/user/:id", handler.UserFindByID)
	v1.Post("/user", handler.UserCreate)
	v1.Put("/user/:id", handler.UserUpdate)
	v1.Delete("/user/:id", handler.UserDelete)
}
