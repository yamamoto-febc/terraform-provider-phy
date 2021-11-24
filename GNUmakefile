TEST?=$$(go list ./...)

default: build

build:
	go build

install:
	go install

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./internal/provider

generate: build
	go generate  ./...

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./internal/provider

test: fmtcheck
	go test -count=1 $(TESTARGS) -timeout=30s $(TEST)

# Run acceptance tests
.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test -count=1 $(TEST) -v $(TESTARGS) -timeout 120m
