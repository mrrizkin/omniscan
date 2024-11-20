package main

import (
	"github.com/mrrizkin/omniscan/bootstrap"
)

func main() {
	app := bootstrap.App()
	if app.Err() != nil {
		panic(app.Err())
	}
	app.Run()
}
