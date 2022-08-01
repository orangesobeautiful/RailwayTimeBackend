package data

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
	"RailwayTime/tdxlib"
	"sort"
	"sync"
	"time"
)

// TrainLiveBoradCache 列車即時動態資料的 cache 控制單元
type TrainLiveBoradCache struct {
	*cacheBaseUnit
	data map[string]*ptxrailwaymodels.TrainLiveBoardInfo
}

func newTrainLiveBoradCache(tdxCtrl *tdxlib.TDXController) (cache *TrainLiveBoradCache) {
	cache = &TrainLiveBoradCache{
		&cacheBaseUnit{
			ptxController: tdxCtrl,
			lock:          &sync.RWMutex{},
		},
		make(map[string]*ptxrailwaymodels.TrainLiveBoardInfo),
	}
	cache.updateTimerFunc = cache.update
	cache.update()

	return
}

// update 更新資料
func (cache *TrainLiveBoradCache) update() {
	// 取得新資料並檢查錯誤
	trainLiveBoardListRsp, err := cache.ptxController.GetTrainLiveBoard()
	if err != nil {
		cache.setNextUpdateTimer(reUpdateSecondWhenErr * time.Second)
		return
	}

	// 更新資料
	var newData = make(map[string]*ptxrailwaymodels.TrainLiveBoardInfo)
	for _, trainLiveBorad := range trainLiveBoardListRsp.TrainLiveBoardList {
		newData[trainLiveBorad.TrainNo] = trainLiveBorad
	}
	cache.lock.Lock()
	cache.data = newData
	cache.lock.Unlock()

	// 設定下次更新資料時間
	cache.setNextUpdateTimerByPTXInfo(trainLiveBoardListRsp.UpdateTime, trainLiveBoardListRsp.UpdateInterval,
		trainLiveBoardListRsp.SrcUpdateTime, trainLiveBoardListRsp.SrcUpdateInterval)
}

// GetTrainLiveBorad 取得指定列車的即時動態資訊
func (cache *TrainLiveBoradCache) GetTrainLiveBorad(trainID string) (trainLiveInfo ptxrailwaymodels.TrainLiveBoardInfo) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	trainLiveInfo = *cache.data[trainID]

	return
}

// GetAllTrainLiveBorad 取得全部列車的即時動態資訊
func (cache *TrainLiveBoradCache) GetAllTrainLiveBorad() (allTrainLiveInfo map[string]ptxrailwaymodels.TrainLiveBoardInfo) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	allTrainLiveInfo = make(map[string]ptxrailwaymodels.TrainLiveBoardInfo)
	for trainID, trainLiveInfo := range cache.data {
		allTrainLiveInfo[trainID] = *trainLiveInfo
	}

	return
}

// ODStationTimeableCache 起訖站時刻表的 cache 控制單元
type ODStationTimeableCache struct {
	*cacheBaseUnit
	cstLocation       *time.Location
	dateTimetableData map[string]map[string]*ptxrailwaymodels.PTXDailyTrainTimeTableListResponse // 不同日期的起訖站點時刻表
}

func newODStationTimeable(tdxCtrl *tdxlib.TDXController) (cache *ODStationTimeableCache) {
	cstLocation, _ := time.LoadLocation("Asia/Taipei")
	cache = &ODStationTimeableCache{
		&cacheBaseUnit{
			ptxController: tdxCtrl,
			lock:          &sync.RWMutex{},
		},
		cstLocation,
		make(map[string]map[string]*ptxrailwaymodels.PTXDailyTrainTimeTableListResponse),
	}

	cache.updateTimerFunc = cache.update
	cache.update()

	return
}

// update 檢查並刪除過舊的資料
func (cache *ODStationTimeableCache) update() {
	var nowCSTTime = time.Now().In(cache.cstLocation)
	cache.lock.RLock()
	todayDate := nowCSTTime.Format("2006-01-02")
	var deleteDateList []string
	for dataDate := range cache.dateTimetableData {
		if todayDate > dataDate {
			deleteDateList = append(deleteDateList, dataDate)
		}
	}
	cache.lock.RUnlock()

	// 刪除過期的資料
	cache.lock.Lock()
	for _, deleteDate := range deleteDateList {
		delete(cache.dateTimetableData, deleteDate)
	}
	cache.lock.Unlock()

	// 設定下次更新資料時間
	nextUpdateTime := time.Date(
		nowCSTTime.Year(), nowCSTTime.Month(), nowCSTTime.Day()+1, 0, 30, 0, nowCSTTime.Nanosecond(), cache.cstLocation)
	cache.setNextUpdateTimer(nextUpdateTime.Sub(nowCSTTime))
}

// odStationKey 起訖站在 cache 中索引的 key
func odStationKey(originStationID, destinationStationID string) string {
	return originStationID + "~" + destinationStationID
}

// GetODStationTimetable 取得指定日期起迄站時刻表
func (cache *ODStationTimeableCache) GetODStationTimetable(originStationID, destinationStationID, trainDate string) (
	trainTimetable *ptxrailwaymodels.PTXDailyTrainTimeTableListResponse, err error) {
	cache.lock.RLock()
	var timetableRsp *ptxrailwaymodels.PTXDailyTrainTimeTableListResponse
	var dateData map[string]*ptxrailwaymodels.PTXDailyTrainTimeTableListResponse
	var dateExist bool
	var odKey = odStationKey(originStationID, destinationStationID)
	if dateData, dateExist = cache.dateTimetableData[trainDate]; dateExist && dateData != nil {
		timetableRsp = dateData[odKey]
	}

	// 有 cache 的資料
	if timetableRsp != nil {
		// 複製資料並回傳
		copyTimetable := *timetableRsp
		trainTimetable = &copyTimetable
		cache.lock.RUnlock()
		return
	}
	cache.lock.RUnlock()

	// 無 cache 的資料, call API 取得
	timetableRsp, err = cache.ptxController.GetODStationTimetable(originStationID, destinationStationID, trainDate)
	if err != nil {
		return
	}

	// 初始化 map
	if !dateExist {
		cache.dateTimetableData[trainDate] = make(map[string]*ptxrailwaymodels.PTXDailyTrainTimeTableListResponse)
	}
	// 根據出發站的離站時間進行排序
	sort.SliceStable(timetableRsp.TrainTimetableList, func(i, j int) bool {
		return timetableRsp.TrainTimetableList[i].StopTimeList[0].DepartureTime < timetableRsp.TrainTimetableList[j].StopTimeList[0].DepartureTime
	})

	// 更新資料
	cache.lock.Lock()
	cache.dateTimetableData[trainDate][odKey] = timetableRsp
	cache.lock.Unlock()
	copyTimetable := *timetableRsp
	trainTimetable = &copyTimetable

	return trainTimetable, nil
}
