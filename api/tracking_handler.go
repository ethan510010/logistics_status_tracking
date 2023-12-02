package api

import (
	"encoding/json"
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
	snoStr := c.QueryParam("sno")
	sno, err := strconv.Atoi(snoStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &QueryResponse{
			Status: "error",
			Data:   nil,
			Error: &QueryError{
				Code:    http.StatusBadRequest,
				Message: "invalid sno input",
			},
		})
	}

	cacheData, err := s.cacheStore.HGetAll(c.Request().Context(), "query", uint32(sno))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &QueryResponse{
			Status: "error",
			Data:   nil,
			Error: &QueryError{
				Code:    http.StatusInternalServerError,
				Message: "query failed",
			},
		})
	}

	data := &QueryData{}
	if len(cacheData) == 0 {
		type QueryRes struct {
			Status              int8   `gorm:"status"`
			EstimatedDelivery   string `gorm:"estimated_delivery"`
			RecipientID         uint32 `gorm:"recipient_id"`
			RecipientName       string `gorm:"recipient_name"`
			RecipientAddress    string `gorm:"recipient_address"`
			RecipientPhone      string `gorm:"recipient_phone"`
			LocationID          uint32 `gorm:"location_id"`
			LocationTitle       string `gorm:"location_title"`
			LocationCity        string `gorm:"location_city"`
			LocationAddress     string `gorm:"location_address"`
			DetailID            uint32 `gorm:"detail_id"`
			DetailDate          string `gorm:"detail_date"`
			DetailTime          string `gorm:"detail_time"`
			DetailStatus        int8   `gorm:"detail_status"`
			DetailLocationID    uint32 `gorm:"detail_location_id"`
			DetailLocationTitle string `gorm:"detail_location_title"`
		}

		var dbRes []QueryRes

		if err := s.db.Table("tracking_statuses").
			Select("tracking_statuses.status as status, "+
				"tracking_statuses.estimated_delivery as estimated_delivery, "+
				"recipients.id as recipient_id, "+
				"recipients.name as recipient_name, "+
				"recipients.address as recipient_address, "+
				"recipients.phone as recipient_phone, "+
				"locations.location_id as location_id, "+
				"locations.title as location_title, "+
				"locations.city as location_city, "+
				"locations.address as location_address, "+
				"details.id as detail_id, "+
				"details.date as detail_date, "+
				"details.time as detail_time, "+
				"details.status as detail_status, "+
				"details.location_id as detail_location_id, "+
				"details.location_title as detail_location_title").
			Joins("inner join details ON tracking_statuses.sno=details.sno").
			Joins("inner join locations ON tracking_statuses.current_location_id=locations.location_id").
			Joins("inner join recipients ON tracking_statuses.recipient=recipients.id").
			Where("details.sno = ?", sno).
			Find(&dbRes).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, &QueryResponse{
				Status: "error",
				Data:   nil,
				Error: &QueryError{
					Code:    http.StatusInternalServerError,
					Message: "query failed",
				},
			})
		}

		details := []QueryDetail{}
		for _, res := range dbRes {
			detail := QueryDetail{
				ID:            res.DetailID,
				Date:          res.DetailDate,
				Hr:            res.DetailTime,
				Status:        po.StatusMsgMapping[po.Status(res.DetailStatus)],
				LocationID:    res.DetailLocationID,
				LocationTitle: res.DetailLocationTitle,
			}
			details = append(details, detail)
			data.Sno = uint32(sno)
			data.TrackingStatus = po.StatusMsgMapping[po.Status(res.Status)]
			data.EstimatedDelivery = res.EstimatedDelivery
			data.Recipient = QueryRecipient{
				ID:      res.RecipientID,
				Name:    res.RecipientName,
				Address: res.RecipientAddress,
				Phone:   res.RecipientPhone,
			}
			data.CurrentLocation = QueryCurrentLocation{
				LocationID: res.LocationID,
				Title:      res.LocationTitle,
				City:       res.LocationCity,
				Address:    res.LocationAddress,
			}
		}
		data.Details = details

		detailsJson, err := json.Marshal(data.Details)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &QueryResponse{
				Status: "error",
				Data:   nil,
				Error: &QueryError{
					Code:    http.StatusInternalServerError,
					Message: "query failed",
				},
			})
		}
		recipientJson, err := json.Marshal(data.Recipient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &QueryResponse{
				Status: "error",
				Data:   nil,
				Error: &QueryError{
					Code:    http.StatusInternalServerError,
					Message: "query failed",
				},
			})
		}
		currentLocationJson, err := json.Marshal(data.CurrentLocation)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &QueryResponse{
				Status: "error",
				Data:   nil,
				Error: &QueryError{
					Code:    http.StatusInternalServerError,
					Message: "query failed",
				},
			})
		}

		if err := s.cacheStore.HSetDataPairs(c.Request().Context(), "query", uint32(sno), map[string]interface{}{
			"tracking_status":    data.TrackingStatus,
			"estimated_delivery": data.EstimatedDelivery,
			"details":            string(detailsJson),
			"recipient":          string(recipientJson),
			"current_location":   string(currentLocationJson),
		}, time.Hour); err != nil {
			return c.JSON(http.StatusInternalServerError, &QueryResponse{
				Status: "error",
				Data:   nil,
				Error: &QueryError{
					Code:    http.StatusInternalServerError,
					Message: "query failed",
				},
			})
		}
	} else {
		data = assembleQueryData(cacheData)
	}

	response := &QueryResponse{
		Status: "success",
		Data:   data,
		Error:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

func assembleQueryData(cacheData map[string]string) *QueryData {
	data := &QueryData{}
	for field, val := range cacheData {
		switch field {
		case "tracking_status":
			data.TrackingStatus = val
		case "estimated_delivery":
			data.EstimatedDelivery = val
		case "details":
			var details []QueryDetail
			json.Unmarshal([]byte(val), &details)
			data.Details = details
		case "recipient":
			var recipient QueryRecipient
			json.Unmarshal([]byte(val), &recipient)
			data.Recipient = recipient
		case "current_location":
			var currentLocation QueryCurrentLocation
			json.Unmarshal([]byte(val), &currentLocation)
			data.CurrentLocation = currentLocation
		}
	}

	return data
}

func (s *Server) fakeHandler(c echo.Context) error {
	numStr := c.QueryParam("num")
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

	allDetails := []po.Detail{}
	statusList := []po.TrackingStatus{}
	for i := 0; i < num; i++ {
		sno := uuid.New().ID()
		allDetails = append(allDetails, genDetails(sno, locations)...)
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
	if err := tx.CreateInBatches(allDetails, len(allDetails)).Error; err != nil {
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
	count := rand.Intn(4) + 1
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
