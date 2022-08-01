package controllers

import (
	"RailwayTime/data"
	"errors"
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
func NewRailwayController(cID, cSEC string) (rwController *RailwayController, err error) {
	dController, err := data.NewDataController(cID, cSEC)
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
	for sIdx := range ptxStationList {
		stationList = append(stationList, &StationInfo{
			Name:      NameType(ptxStationList[sIdx].StationName),
			StationID: ptxStationList[sIdx].StationID,
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
func (ctrl *RailwayController) GetTimetableOD(originStationID, destinationStationID string) (
	trainTimetableList []*TrainTimetable, err error) {
	// 取得車站時刻表
	ptxTimetableList, err := ctrl.dataController.GetStationTimetableOD(
		originStationID, destinationStationID, time.Now().In(ctrl.cstLocation).Format("2006-01-02"))
	if err != nil {
		return
	}
	// 取得列車延誤資訊
	allTrainLiveInfo, err := ctrl.dataController.GetAllTrainLiveBoard()
	if err != nil {
		log.Printf("GetAllTrainLiveBoard failed, err=%s", err)
		err = errors.New("internal server error")
		return
	}

	// 資料轉換
	for ttIdx := range ptxTimetableList {
		if len(ptxTimetableList[ttIdx].StopTimeList) < 2 {
			log.Printf("the len of ptxTimetable.StopTimeList < 2 in GetTimetableOD")
			err = fmt.Errorf("internal server error")
			return
		}

		trainNo := ptxTimetableList[ttIdx].TrainInfo.TrainNo
		var delayTime int32
		if trainLiveBorad, exist := allTrainLiveInfo[trainNo]; exist {
			delayTime = trainLiveBorad.DelayMinute
		}

		trainTimetableList = append(trainTimetableList, &TrainTimetable{
			TrainInfo: &TrainInfo{
				TrainNo:             ptxTimetableList[ttIdx].TrainInfo.TrainNo,
				Direction:           ptxTimetableList[ttIdx].TrainInfo.Direction,
				TrainTypeName:       NameType(ptxTimetableList[ttIdx].TrainInfo.TrainTypeName),
				TripHeadSign:        ptxTimetableList[ttIdx].TrainInfo.TripHeadSign,
				StartingStationID:   ptxTimetableList[ttIdx].TrainInfo.StartingStationID,
				StartingStationName: NameType(ptxTimetableList[ttIdx].TrainInfo.StartingStationName),
				EndingStationID:     ptxTimetableList[ttIdx].TrainInfo.EndingStationID,
				EndingStationName:   NameType(ptxTimetableList[ttIdx].TrainInfo.EndingStationName),
				TripLine:            ptxTimetableList[ttIdx].TrainInfo.TripLine,
				DailyFlag:           ptxTimetableList[ttIdx].TrainInfo.DailyFlag,
				ExtraTrainFlag:      ptxTimetableList[ttIdx].TrainInfo.ExtraTrainFlag,
				Note:                ptxTimetableList[ttIdx].TrainInfo.Note,
			},
			OriginStationStopTime:      ptxTimetableList[ttIdx].StopTimeList[0],
			DestinationStationStopTime: ptxTimetableList[ttIdx].StopTimeList[1],
			DelayMinute:                delayTime,
		})
	}

	return trainTimetableList, nil
}
