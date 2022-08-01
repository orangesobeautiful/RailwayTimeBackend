package data

import (
	ptxrailwaymodels "RailwayTime/models/ptx-railway-models"
	"RailwayTime/tdxlib"
	"fmt"
	"regexp"
	"sync"
	"time"
)

// StationCache 車站資料的 cache 控制單元
type StationCache struct {
	*cacheBaseUnit
	stationData      map[string]*ptxrailwaymodels.StationInfo // 站點資料
	stationNameOrder []string                                 // 紀錄站點的順序
	regionData       map[string]*ptxrailwaymodels.RegionInfo  // 地區站點資料
	regionNameOrder  []string                                 // 紀錄地區的順序
}

func newStationCache(tdxCtrl *tdxlib.TDXController) (cache *StationCache) {
	cache = &StationCache{
		&cacheBaseUnit{
			ptxController: tdxCtrl,
			lock:          &sync.RWMutex{},
		},
		make(map[string]*ptxrailwaymodels.StationInfo),
		[]string{},
		make(map[string]*ptxrailwaymodels.RegionInfo),
		[]string{},
	}
	cache.updateTimerFunc = cache.update
	cache.update()

	return
}

// update 更新資料
func (cache *StationCache) update() {
	timeNow := time.Now()

	// 取得新資料並檢查錯誤
	stationListResponse, err := cache.ptxController.GetStationInfo()
	cache.lastUpdateTime = timeNow
	cache.lastUpdateError = err
	if err != nil {
		cache.setNextUpdateTimer(reUpdateSecondWhenErr * time.Second)
		return
	}

	// 更新資料
	cache.dataTime = timeNow
	// 初始化新資料
	var newStationData = make(map[string]*ptxrailwaymodels.StationInfo)
	var newStationNameOrder []string
	var newRegionData = make(map[string]*ptxrailwaymodels.RegionInfo)
	var newRegionNameOrder []string

	for _, stationInfo := range stationListResponse.StationList {
		newStationData[stationInfo.StationID] = stationInfo
		newStationNameOrder = append(newStationNameOrder, stationInfo.StationID)
		var addressRegex *regexp.Regexp
		addressRegex, err = regexp.Compile("([0-9]+)([^縣市]+)(.*)")
		if err != nil {
			err = fmt.Errorf("regexp compile failed, error=%s", err)
			return
		}
		ZhTwAddress := stationInfo.StationAddress
		addressMatch := addressRegex.FindSubmatch([]byte(ZhTwAddress))
		if len(addressMatch) > 2 {
			zhtwRegionName := string(addressMatch[2])
			if _, regionExist := newRegionData[zhtwRegionName]; regionExist {
				newRegionData[zhtwRegionName].StationList = append(newRegionData[zhtwRegionName].StationList, stationInfo)
			} else {
				newRegionData[zhtwRegionName] = &ptxrailwaymodels.RegionInfo{
					Name: ptxrailwaymodels.NameType{
						ZhTW: zhtwRegionName,
					},
					StationList: []*ptxrailwaymodels.StationInfo{stationInfo},
				}
				newRegionNameOrder = append(newRegionNameOrder, zhtwRegionName)
			}

		}
	}
	cache.lock.Lock()
	cache.stationData = newStationData
	cache.stationNameOrder = newStationNameOrder
	cache.regionData = newRegionData
	cache.regionNameOrder = newRegionNameOrder
	cache.lock.Unlock()

	// 設定下次更新資料時間
	cache.setNextUpdateTimerByPTXInfo(stationListResponse.UpdateTime, stationListResponse.UpdateInterval,
		stationListResponse.SrcUpdateTime, stationListResponse.SrcUpdateInterval)

	return
}

// GetStationList 取得車站列表
func (cache *StationCache) GetStationList() (stationList []ptxrailwaymodels.StationInfo, err error) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	for _, stationID := range cache.stationNameOrder {
		stationList = append(stationList, *cache.stationData[stationID])
	}

	return
}

// GetRegionList 取得地區站點列表
func (cache *StationCache) GetRegionList() (regionList []ptxrailwaymodels.RegionInfo, err error) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	for _, regionZhTwName := range cache.regionNameOrder {
		regionList = append(regionList, *cache.regionData[regionZhTwName])
	}

	return
}
