.PHONY: warm-up build deploy status unload run infra-up infra-down

WARMUP_CSV ?= class.csv
APP_PATH ?= ./main.go

infra-up:
	docker-compose up -d redis mongo

infra-down:
	docker-compose down

warm-up:
	go run ./cmd/warmup --csv=$(WARMUP_CSV)

unload:
	go run ./cmd/warmup --unload

run:
	go run ./main.go

build:
	go build -o paralisus .