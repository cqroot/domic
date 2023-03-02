PROJ_NAME=gmdots
BUILD_DIR=$(CURDIR)/.build

.PHONY: build
build:
	go build -o "$(BUILD_DIR)/$(PROJ_NAME)" main.go

.PHONY: clean
clean:
	rm -f "$(BUILD_DIR)"

.PHONY: install
install:
	go install github.com/cqroot/$(PROJ_NAME)@latest

.PHONY: uninstall
uninstall:
	rm -f $${GOPATH}/bin/$(PROJ_NAME)

.PHONY: test
test:
	go test -v ./...
