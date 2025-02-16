GO_BUILD_PATH ?= $(CURDIR)/bin
GO_BUILD_APP_PATH ?= $(GO_BUILD_PATH)/avito/

# Цели для кросс-компиляции
GOOS ?= linux
GOARCH ?= amd64
CGO ?= 0

# Цель для сборки
build:
	go env -w GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO)
	go build -o $(GO_BUILD_APP_PATH) ./cmd/

# Запуск docker compose
up:
	docker-compose up --build -d

# Остановка docker compose
down:
	@docker-compose down

# Очистка сгенерированных файлов
clean:
	@rm -rf $(GO_BUILD_PATH)

test-docker:
	docker-compose -f docker-compose.test.yml up -d
	go env -w GOOS=windows
	go test ./... -v
	docker-compose -f docker-compose.test.yml down -v