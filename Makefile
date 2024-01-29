PSQL_USER ?= postgres
PSQL_PASSWORD ?= postgres
PSQL_DB ?= postgres
PSQL_PORT ?= 5432

NATS_PORT ?= 4222
NATS_MONITOR_PORT ?= 8222

include .env

export

run:
	docker-compose up -d


stop:
	docker-compose down