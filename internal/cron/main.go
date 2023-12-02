package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"logistics_status_tracking/internal/infra/po"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(report); err != nil {
		return
	}

	uploader := NewS3Uploader()

	bucket := viper.GetString("AWS_S3_TRACKING_BUCKET")
	filename := fmt.Sprintf("report-%d.json", time.Now().Unix())
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   buf,
	}); err != nil {
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
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

func NewS3Uploader() *s3manager.Uploader {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("AWS_ACCESS_KEY"), viper.GetString("AWS_SECRET"), ""),
	})
	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(sess)
	return uploader
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
