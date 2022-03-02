package controllers

import (
	"RailwayTime/data"
	"fmt"
	"log"
	"time"

	"github.com/kataras/iris/v12/mvc"
)

// RailwayController 鐵路相關 API controller
type RailwayController struct {
	dataController *data.Controller
	cstLocation    *time.Location
}

// NewRailwayController create a new RailwayController
func NewRailwayController(appID string, appKey string) (rwController *RailwayController, err error) {

	dController, err := data.NewDataController(appID, appKey)
	if err != nil {
		return
	}
	cstLocation, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return
	}
	rwController = &RailwayController{
		dataController: dController,
		cstLocation:    cstLocation,
	}

	return
}

// MVC Railway MVC
func (ctrl *RailwayController) MVC(app *mvc.Application) {
	app.Handle(ctrl)
}

// BeforeActivation called once, before the controller adapted to the main application
func (ctrl *RailwayController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/timetable/OD/{originStationID}/to/{destinationStationID}", "GetTimetableOD")
}

// GetStation serves
// Method:   GET
// Resource: /station
func (ctrl *RailwayController) GetStation() (stationList []*StationInfo, err error) {
	ptxStationList, err := ctrl.dataController.GetStationList()
	if err != nil {
		return
	}

	// 資料轉換
	for _, ptxStation := range ptxStationList {
		stationList = append(stationList, &StationInfo{
			Name:      NameType(ptxStation.StationName),
			StationID: ptxStation.StationID,
		})

	}

	return
}

// GetRegion serves
// Method:   GET
// Resource: /region
func (ctrl *RailwayController) GetRegion() (regionList []*RegionInfo, err error) {
	ptxRegionList, err := ctrl.dataController.GetRegionList()
	if err != nil {
		return
	}

	// 資料轉換
	for _, ptxRegion := range ptxRegionList {
		var stationList []*StationInfo
		for _, ptxStation := range ptxRegion.StationList {
			stationList = append(stationList, &StationInfo{
				Name:      NameType(ptxStation.StationName),
				StationID: ptxStation.StationID,
			})
		}

		regionList = append(regionList, &RegionInfo{
			Name:        ptxRegion.Name.ZhTW,
			StationList: stationList,
		})
	}

	return
}

// GetTimetableOD serves
// Method:   GET
// Resource: /timetable/OD/{originStationID}/to/{destinationStationID}
func (ctrl *RailwayController) GetTimetableOD(originStationID string, destinationStationID string) (trainTimetableList []*TrainTimetable, err error) {
	// 取得車站時刻表
	ptxTimetableList, err := ctrl.dataController.GetStationTimetableOD(originStationID, destinationStationID, time.Now().In(ctrl.cstLocation).Format("2006-01-02"))
	if err != nil {
		return
	}
	// 取得列車延誤資訊
	allTrainLiveInfo, err := ctrl.dataController.GetAllTrainLiveBoard()

	// 資料轉換
	for _, ptxTimetable := range ptxTimetableList {
		if len(ptxTimetable.StopTimeList) < 2 {
			log.Printf("the len of ptxTimetable.StopTimeList < 2 in GetTimetableOD")
			err = fmt.Errorf("internal server error")
			return
		}

		trainNo := ptxTimetable.TrainInfo.TrainNo
		var delayTime int32
		if trainLiveBorad, exist := allTrainLiveInfo[trainNo]; exist {
			delayTime = trainLiveBorad.DelayTime
		}

		trainTimetableList = append(trainTimetableList, &TrainTimetable{
			TrainInfo: &TrainInfo{
				TrainNo:             ptxTimetable.TrainInfo.TrainNo,
				Direction:           ptxTimetable.TrainInfo.Direction,
				TrainTypeName:       NameType(ptxTimetable.TrainInfo.TrainTypeName),
				TripHeadSign:        ptxTimetable.TrainInfo.TripHeadSign,
				StartingStationID:   ptxTimetable.TrainInfo.StartingStationID,
				StartingStationName: NameType(ptxTimetable.TrainInfo.StartingStationName),
				EndingStationID:     ptxTimetable.TrainInfo.EndingStationID,
				EndingStationName:   NameType(ptxTimetable.TrainInfo.EndingStationName),
				TripLine:            ptxTimetable.TrainInfo.TripLine,
				DailyFlag:           ptxTimetable.TrainInfo.DailyFlag,
				ExtraTrainFlag:      ptxTimetable.TrainInfo.ExtraTrainFlag,
				Note:                ptxTimetable.TrainInfo.Note,
			},
			OriginStationStopTime:      ptxTimetable.StopTimeList[0],
			DestinationStationStopTime: ptxTimetable.StopTimeList[1],
			DelayTime:                  delayTime,
		})
	}

	return
}
