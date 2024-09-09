include .env
url="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
migration_dir="./migration"

status:
	GOOSE_MIGRATION_DIR=${migration_dir} goose postgres ${url} status

up:
	GOOSE_MIGRATION_DIR=${migration_dir} goose postgres ${url} up

down:
	GOOSE_MIGRATION_DIR=${migration_dir} goose postgres ${url} down

