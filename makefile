.PHONY: run
run:
	go run cmd/lesson4/main.go

.PHONY: build
build:
	go build cmd/lesson4/main.go


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -race ./...