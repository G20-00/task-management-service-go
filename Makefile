# Formateo y linting automático para Go

.PHONY: format lint

format:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

test:
	go test ./...
