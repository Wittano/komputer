DEST_DIR = /opt/komputer
ARCH = $(shell uname -m)

ifeq ($(ARCH), x86_64)
	GOARCH="amd64"
else
  	GOARCH="arm64"
endif

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o ./build/komputer ./cmd/komputer/main.go

install: build
	mkdir -p $(DEST_DIR)
	cp -r assets $(DEST_DIR)
	cp build/komputer $(DEST_DIR)


uninstall:
	rm -r $(DEST_DIR)

clean:
	rm -r build