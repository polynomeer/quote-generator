APP=quote-generator
BIN=bin/$(APP)
PKG=github.com/your-org/quote-generator

.PHONY: all build run test clean docker

all: build

build:
	@mkdir -p bin
	GO111MODULE=on go build -o $(BIN) ./cmd/quote-generator

run:
	QG_CONFIG=configs/config.yaml $(BIN)

test:
	go test ./...

clean:
	rm -rf bin

docker:
	docker build -t $(APP):dev .