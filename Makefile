.PHONY: docs build build-cli test test-v test-integration cover cover-check lint fmt vet tidy snapshot install-tools

MODULE := github.com/Edge-Center/edgecentercdn-go
COVER_PKGS := $(shell go list ./... | grep -v '/cmd')
COVER_MIN := 70

build:
	go build ./...

build-cli:
	go build -o bin/edge-cli ./cmd

test:
	go test -race ./...

test-v:
	go test -race -v ./...

test-integration:
	go test -tags=integration -race -count=1 ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

cover-check:
	go test -race -coverprofile=coverage.out $(COVER_PKGS)
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}' | \
		awk -F. '{if ($$1 < $(COVER_MIN)) {print "Coverage " $$0 " is below $(COVER_MIN)%"; exit 1}}'

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w -local $(MODULE) .

vet:
	go vet ./...

tidy:
	go mod tidy

snapshot:
	goreleaser release --snapshot --clean

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/goreleaser/goreleaser/v2@latest

docs:
	go list ./... | xargs -n1 go doc -all
