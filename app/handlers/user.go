package handlers

import (
	"math"

	"github.com/gofiber/fiber/v2"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/stypes"
)

func (h *Handlers) UserCreate(c *fiber.Ctx) error {
	payload := new(models.User)
	err := c.BodyParser(payload)
	if err != nil {
		h.System.Logger.Error(err, "failed to parse payload")
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	validationError := h.System.Validator.MustValidate(payload)
	if validationError != nil {
		return validationError
	}

	user, err := h.userService.Create(payload)
	if err != nil {
		h.System.Logger.Error(err, "failed create user")
		return &fiber.Error{
			Code:    500,
			Message: "failed create user",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success create user",
		Data:    user,
	})
}

func (h *Handlers) UserFindAll(c *fiber.Ctx) error {
	pagination := h.GetPaginationQuery(c)
	users, err := h.userService.FindAll(pagination)
	if err != nil {
		h.System.Logger.Error(err, "failed get users")
		return &fiber.Error{
			Code:    500,
			Message: "failed get users",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success get users",
		Data:    users.Result,
		Meta: &stypes.PaginationMeta{
			Page:      pagination.Page,
			PerPage:   pagination.PerPage,
			Total:     users.Total,
			PageCount: int(math.Ceil(float64(users.Total) / float64(pagination.PerPage))),
		},
	})
}

func (h *Handlers) UserFindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.System.Logger.Error(err, "failed to parse id")
		return &fiber.Error{
			Code:    400,
			Message: "id not valid",
		}
	}

	user, err := h.userService.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return &fiber.Error{
				Code:    404,
				Message: "user not found",
			}
		}

		h.System.Logger.Error(err, "failed get user")
		return &fiber.Error{
			Code:    500,
			Message: "failed get user",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success get user",
		Data:    user,
	})
}

func (h *Handlers) UserUpdate(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.System.Logger.Error(err, "failed to parse id")
		return &fiber.Error{
			Code:    400,
			Message: "id not valid",
		}
	}

	payload := new(models.User)
	err = c.BodyParser(payload)
	if err != nil {
		h.System.Logger.Error(err, "failed to parse payload")
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	validationError := h.System.Validator.MustValidate(payload)
	if validationError != nil {
		return validationError
	}

	user, err := h.userService.Update(uint(id), payload)
	if err != nil {
		h.System.Logger.Error(err, "failed update user")
		return &fiber.Error{
			Code:    500,
			Message: "failed update user",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success update user",
		Data:    user,
	})
}

func (h *Handlers) UserDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		h.System.Logger.Error(err, "failed to parse id")
		return &fiber.Error{
			Code:    400,
			Message: "id not valid",
		}
	}

	user := h.GetUser(c)
	if user == nil {
		return &fiber.Error{
			Code:    401,
			Message: "unauthorized",
		}
	}

	if user.ID == uint(id) {
		return &fiber.Error{
			Code:    400,
			Message: "cannot delete yourself",
		}
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		h.System.Logger.Error(err, "failed delete user")
		return &fiber.Error{
			Code:    500,
			Message: "failed delete user",
		}
	}

	return h.SendJson(c, stypes.Response{
		Status:  "success",
		Title:   "Success",
		Message: "success delete user",
	})
}
