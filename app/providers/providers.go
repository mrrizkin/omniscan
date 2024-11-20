package providers

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/omniscan/pkg/boot/constructor"

	"github.com/mrrizkin/omniscan/app/providers/app"
	"github.com/mrrizkin/omniscan/app/providers/asset"
	"github.com/mrrizkin/omniscan/app/providers/cache"
	"github.com/mrrizkin/omniscan/app/providers/database"
	"github.com/mrrizkin/omniscan/app/providers/hashing"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/app/providers/scheduler"
	"github.com/mrrizkin/omniscan/app/providers/session"
	"github.com/mrrizkin/omniscan/app/providers/validator"
	"github.com/mrrizkin/omniscan/app/providers/view"
)

func New() fx.Option {
	return constructor.Load(
		&app.App{},
		&asset.Asset{},
		&cache.Cache{},
		&database.Database{},
		&hashing.Hashing{},
		&logger.Logger{},
		&scheduler.Scheduler{},
		&session.Session{},
		&validator.Validator{},
		&view.View{},
	)
}
