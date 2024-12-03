package controllers

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/omniscan/app/controllers/api"
	"github.com/mrrizkin/omniscan/app/controllers/scan"
	"github.com/mrrizkin/omniscan/pkg/boot/constructor"
)

func New() fx.Option {
	return constructor.Load(
		&api.EStatementController{},

		&scan.EStatementController{},

		&UserController{},
		&DashboardController{},
	)
}
