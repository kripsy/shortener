
.PHONY: build
build:
	go build -o ./bin/shortener ./cmd/shortener/*.go
	chmod +x ./bin/shortener


.PHONY: db
db:
	docker-compose up postgres_db


.PHONY: migration_new
migration_new:
	docker-compose up migration_db_add

.PHONY: migration_up
migration_up:
	docker-compose up migration_db_up


.PHONY: migration_down
migration_down:
	docker-compose up migration_db_down


.PHONY: migration_force
migration_force:
	docker-compose up migration_db_force