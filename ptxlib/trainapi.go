package ptxlib

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
)

// GetTrainLiveBoard 取得列車即時位置動態資料
func (p *PTXController) GetTrainLiveBoard() (trainLiveBoradListRsp *ptxrailwaymodels.PTXTrainLiveBoardListResponse, err error) {
	apiURL := "https://ptx.transportdata.tw/MOTC/v3/Rail/TRA/TrainLiveBoard"
	trainLiveBoradListRsp = &ptxrailwaymodels.PTXTrainLiveBoardListResponse{}
	err = p.JSONGet(apiURL, nil, &trainLiveBoradListRsp)
	return
}
