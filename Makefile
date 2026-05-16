.PHONY: test lint tidy vet

test:
	go test ./... -race -v

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

vet:
	go vet ./...
