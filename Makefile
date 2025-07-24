# Название бинарного файла
BIN_NAME=auth
VERSION?=0.1.0

# Go параметры
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
GOMOD=$(GO) mod

# Пути
SRC_DIR=./cmd
BUILD_DIR=./bin

.PHONY: all build clean test run fmt vet lint docker-build help

all: build

## build: Скомпилировать пример
build:
	$(GOBUILD) -C $(SRC_DIR) -o $(BUILD_DIR)/$(BIN_NAME)

## clean: Удалить скомпилированные файлы
clean:	
	$(GOCLEAN)
	rm -rf $(SRC_DIR)/bin
	
## test: Запустить тесты
test:
	$(GOTEST) ./... -v -cover -count=1

## fmt: Форматировать исходный код
fmt:
	$(GO) fmt ./...

## vet: Проверить код на наличие подозрительных конструкций
vet:
	$(GO) vet ./...

## lint: Запустить линтер (golangci-lint)
lint:
	golangci-lint run ./...

## docker-build: Собрать Docker-образ
docker-build:
	docker build -t $(BIN_NAME):$(VERSION) .

## help: Показать справку по командам
help:
	@echo "Доступные команды:"
	@echo
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo