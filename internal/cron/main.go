package main

import (
	"fmt"
	"logistics_status_tracking/internal/infra/po"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	loadConfig()
	db := NewDB()

	type Res struct {
		Total  int
		Status int8
	}
	var results []Res
	if err := db.Table("tracking_statuses").
		Select("COUNT(*) as total, status").
		Group("status").
		Scan(&results).Error; err != nil {
		return
	}

	currentTimeStr := time.Now().Format("2006-01-02T15:04:05Z")

	type Report struct {
		CreatedAt       string         `json:"created_at"`
		TrackingSummary map[string]int `json:"tracking_summary"`
	}
	report := Report{
		CreatedAt: currentTimeStr,
	}
	trackingSummary := make(map[string]int)

	for _, res := range results {
		statusStr := po.StatusMsgMapping[po.Status(res.Status)]
		trackingSummary[statusStr] = res.Total
	}
	report.TrackingSummary = trackingSummary

	fmt.Fprintln(os.Stdout, report)
}

func loadConfig() {
	// Setup
	viper.SetConfigName("api")
	viper.SetConfigType("env")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.AutomaticEnv()
}

func NewDB() *gorm.DB {
	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/tracking_status_storage?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
