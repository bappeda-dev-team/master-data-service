APP_NAME=master-data-service

.PHONY: build run clean test swagger build-image build-docker format dev check-format

build:
	@mkdir -p bin
	go build -o ./bin/$(APP_NAME) ./cmd/api

build-docker:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
	-o ./bin/app ./cmd/api

format:
	goimports -w .

check-format:
	@test -z "$$(gofmt -l .)" || (echo "gofmt required:" && gofmt -l . && exit 1)
	@test -z "$$(goimports -l .)" || (echo "goimports required:" && goimports -l . && exit 1)

run:
	go run ./cmd/api

dev: format run

test:
	go test -race ./...

swagger:
	swag init -d cmd/api,internal/api,internal/web,internal/master_sync/domain,internal/master_sync/web,internal/opd/domain,internal/opd/web,internal/jabatan/domain,internal/jabatan/web,internal/pegawai/domain,internal/pegawai/web

clean:
	rm -f ./bin/$(APP_NAME)

build-image:
	@docker build . -t $(APP_NAME)
