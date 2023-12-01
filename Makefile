.PHONY: all build clean test

all: build

build:
	go build -o logistics_status_tracking ./cmd/main.go
