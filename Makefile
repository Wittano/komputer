ARCH = $(shell uname -m)
OUTPUT_DIR=./build
PROTOBUF_API_DEST=./api

ifeq ($(ARCH), x86_64)
	GOARCH="amd64"
else
  	GOARCH="arm64"
endif

.PHONY: test clean

bot-dev:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -tags dev -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go


bot-prod:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go

test: test-bot

test-bot:
	CGO_CFLAGS="-w" go test ./bot/...;

all: bot-prod test
clean:
ifneq ("$(wildcard $(OUTPUT_DIR))", "")
	rm -r $(OUTPUT_DIR)
endif
