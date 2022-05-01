package main

import (
	"github.com/CyrilKuzmin/itpath69/config"
	"github.com/CyrilKuzmin/itpath69/server"
)

func main() {
	conf := config.Get()
	s := server.NewApp(conf)
	s.Start()
}
