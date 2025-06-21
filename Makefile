.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: bootstrap
bootstrap:
	cd ./deploy/docker-compose && docker compose up -d && cd ../../
	go run ./cmd/migration
	nunu run ./cmd/server

.PHONY: run
run:
	nunu run ./cmd/server

.PHONY: mock
mock:
	mockgen -source=internal/query/gen.go -destination test/mocks/query/gen.go
	mockgen -source=internal/query/shop.gen.go -destination test/mocks/query/shop.gen.go
	mockgen -source=internal/repository/repository.go -destination test/mocks/repository/repository.go
	mockgen -source=internal/service/shop.go -destination test/mocks/service/shop.go

.PHONY: test
test:
	go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...
	go tool cover -html=./coverage.out -o coverage.html

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: docker
docker:
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/task -t 1.1.1.1:5000/demo-task:v1 .
	docker run --rm -i 1.1.1.1:5000/demo-task:v1

.PHONY: swag
swag:
	swag init -g cmd/server/main.go -o ./docs --parseDependency