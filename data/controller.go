package data

import (
	"RailwayTime/tdxlib"
)

// Controller Controller
type Controller struct {
	tdxController *tdxlib.TDXController // 讀取 PTX API 的 controller

	stationCache            *StationCache
	odStationTimetableCache *ODStationTimeableCache
	trainLiveBoradCache     *TrainLiveBoradCache
}

// NewDataController 初始化資料和設定 timer 自動更新資料
func NewDataController(cID, cSEC string) (ctrl *Controller, err error) {
	ctrl = &Controller{}
	ctrl.tdxController, err = tdxlib.NewTDXController(cID, cSEC)
	if err != nil {
		return
	}

	ctrl.trainLiveBoradCache = newTrainLiveBoradCache(ctrl.tdxController)
	ctrl.odStationTimetableCache = newODStationTimeable(ctrl.tdxController)
	ctrl.stationCache = newStationCache(ctrl.tdxController)

	return
}
