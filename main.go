package main

import (
	"flag"

	"github.com/labstack/echo/v4"
)

func main() {
	flag.Parse()

	e := echo.New()
	config, err := InitConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	InitRoute(e, config)
	e.Logger.Fatal(e.Start(config.Addr))
}
