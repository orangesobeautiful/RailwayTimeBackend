package ptxlib

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
	"fmt"
)

// GetStationInfo 取得站點資料
func (p *PTXController) GetStationInfo() (stationListRsp *ptxrailwaymodels.PTXStationListResponse, err error) {
	apiURL := "https://ptx.transportdata.tw/MOTC/v3/Rail/TRA/Station"
	stationListRsp = &ptxrailwaymodels.PTXStationListResponse{}
	err = p.JSONGet(apiURL, nil, &stationListRsp)
	return
}

// GetODStationTimetable 取得指定日期起訖站時刻表
func (p *PTXController) GetODStationTimetable(originStationID string, destinationStationID string, trainDate string) (timetableListRsp *ptxrailwaymodels.PTXDailyTrainTimeTableListResponse, err error) {
	apiURL := fmt.Sprintf("https://ptx.transportdata.tw/MOTC/v3/Rail/TRA/DailyTrainTimetable/OD/%s/to/%s/%s", originStationID, destinationStationID, trainDate)
	timetableListRsp = &ptxrailwaymodels.PTXDailyTrainTimeTableListResponse{}
	err = p.JSONGet(apiURL, nil, &timetableListRsp)
	return
}
