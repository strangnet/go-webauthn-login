export CGO_ENABLED=0

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: docs build internal

dep:
	go get -u golang.org/x/lint/golint

lint:
	go vet ./...
	go list ./... | xargs -L1 golint -set_exit_status
	staticcheck -version || go get honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

test:
	go test -v ./...

race: dep ## Run data race detector
	CGO_ENABLED=1 go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	./build/ci/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./build/ci/coverage.sh html;
