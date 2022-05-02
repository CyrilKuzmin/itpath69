package main

import (
	"github.com/CyrilKuzmin/itpath69/config"
	"github.com/CyrilKuzmin/itpath69/service"
)

func main() {
	conf := config.Get()
	s := service.NewApp(conf)
	s.Start()
}
