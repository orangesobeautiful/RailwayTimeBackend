package data

import (
	"RailwayTime/ptxlib"
	"math"
	"sync"
	"time"
)

const reUpdateSecondWhenErr = 15

type cacheBaseUnit struct {
	ptxController *ptxlib.PTXController // 讀取 PTX API 的 controller

	dataTime        time.Time     // 資料最後更新時間(發生錯誤時不會記錄)
	lastUpdateTime  time.Time     // 最後更新時間
	lastUpdateError error         // 最後更新時發生的錯誤
	nextUpdateTime  time.Time     // 預計下次的更新時間
	updateTimer     *time.Timer   // 更新計時器
	updateTimerFunc func()        // 更新資料的 timer 需要呼叫的 func
	lock            *sync.RWMutex // 讀寫鎖
}

func (unit *cacheBaseUnit) setNextUpdateTimer(d time.Duration) {
	unit.nextUpdateTime = time.Now().Add(d)
	unit.updateTimer = time.AfterFunc(d, unit.updateTimerFunc)
}

// setNextUpdateTimerByPTXRsp 根據 PTX 資料回傳的資料更新時間設定下次更新 timer
func (unit *cacheBaseUnit) setNextUpdateTimerByPTXInfo(ptxUpdateTime time.Time, ptxUpdateInterval int32, srcUpdateTime time.Time, srcUpdateInterval int32) {
	nowTime := time.Now()
	var nextUpdateInterval time.Duration
	if srcUpdateInterval > 0 {
		// 資料來源端下次更新時間
		srcNextUpdateTime := srcUpdateTime.Add(time.Duration(srcUpdateInterval) * time.Second)
		// PTX 平台下次更新來源端最新資料的時間間隔次數
		intervalTimes :=
			int32(math.Ceil(srcNextUpdateTime.Sub(ptxUpdateTime).Seconds() / float64(ptxUpdateInterval)))
		nextUpdateInterval = ptxUpdateTime.Add(time.Duration(intervalTimes*ptxUpdateInterval)*time.Second).Sub(nowTime) + (3 * time.Second)
	} else {
		// 小於等於零代表不定期更新
		nextUpdateInterval = ptxUpdateTime.Add(time.Duration(ptxUpdateInterval)*time.Second).Sub(nowTime) + (3 * time.Second)
	}

	if nextUpdateInterval < 0 {
		nextUpdateInterval = 33 * time.Second
	} else if nextUpdateInterval > 86000*time.Second {
		nextUpdateInterval = 86000 * time.Second
	}

	unit.setNextUpdateTimer(nextUpdateInterval)
}
