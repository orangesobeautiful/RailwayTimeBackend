package controllers

import ptxrailwaymodels "RailwayTime/models/ptx-railway-models"

// StationInfo 站點資訊
type StationInfo struct {
	Name      NameType // 車站名稱
	StationID string   // 車站 ID
}

// RegionInfo 地區站點資料
type RegionInfo struct {
	Name        string         // 地區名稱
	StationList []*StationInfo // 車站列表
}

// TrainInfo 車次資料
type TrainInfo struct {
	TrainNo             string   // 車次代碼
	Direction           int32    // 行駛方向 : [0:'順行',1:'逆行']
	TrainTypeName       NameType // 車種名稱
	TripHeadSign        string   // 車次之目的地方向描述
	StartingStationID   string   // 車次之起始站車站代號
	StartingStationName NameType // 車次之起始站車站名稱
	EndingStationID     string   // 車次之終點站車站代號
	EndingStationName   NameType // 車次之終點站車站名稱
	TripLine            int32    // 山海線類型 : [0:'不經山海線',1:'山線',2:'海線',3:'成追線']
	DailyFlag           int32    // 是否每日行駛 : [0:'否',1:'是']
	ExtraTrainFlag      int32    // 是否為加班車 : [0:'否',1:'是']
	Note                string   // 附註說明
}

// TrainTimetable 列車時刻表
type TrainTimetable struct {
	TrainInfo                  *TrainInfo                     // 車次資訊
	OriginStationStopTime      *ptxrailwaymodels.StopTimeInfo // 出發站時刻
	DestinationStationStopTime *ptxrailwaymodels.StopTimeInfo // 抵達站時刻
	DelayMinute                int32                          // 延誤時間(分鐘)
}
