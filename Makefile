PROJ_NAME=dm
BUILD_DIR=$(CURDIR)/.build

.PHONY: build
build:
	go build -o "$(BUILD_DIR)/$(PROJ_NAME)" main.go

.PHONY: clean
clean:
	rm -f "$(BUILD_DIR)"

.PHONY: install
install: build
	cp "$(BUILD_DIR)/$(PROJ_NAME)" $${GOPATH}/bin/

.PHONY: uninstall
uninstall:
	rm -f $${GOPATH}/bin/$(PROJ_NAME)
