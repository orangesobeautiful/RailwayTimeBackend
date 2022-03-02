package ptxrailwaymodels

import "time"

// PTXStationListResponse 車站回應結構
type PTXStationListResponse struct {
	PTXCommonResponseInfo `json:",inline"`
	StationList           []*StationInfo `json:"Stations"` // 車站資料
}

// StationInfo 車站資訊結構
type StationInfo struct {
	StationUID      string    `json:"StationUID"`      // 車站唯一識別代碼
	StationID       string    `json:"StationID"`       // 車站代碼
	ReservationCode string    `json:"ReservationCode"` // 訂票車站代碼
	StationName     NameType  `json:"StationName"`     // 車站名稱
	StationPosition PointType `json:"StationPosition"` // 車站座標
	StationAddress  string    `json:"StationAddress"`  // 車站地址
	StationPhone    string    `json:"StationPhone"`    // 車站聯絡電話
	StationClass    string    `json:"StationClass"`    // 車站級別 = ['0: 特等', '1: 一等', '2: 二等', '3: 三等', '4: 簡易', '5: 招呼', '6: 號誌', 'A: 貨運', 'B: 基地', 'X: 非車']
	StationURL      string    `json:"StationURL"`      // 車站資訊說明網址

	Region string `json:"-"` // 自定義地區(ex: 台北、新竹、高雄...)
}

// RegionInfo 地區站點資訊
type RegionInfo struct {
	Name        NameType       // 地區名稱
	StationList []*StationInfo // StationList 站點列表
}

// StopTimeInfo 停靠站資料
type StopTimeInfo struct {
	StopSequence  int32    `json:"StopSequence"`  // 停靠站序
	StationID     string   `json:"StationID"`     // 車站代碼
	StationName   NameType `json:"StationName"`   // 車站名稱
	ArrivalTime   string   `json:"ArrivalTime"`   // 到站時間
	DepartureTime string   `json:"DepartureTime"` // 離站時間
}

// SectionAmongInfo 站點所屬區間
type SectionAmongInfo struct {
	StartStationID string `json:"StartStationID"` // 起站車站代碼
	EndStationID   string `json:"EndStationID"`   // 迄站車站代碼
}

// DiningFlagSectionInfo 提供訂便當服務之車站區間
type DiningFlagSectionInfo struct {
	StartSection SectionAmongInfo `json:"StartSection"` // 乘客出發站所屬區間
	EndSection   SectionAmongInfo `json:"EndSection"`   // 乘客目的站所屬區間
}

// TrainInfo 車次資料
type TrainInfo struct {
	TrainNo               string                   `json:"TrainNo"`             // 車次代碼
	RouteID               string                   `json:"RouteID"`             // 營運路線代碼
	Direction             int32                    `json:"Direction"`           // 行駛方向 : [0:'順行',1:'逆行']
	TrainTypeID           string                   `json:"TrainTypeID"`         // 車種代嗎
	TrainTypeCode         string                   `json:"TrainTypeCode"`       // 車種簡碼 = ['1: 太魯閣', '2: 普悠瑪', '3: 自強', '4: 莒光', '5: 復興', '6: 區間', '7: 普快', '10: 區間快']
	TrainTypeName         NameType                 `json:"TrainTypeName"`       // 車種名稱
	TripHeadSign          string                   `json:"TripHeadSign"`        // 車次之目的地方向描述
	StartingStationID     string                   `json:"StartingStationID"`   // 車次之起始站車站代號
	StartingStationName   NameType                 `json:"StartingStationName"` // 車次之起始站車站名稱
	EndingStationID       string                   `json:"EndingStationID"`     // 車次之終點站車站代號
	EndingStationName     NameType                 `json:"EndingStationName"`   // 車次之終點站車站名稱
	OverNightStationID    string                   `json:"OverNightStationID"`  // 跨夜車站代碼
	TripLine              int32                    `json:"TripLine"`            // 山海線類型 : [0:'不經山海線',1:'山線',2:'海線',3:'成追線']
	WheelChairFlag        int32                    `json:"WheelChairFlag"`      // 是否設身障旅客專用座位車 : [0:'否',1:'是']
	PackageServiceFlag    int32                    `json:"PackageServiceFlag"`  // 是否提供行李服務 : [0:'否',1:'是']
	DiningFlag            int32                    `json:"DiningFlag"`          // 是否提供訂便當服務 : [0:'否',1:'是']
	DiningFlagSectionList []*DiningFlagSectionInfo `json:"DiningFlagSections"`  // 提供訂便當服務之車站區間
	BreastFeedFlag        int32                    `json:"BreastFeedFlag"`      // 是否設有哺(集)乳室車廂 : [0:'否',1:'是']
	BikeFlag              int32                    `json:"BikeFlag"`            // 是否人車同行班次(置於攜車袋之自行車各級列車均可乘車) : [0:'否',1:'是']
	CarFlag               int32                    `json:"CarFlag"`             // 是否提供小汽車上火車服務 : [0:'否',1:'是']
	DailyFlag             int32                    `json:"DailyFlag"`           // 是否每日行駛 : [0:'否',1:'是']
	ExtraTrainFlag        int32                    `json:"ExtraTrainFlag"`      // 是否為加班車 : [0:'否',1:'是']
	Note                  string                   `json:"Note"`                // 附註說明
}

// TrainTimetableInfo 車站時刻表資訊
type TrainTimetableInfo struct {
	TrainInfo    TrainInfo       `json:"TrainInfo"` // 車次資料
	StopTimeList []*StopTimeInfo `json:"StopTimes"` // 停靠站資料
}

// PTXDailyTrainTimeTableListResponse 車站時刻表回應結構
type PTXDailyTrainTimeTableListResponse struct {
	PTXCommonResponseInfo `json:",inline"`
	TrainTimetableList    []*TrainTimetableInfo `json:"TrainTimetables"` // 車站時刻列表
}

// TrainLiveBoardInfo 列車即時動態位置資訊
type TrainLiveBoardInfo struct {
	TrainNo            string    `json:"TrainNo"`            // 車次代碼
	TrainTypeID        string    `json:"TrainTypeID"`        // 車種代嗎
	TrainTypeCode      string    `json:"TrainTypeCode"`      // 車種簡碼 = ['1: 太魯閣', '2: 普悠瑪', '3: 自強', '4: 莒光', '5: 復興', '6: 區間', '7: 普快', '10: 區間快']
	TrainTypeName      NameType  `json:"TrainTypeName"`      // 車種名稱
	StationID          string    `json:"StationID"`          // 經過/停靠車站代號
	StationName        NameType  `json:"StationName"`        // 經過/停靠車站名稱
	TrainStationStatus int32     `json:"TrainStationStatus"` // 列車目前所在之車站狀態 : [0:'進站中',1:'在站上',2:'已離站']
	DelayTime          int32     `json:"DelayTime"`          // 延誤分鐘
	UpdateTime         time.Time `json:"UpdateTime"`         // 本筆位置資料之更新日期時間
}

// PTXTrainLiveBoardListResponse 列車即時位置動態回傳結構
type PTXTrainLiveBoardListResponse struct {
	PTXCommonResponseInfo `json:",inline"`
	TrainLiveBoardList    []*TrainLiveBoardInfo `json:"TrainLiveBoards"` // 列車即時位置列表
}
