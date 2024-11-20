package services

import (
	"go.uber.org/fx"

	estatement "github.com/mrrizkin/omniscan/app/services/e-statement"
	"github.com/mrrizkin/omniscan/pkg/boot/constructor"
)

func New() fx.Option {
	return constructor.Load(
		&UserService{},
		&estatement.EStatementService{},
	)
}
