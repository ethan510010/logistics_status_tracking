.PHONY: all build clean test

all: build run

build:
	go build -o logistics_status_tracking ./cmd/main.go


run:
	./logistics_status_tracking
