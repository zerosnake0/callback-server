package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"callback-server/pkg/log"
	"callback-server/pkg/server"
)

func main() {
	log.InitLog()

	var port int
	var debug bool
	flag.IntVar(&port, "port", 80, "listen port")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server.Run(port)
}
