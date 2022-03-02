package data

import (
	"RailwayTime/ptxlib"
)

// Controller Controller
type Controller struct {
	ptxController *ptxlib.PTXController // 讀取 PTX API 的 controller

	stationCache            *StationCache
	odStationTimetableCache *ODStationTimeableCache
	trainLiveBoradCache     *TrainLiveBoradCache
}

// NewDataController 初始化資料和設定 timer 自動更新資料
func NewDataController(id, key string) (ctrl *Controller, err error) {
	ctrl = &Controller{}
	ctrl.ptxController = ptxlib.NewPTXController(id, key)

	ctrl.trainLiveBoradCache = newTrainLiveBoradCache(ctrl.ptxController)
	ctrl.odStationTimetableCache = newODStationTimeable(ctrl.ptxController)
	ctrl.stationCache = newStationCache(ctrl.ptxController)

	return
}
