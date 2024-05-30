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


bot-prod: protobuf
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go

sever: protobuf
	go build -o $(OUTPUT_DIR)/server ./cmd/server/main.go

tui: protobuf
	go build -o $(OUTPUT_DIR)/tui ./cmd/tui/main.go

protobuf:
	mkdir -p $(PROTOBUF_API_DEST)
	protoc --go_out=./api --go_opt=paths=source_relative --go-grpc_out=./api --go-grpc_opt=paths=source_relative proto/*

test: protobuf
	find . -name go.mod -execdir go test ./... \;

all: bot-prod sever tui test

# FIXME Replaced invalid path
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
