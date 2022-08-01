package tdxlib

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
)

// GetTrainLiveBoard 取得列車即時位置動態資料
func (tc *TDXController) GetTrainLiveBoard() (trainLiveBoradListRsp *ptxrailwaymodels.PTXTrainLiveBoardListResponse, err error) {
	const apiURL string = apiBasicBaseURL + "/v3/Rail/TRA/TrainLiveBoard"
	trainLiveBoradListRsp = &ptxrailwaymodels.PTXTrainLiveBoardListResponse{}
	err = tc.JSONGet(apiURL, nil, &trainLiveBoradListRsp)
	return
}
