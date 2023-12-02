package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"logistics_status_tracking/api"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Setup
	viper.SetConfigName("api")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.AutomaticEnv()
	server := api.New(NewDB())

	// Start server
	go func() {
		if err := server.Start(":5000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.Logger.Fatal("shutting down the server", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
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
