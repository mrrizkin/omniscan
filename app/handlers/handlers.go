package handlers

import (
	"github.com/mrrizkin/omniscan/app/domains/ocr"
	"github.com/mrrizkin/omniscan/app/domains/user"
	"github.com/mrrizkin/omniscan/system/stypes"
)

type Handlers struct {
	*stypes.App

	userRepo    *user.Repo
	userService *user.Service

	ocrRepo    *ocr.Repo
	ocrService *ocr.Service
}

func New(
	app *stypes.App,
) *Handlers {
	userRepo := user.NewRepo(app.System.Database)
	userService := user.NewService(userRepo, app.Library.Hashing)

	ocrRepo := ocr.NewRepo(app.System.Database)
	ocrService := ocr.NewService(ocrRepo, app.Library.MutasiScanner)

	return &Handlers{
		App: app,

		userRepo:    userRepo,
		userService: userService,

		ocrRepo:    ocrRepo,
		ocrService: ocrService,
	}
}
