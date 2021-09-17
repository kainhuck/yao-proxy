GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

.PHONY: build-local build-remote build-cli build-all-platform

build-local:
	$(GO) build -o ./cmd/local/yp-local ./cmd/local

build-remote:
	$(GO) build -o ./cmd/remote/yp-remote ./cmd/remote

build-cli:
	$(GO) build -o ./cmd/cli/cli ./cmd/cli

build-all-platform:
	./bin/go_build_all.py -o ./output/cli cmd/cli/main.go
	./bin/go_build_all.py -o ./output/local cmd/local/main.go
	./bin/go_build_all.py -o ./output/remote cmd/remote/main.go

.PHONY: run-local run-remote run-cli run-test-local

run-cli:
	$(GO) run ./cmd/cli/main.go

run-local:
	$(GO) run ./cmd/local/main.go -c ./cmd/local/res/config.json

run-remote:
	$(GO) run ./cmd/remote/main.go -c ./cmd/remote/res/config.json

run-test-local:
	$(GO) run ./cmd/local/main.go -c /etc/yao/config.json

.PHONY: docker-local docker-remote

docker-local:
	docker build -t docker.pkg.github.com/kainhuck/yao-proxy/local -f cmd/local/Dockerfile .

docker-remote:
	docker build -t docker.pkg.github.com/kainhuck/yao-proxy/remote -f cmd/remote/Dockerfile .

# build and push docker
.PHONY: docker-all

docker-all:
	./bin/docker-build.sh all