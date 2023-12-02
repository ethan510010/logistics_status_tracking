package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"logistics_status_tracking/internal/infra/po"
	"logistics_status_tracking/internal/service/handlers"

	"github.com/labstack/echo/v4"
)

func (s *Server) queryHandler(c echo.Context) error {
	return nil
}

func (s *Server) fakeHandler(c echo.Context) error {
	numStr := c.Param("num")
	number, err := strconv.Atoi(numStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &FakeList{
			ErrorMsg: "invalid number input",
			Data:     nil,
		})
	}
	fakeDataService := handlers.NewFakeDataService(s.db)
	recipients, err := fakeDataService.ParseRecipient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &FakeList{
			ErrorMsg: "gen failed",
			Data:     nil,
		})
	}
	locations, err := fakeDataService.ParseLocation()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &FakeList{
			ErrorMsg: "gen failed",
			Data:     nil,
		})
	}
	res := []FakeResponse{}
	statusList, err := s.genTrackingStatusList(number, locations, recipients)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &FakeList{
			ErrorMsg: "gen failed",
		})
	}

	for _, status := range statusList {
		res = append(res, FakeResponse{
			Sno:            status.Sno,
			TrackingStatus: status.Status,
		})
	}

	return c.JSON(http.StatusOK, &FakeList{
		Data: res,
	})
}

func (s *Server) genTrackingStatusList(num int, locations []po.Location, recipients []po.Recipient) ([]po.TrackingStatus, error) {
	recipientIDs := []uint32{}
	locationIDs := []uint32{}

	for _, location := range locations {
		locationIDs = append(locationIDs, location.LocationID)
	}
	for _, recipient := range recipients {
		recipientIDs = append(recipientIDs, recipient.ID)
	}

	details := []po.Detail{}
	statusList := []po.TrackingStatus{}
	sno := uuid.New().ID()
	for i := 0; i < num; i++ {
		details = genDetails(sno, locations)
		status := int8(po.DeliverStatusList[rand.Intn(len(po.DeliverStatusList))])
		t := po.TrackingStatus{
			Sno:                   sno,
			Status:                status,
			EstimatedDeliveryTime: randomDateStr(),
			RecipientID:           recipientIDs[rand.Intn(len(recipientIDs))],
			CurrentLocationID:     locationIDs[rand.Intn(len(locationIDs))],
		}
		statusList = append(statusList, t)
	}
	tx := s.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("begin tx failed: %s", tx.Error.Error())
	}
	if err := tx.CreateInBatches(details, len(details)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create fake details failed: %s", err.Error())
	}
	if err := tx.CreateInBatches(statusList, len(statusList)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create fake tracking status list failed: %s", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("tx commit failed: %s", err.Error())
	}
	return statusList, nil
}

func randomDateStr() string {
	year := rand.Intn(2023-2000) + 2000 // 假資料是 2000 到 2023 年
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1 // 簡單處理，假设每个月都有 28 天

	// 创建日期
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// 格式化日期
	dateStr := date.Format("2006-01-02")
	return dateStr
}

func randomHrStr() string {
	hour := rand.Intn(24)   // 小时范围：0 到 23
	minute := rand.Intn(60) // 分钟范围：0 到 59

	// 格式化为 HH:MM 形式
	randomTime := fmt.Sprintf("%02d:%02d", hour, minute)

	return randomTime
}

func genDetails(sno uint32, locations []po.Location) []po.Detail {
	var details []po.Detail
	count := rand.Intn(5)
	for i := 0; i < count; i++ {
		details = append(details, po.Detail{
			Date:          randomDateStr(),
			TimeHour:      randomHrStr(),
			Status:        int8(po.DeliverStatusList[rand.Intn(len(po.DeliverStatusList))]),
			LocationID:    locations[rand.Intn(len(locations))].LocationID,
			LocationTitle: locations[rand.Intn(len(locations))].Title,
			Sno:           sno,
		})
	}

	return details
}
