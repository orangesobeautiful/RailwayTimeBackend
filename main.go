package main

import (
	"RailwayTime/server"
	"flag"
)

func main() {
	var cliDebugMode bool
	flag.BoolVar(&cliDebugMode, "debug", false, "")
	flag.Parse()
	server.StartServer(cliDebugMode)
}
