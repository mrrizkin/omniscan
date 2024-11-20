package controllers

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"

	"github.com/mrrizkin/omniscan/app/controllers/types"
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/providers/app"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/repositories"
	"github.com/mrrizkin/omniscan/app/services"
)

type UserController struct {
	*app.App

	log *logger.Logger

	userService *services.UserService
	userRepo    *repositories.UserRepository
}

func (*UserController) Construct() interface{} {
	return func(
		app *app.App,
		log *logger.Logger,

		userService *services.UserService,
		userRepo *repositories.UserRepository,
	) (*UserController, error) {
		return &UserController{
			App: app,
			log: log,

			userService: userService,
			userRepo:    userRepo,
		}, nil
	}
}

// UserCreate godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided information
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User							true	"User information"
//	@Success		200		{object}	types.Response{data=models.User}	"Successfully created user"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user [post]
func (c *UserController) UserCreate(ctx *fiber.Ctx) error {
	payload := new(models.User)
	err := c.ParseBodyAndValidate(ctx, payload)
	if err != nil {
		return err
	}

	user, err := c.userService.Create(payload)
	if err != nil {
		c.log.Error("failed to create user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to create user: %s", err),
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	})
}

// UserFindAll godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"
//	@Param			per_page	query		int																false	"Number of items per page"
//	@Success		200			{object}	types.Response{data=[]models.User,meta=types.PaginationMeta}	"Successfully retrieved users"
//	@Failure		500			{object}	validator.GlobalErrorResponse									"Internal server error"
//	@Router			/user [get]
func (c *UserController) UserFindAll(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	perPage := ctx.QueryInt("per_page", 10)

	users, err := c.userService.FindAll(page, perPage)
	if err != nil {
		c.log.Error("failed to get users", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get users: %s", err),
		}
	}

	total, ok := users["total"].(int)
	if !ok {
		total = 0
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Message: "users retrieved successfully",
		Data:    users["result"],
		Meta: &types.PaginationMeta{
			Page:      page,
			PerPage:   perPage,
			Total:     total,
			PageCount: int(math.Ceil(float64(total) / float64(perPage))),
		},
	})
}

// UserFindByID godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int									true	"User ID"
//	@Success		200	{object}	types.Response{data=models.User}	"Successfully retrieved user"
//	@Failure		400	{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		404	{object}	validator.GlobalErrorResponse		"User not found"
//	@Failure		500	{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [get]
func (c *UserController) UserFindByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	user, err := c.userService.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "user not found",
			}
		}

		c.log.Error("failed to get user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get user: %s", err),
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UserUpdate godoc
//
//	@Summary		Update a user
//	@Description	Update a user's information by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"User ID"
//	@Param			user	body		models.User							true	"Updated user information"
//	@Success		200		{object}	types.Response{data=models.User}	"Successfully updated user"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [put]
func (c *UserController) UserUpdate(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	payload := new(models.User)
	err = c.ParseBodyAndValidate(ctx, payload)
	if err != nil {
		return err
	}

	user, err := c.userService.Update(uint(id), payload)
	if err != nil {
		c.log.Error("failed to update user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to update user: %s", err),
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	})
}

// UserDelete godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"User ID"
//	@Success		200	{object}	types.Response{}	"Successfully deleted user"
//	@Failure		400	{object}	validator.GlobalErrorResponse	"Bad request"
//	@Failure		401	{object}	validator.GlobalErrorResponse	"Unauthorized"
//	@Failure		500	{object}	validator.GlobalErrorResponse	"Internal server error"
//	@Router			/user/{id} [delete]
func (c *UserController) UserDelete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	user := c.getUser(ctx)
	if user == nil {
		return &fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "unauthorized",
		}
	}

	err = c.userService.Delete(user, uint(id))
	if err != nil {
		c.log.Error("failed to delete user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete user: %s", err),
		}
	}

	return ctx.JSON(types.Response{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

func (c *UserController) getUser(ctx *fiber.Ctx) *models.User {
	userId, ok := ctx.Locals("uid").(uint)
	if !ok {
		return nil
	}

	user, err := c.userRepo.FindByID(userId)
	if err != nil {
		return nil
	}

	return user
}
