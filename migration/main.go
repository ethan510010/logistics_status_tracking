package main

import (
	"fmt"
	"logistics_status_tracking/internal/infra/po"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("root:Ethan0909@tcp(127.0.0.1:3306)/tracking_status_storage?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&po.Detail{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&po.Location{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&po.Recipient{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&po.TrackingStatus{}); err != nil {
		panic(err)
	}
}
