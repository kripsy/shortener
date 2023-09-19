
.PHONY: build
build:
	go build -o ./bin/shortener ./cmd/shortener/*.go
	chmod +x ./bin/shortener

.PHONY: build_linter
build_linter:
	go build -o ./bin/staticlint ./cmd/staticlint/*.go
	chmod +x ./bin/staticlint

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


.PHONY: run_autotest
run_autotest: build
	~/go/shortenertest-darwin-arm64 -test.v -test.run=TestIteration1 -binary-path=./cmd/shortener/shortener -source-path=./

.PHONY: run_unittest
run_unittest:
	go test ./... -v -coverprofile=coverage.out; go tool cover -html=coverage.out


	