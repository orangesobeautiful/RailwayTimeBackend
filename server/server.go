package server

import (
	"RailwayTime/controllers"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
)

var serverSetting SettingProperties

func prettyPrint(v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Println("failed to convert to json")
		return
	}
	fmt.Println(string(jsonBytes))
}

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
	}
	app.Use(recover.New())
	app.Use(logger.New())

	// init railway controller
	rwController, err := controllers.NewRailwayController(serverSetting.APPID, serverSetting.APPKey)
	if err != nil {
		log.Fatal(err)
	}
	m := mvc.New(app.Party("/api/railway"))
	m.Handle(rwController)

	// iris listen
	hostPort := fmt.Sprintf("%s:%d", serverSetting.ServerHost, serverSetting.ServerPort)
	err = app.Listen(hostPort)
	if err != nil {
		log.Fatal(err)
	}
	return
}
