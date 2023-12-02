package handlers

import (
	"encoding/json"
	"logistics_status_tracking/internal/infra/po"
	"os"

	"gorm.io/gorm"
)

type FakeDataService struct {
	db *gorm.DB
}

func NewFakeDataService(db *gorm.DB) *FakeDataService {
	return &FakeDataService{
		db: db,
	}
}

func (f *FakeDataService) ParseLocation() ([]po.Location, error) {
	data, err := os.ReadFile("internal/service/handlers/locations.json")
	if err != nil {
		return nil, err
	}
	var locations []po.Location
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

func (f *FakeDataService) ParseRecipient() ([]po.Recipient, error) {
	data, err := os.ReadFile("internal/service/handlers/recipients.json")
	if err != nil {
		return nil, err
	}
	var recipients []po.Recipient
	if err := json.Unmarshal(data, &recipients); err != nil {
		return nil, err
	}
	return recipients, nil
}
