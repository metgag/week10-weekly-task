include ./.env
DBURL=$(DB_URL)?sslmode=disable
DBURLM=${DB_URL_M}?sslmode=disable
MIGRATIONPATH=db/migrations
SEED_PATH=db/seeds

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DBURLM) -path $(MIGRATIONPATH) up

migrate-down:
	migrate -database $(DBURLM) -path $(MIGRATIONPATH) down

insert-seed:
	for file in $$(ls $(SEED_PATH)/*.sql | sort); do \
		psql "$(DBURLM)" -f $$file; \
	done

swag-all:
	swag fmt
	swag init -d ./cmd
	swag init -g ./cmd/main.go