package server

import (
	"RailwayTime/controllers"
	"fmt"
	"log"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
)

var serverSetting SettingProperties

// StartServer start server
func StartServer(debugMode bool) (err error) {
	// read setting
	serverSetting, err = readSettingFile()
	if err != nil {
		log.Fatal(err)
	}

	// init iris router
	app := iris.New()
	if debugMode {
		app.Logger().SetLevel("debug")
		crs := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		})
		app.Use(crs)
	}
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(iris.Compression)

	// init railway controller
	rwController, err := controllers.NewRailwayController(serverSetting.CID, serverSetting.CSEC)
	if err != nil {
		log.Fatal(err)
	}
	m := mvc.New(app.Party("/api/railway"))
	m.Handle(rwController)

	// iris listen
	hostPort := fmt.Sprintf("%s:%d", serverSetting.ServerHost, serverSetting.ServerPort)
	err = app.Listen(hostPort, iris.WithOptimizations)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
