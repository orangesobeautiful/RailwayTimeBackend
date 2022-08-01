package tdxlib

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
	"fmt"
)

// GetStationInfo 取得站點資料
func (tc *TDXController) GetStationInfo() (stationListRsp *ptxrailwaymodels.PTXStationListResponse, err error) {
	const apiURL string = apiBasicBaseURL + "/v3/Rail/TRA/Station"
	stationListRsp = &ptxrailwaymodels.PTXStationListResponse{}
	err = tc.JSONGet(apiURL, nil, &stationListRsp)
	return
}

// GetODStationTimetable 取得指定日期起訖站時刻表
func (tc *TDXController) GetODStationTimetable(originStationID, destinationStationID, trainDate string) (
	timetableListRsp *ptxrailwaymodels.PTXDailyTrainTimeTableListResponse, err error) {
	apiURL := fmt.Sprintf(apiBasicBaseURL+"/v3/Rail/TRA/DailyTrainTimetable/OD/%s/to/%s/%s", originStationID, destinationStationID, trainDate)
	timetableListRsp = &ptxrailwaymodels.PTXDailyTrainTimeTableListResponse{}
	err = tc.JSONGet(apiURL, nil, &timetableListRsp)
	return
}
