package system

import (
	"fmt"

	"github.com/mrrizkin/omniscan/routes"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/config"
	"github.com/mrrizkin/omniscan/system/database"
	"github.com/mrrizkin/omniscan/system/server"
	"github.com/mrrizkin/omniscan/system/session"
	"github.com/mrrizkin/omniscan/system/stypes"
	"github.com/mrrizkin/omniscan/system/validator"
	"github.com/mrrizkin/omniscan/third-party/hashing"
	"github.com/mrrizkin/omniscan/third-party/logger"
	mutasi_scanner "github.com/mrrizkin/omniscan/third-party/mutasi-scanner"
)

func Run() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	log, err := logger.Zerolog(conf)
	if err != nil {
		panic(err)
	}
	sess, err := session.New(conf)
	if err != nil {
		panic(err)
	}
	defer sess.Stop()
	hash := hashing.Argon2(*conf)

	model := models.New(conf, hash)
	db, err := database.New(conf, model, log)
	if err != nil {
		panic(err)
	}
	defer db.Stop()
	err = db.Start()
	if err != nil {
		panic(err)
	}

	valid := validator.New()
	serv := server.New(conf, log)

	mutasiscan := mutasi_scanner.New()

	routes.Setup(&stypes.App{
		App: serv.App,
		System: &stypes.System{
			Logger:    log,
			Database:  db,
			Config:    conf,
			Session:   sess,
			Validator: valid,
		},
		Library: &stypes.Library{
			Hashing:       hash,
			MutasiScanner: mutasiscan,
		},
	}, sess)

	log.Info(fmt.Sprintf("Server is running on port %d", conf.PORT))

	if err := serv.Serve(); err != nil {
		panic(err)
	}
}
