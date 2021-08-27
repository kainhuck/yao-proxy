GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

.PHONY: build-local build-remote

build-local:
	$(GO) build -o ./cmd/local/yp-local ./cmd/local

build-remote:
	$(GO) build -o ./cmd/remote/yp-remote ./cmd/remote

.PHONY: run-local run-remote

run-local:
	$(GO) run ./cmd/local/main.go -c ./cmd/local/res/config.json

run-remote:
	$(GO) run ./cmd/remote/main.go -c ./cmd/remote/res/config.json
