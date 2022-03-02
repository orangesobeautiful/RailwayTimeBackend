package data

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
)

// GetStationList 獲取站點列表
func (ctrl *Controller) GetStationList() (stationList []ptxrailwaymodels.StationInfo, err error) {
	return ctrl.stationCache.GetStationList()
}

// GetRegionList 獲取地區站點列表
func (ctrl *Controller) GetRegionList() (regionList []ptxrailwaymodels.RegionInfo, err error) {
	return ctrl.stationCache.GetRegionList()
}

// GetStationTimetableOD 獲取指定日期起迄站時刻表
func (ctrl *Controller) GetStationTimetableOD(originStationID string, destinationStationID string, trainDate string) (trainTimetableList []ptxrailwaymodels.TrainTimetableInfo, err error) {
	timetableListRspInfo, err := ctrl.odStationTimetableCache.GetODStationTimetable(originStationID, destinationStationID, trainDate)
	if err != nil {
		return
	}
	for _, trainTimetable := range timetableListRspInfo.TrainTimetableList {
		trainTimetableList = append(trainTimetableList, *trainTimetable)
	}
	return
}

// GetAllTrainLiveBoard 獲取列車即時位置動態資料
func (ctrl *Controller) GetAllTrainLiveBoard() (allTrainLiveInfo map[string]ptxrailwaymodels.TrainLiveBoardInfo, err error) {
	allTrainLiveInfo = ctrl.trainLiveBoradCache.GetAllTrainLiveBorad()
	if err != nil {
		return
	}

	return
}
