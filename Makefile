DEST_DIR = /opt/komputer
ARCH = $(shell uname -m)
OUTPUT_DIR=./build
PROTOBUF_API_DEST=./api

ifeq ($(ARCH), x86_64)
	GOARCH="amd64"
else
  	GOARCH="arm64"
endif

.PHONY: test clean

bot-dev: proto
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -tags dev -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go


bot-prod: proto
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go

sever: proto
	go build -o $(OUTPUT_DIR)/server ./cmd/server/main.go

tui: proto
	go build -o $(OUTPUT_DIR)/tui ./cmd/tui/main.go

protobuf: cleanProto
	mkdir -p $(PROTOBUF_API_DEST)
	protoc --go_out=. ./proto/*

test: proto
	go test -race ./...

all: bot-dev bot-prod sever tui

install: prod test
	mkdir -p $(PROTOBUF_API_DEST)
	cp -r assets $(PROTOBUF_API_DEST)
	cp $(OUTPUT_DIR)/komputer $(PROTOBUF_API_DEST)

uninstall:
ifneq ("$(wildcard $(PROTOBUF_API_DEST))", "")
	rm -r $(PROTOBUF_API_DEST)
endif

clean: cleanProto
ifneq ("$(wildcard $(OUTPUT_DIR))", "")
	rm -r $(OUTPUT_DIR)
endif

cleanProto:
ifneq ("$(wildcard $(PROTOBUF_API_DEST))", "")
	rm -r $(PROTOBUF_API_DEST)
endif
