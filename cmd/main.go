package main

import (
	"github.com/CyrilKuzmin/itpath69/app"
	"github.com/CyrilKuzmin/itpath69/config"
)

func main() {
	conf := config.Get()
	a := app.NewApp(conf)
	a.Start()
}
