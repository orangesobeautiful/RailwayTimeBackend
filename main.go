package main

import (
	"RailwayTime/server"
	"flag"
	"log"
)

func main() {
	var cliDebugMode bool
	flag.BoolVar(&cliDebugMode, "debug", false, "")
	flag.Parse()

	err := server.StartServer(cliDebugMode)
	if err != nil {
		log.Fatal("server.StartServer failed, err", err)
	}
}
