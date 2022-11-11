include .env

DB_CONFIG := configs/database.yml
# SODA_PATH := storage/migrations
MODEL_PATH := internal/model

soda-create:
	soda create -e ${SODA_ENV} -c ${DB_CONFIG}

soda-drop:
	soda drop -e ${SODA_ENV} -c ${DB_CONFIG}

# soda-generate-fizz:
# 	soda generate fizz ${FIZZ_NAME} -c ${DB_CONFIG} -e ${SODA_ENV} -p ${MIGRATION_PATH}
new-migrations:
	@echo "New directory if not exist: ${MIGRATION_PATH}"
	mkdir -p ${MIGRATION_PATH}

soda-generate-fizz:
	@read -p "fizz name? : " FIZZ_NAME \
	&& soda generate fizz $${FIZZ_NAME} -c ${DB_CONFIG} -e ${SODA_ENV} -p ${MIGRATION_PATH}

soda-generate-sql:
	@read -p "sql cmd name? : " SQLCMD_NAME \
	&& soda generate fizz $${SQLCMD_NAME} -c ${DB_CONFIG} -e ${SODA_ENV} -p ${MIGRATION_PATH}

soda-migrate-up:
	soda migrate up -c ${DB_CONFIG} -e ${SODA_ENV} -p ${MIGRATION_PATH}

soda-migrate-down:
	soda migrate down -c ${DB_CONFIG} -e ${SODA_ENV} -p ${MIGRATION_PATH}

swag-init:
	swag init --parseDependency --parseInternal --dir .

gen-configs:
	@echo generating sqlc.yaml
	@if [ -f ./configs/sqlc.yaml ]; then \
		rm ./configs/sqlc.yaml; \
	fi
	@if [ -f ./config_generator ]; then \
		rm ./config_generator; \
	fi
	@go build -o config_generator third_party/config_generator/*.go
	@./config_generator >> ./configs/sqlc.yaml
	@rm ./config_generator

sqlc-generate: gen-configs
	sqlc generate -f ./configs/sqlc.yaml

sqlc-clean:
	rm ./internal/models/*.sql.go
	rm ./internal/models/db.go
	rm ./internal/models/models.go

build:
	go build ./main.go

run:
	go run ./main.go

about: ## Display info related to the build
	@echo "- Protoc version  : $(shell protoc --version)"
	@echo "- Go version      : $(shell go version)"
	@echo "- Soda version    : ${shell soda --version}"