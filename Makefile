PROJ_NAME=domic
BUILD_DIR=$(CURDIR)/.build

.PHONY: build
build:
	go build -o "$(BUILD_DIR)/$(PROJ_NAME)" main.go

.PHONY: clean
clean:
	rm -f "$(BUILD_DIR)"

.PHONY: install
install:
	# go install github.com/cqroot/$(PROJ_NAME)@latest
	go install $(CURDIR)

.PHONY: uninstall
uninstall:
	rm -f $${GOPATH}/bin/$(PROJ_NAME)

.PHONY: test
test:
	go test -v ./...

.PHONY: docs
docs:
	go run $(CURDIR)/docs/apps/main.go > $(CURDIR)/docs/apps.md

.PHONY: check
check:
	golangci-lint run
	@echo
	gofumpt -l .
