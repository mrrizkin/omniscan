package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/omniscan/app/providers/session"
)

type Authentication struct {
	session *session.Session
}

func (*Authentication) Construct() interface{} {
	return func(session *session.Session) *Authentication {
		return &Authentication{
			session: session,
		}
	}
}

func (a *Authentication) Protect(c *fiber.Ctx) error {
	session, err := a.session.Get(c)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get session: %s", err),
		}
	}

	uid, ok := session.Get("uid").(uint)
	if !ok {
		return &fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	sid, ok := session.Get("sid").(string)
	if !ok {
		return &fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	c.Locals("uid", uid)
	c.Locals("sid", sid)

	return c.Next()
}
