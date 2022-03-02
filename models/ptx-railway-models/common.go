package ptxrailwaymodels

import "time"

// NameType 多語名稱結構
type NameType struct {
	ZhTW string `json:"Zh_tw"` // 繁中名稱
	En   string `json:"En"`    // 英文名稱
}

// PointType 座標結構
type PointType struct {
	PositionLon float64 `json:"PositionLon"` // 位置經度(WGS84)
	PositionLat float64 `json:"PositionLat"` // 位置緯度(WGS84)
}

// PTXCommonResponseInfo PTX API 回應結構共有數據結構
type PTXCommonResponseInfo struct {
	UpdateTime        time.Time `json:"UpdateTime"`        // PTX 平台資料更新時間(ISO8601格式:yyyy-MM-ddTHH:mm:sszzz)
	UpdateInterval    int32     `json:"UpdateInterval"`    // PTX 平台資料更新週期(秒)
	SrcUpdateTime     time.Time `json:"SrcUpdateTime"`     // 來源端平台資料更新時間(ISO8601格式:yyyy-MM-ddTHH:mm:sszzz)
	SrcUpdateInterval int32     `json:"SrcUpdateInterval"` // 來源端平台資料更新週期(秒)['-1: 不定期更新']
	AuthorityCode     string    `json:"AuthorityCode"`     // 業管機關簡碼
	Count             int64     `json:"Count"`             // 資料總筆數
}
