DEST_DIR = /opt/komputer
ARCH = $(shell uname -m)

ifeq ($(ARCH), x86_64)
	GOARCH="amd64"
else
  	GOARCH="arm64"
endif

dev:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -tags dev -o ./build/komputer ./cmd/komputer/main.go

prod:
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -o ./build/komputer ./cmd/komputer/main.go

test:
	go test ./...

testContainers:
ifeq (,$(shell command -v docker 2> /dev/null))
	$(error "No docker in $(PATH)")
endif

	go test -tags testcontainers ./pkgs/db

install: prod
	mkdir -p $(DEST_DIR)
	cp -r assets $(DEST_DIR)
	cp build/komputer $(DEST_DIR)

uninstall:
	rm -r $(DEST_DIR)

clean:
	rm -r build