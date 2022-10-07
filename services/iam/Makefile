ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
MIGRATE_COMMAND:=docker run -v $(ROOT_DIR)/internal/database/postgresql/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://collectable_admin:collectable_admin_password@localhost:5432/iam?sslmode=disable
SQLC_COMMAND:=docker run -v $(ROOT_DIR):/src -w /src --network host kjconroy/sqlc

run:
	@go run main.go
build:
	@docker build . --target bin
gen-keys:
	@ openssl genpkey -algorithm RSA -out ./assets/jwt/private.pem -pkeyopt rsa_keygen_bits:2048
	@ openssl rsa -in ./assets/jwt/private.pem -pubout -out ./assets/jwt/public.pem
migrate-create:
	@ $(MIGRATE_COMMAND) create -ext=sql -dir=/migrations $(name)
migrate:
	@ $(MIGRATE_COMMAND) up
migrate-down:
	@ $(MIGRATE_COMMAND) down -all
sqlc-generate:
	@ $(SQLC_COMMAND) generate
