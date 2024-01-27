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
#	PSQL_USER=$(PSQL_USER) \
#	PSQL_PASSWORD=$(PSQL_PASSWORD) \
#	PSQL_DB=$(PSQL_DB) \
#    POSTGRES_PORT=$(POSTGRES_PORT) \
#    NATS_PORT=$(NATS_PORT) \
#    NATS_MONITOR_PORT=$(NATS_MONITOR_PORT) \

stop:
	docker-compose down