package controllers

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/omniscan/pkg/boot/constructor"
)

func New() fx.Option {
	return constructor.Load(
		&EStatementController{},
		&UserController{},
		&WelcomeController{},
	)
}
