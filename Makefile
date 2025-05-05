.PHONY: test
test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

.PHONY: check
check:
	grep -L "Copyright" $$(find . -type f -name "*.go")
	@echo
	golangci-lint run
	@echo
	gofumpt -l .
