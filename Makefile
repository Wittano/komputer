DEST_DIR = /opt/komputer
ARCH = $(shell uname -m)
OUTPUT_DIR=./build

ifeq ($(ARCH), x86_64)
	GOARCH="amd64"
else
  	GOARCH="arm64"
endif

dev:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -tags dev -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go

prod:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/komputer ./cmd/komputer/main.go

api:
	GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/api ./cmd/web/main.go

test-bot:
	go test -race ./bot/...

test-web:
	go test -race ./web/...

test: test-bot test-web

install: prod
	mkdir -p $(DEST_DIR)
	cp -r assets $(DEST_DIR)
	cp $(OUTPUT_DIR)/komputer $(DEST_DIR)

uninstall:
ifneq ("$(wildcard $(DEST_DIR))", "")
	rm -r $(DEST_DIR)
endif

clean:
ifneq ("$(wildcard $(OUTPUT_DIR))", "")
	rm -r build
endif