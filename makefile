.PHONY: runsrv
runsrv:
	go run cmd/lesson/chatsrv/main.go

.PHONY: runcl
runcl:
	go run cmd/lesson/chatcli/main.go


.PHONY: build
build:
	go build cmd/lesson/chatcli/main.go
	go build cmd/lesson/chatsrv/main.go


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -race ./...