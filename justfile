### Variables

TARGETOS := `go env GOOS`
TARGETARCH := `go env GOARCH`
BIN_SUFFIX := if TARGETOS == "windows" { ".exe" } else { "" }
DIST_DIR := "dist"
CI_COMMIT_SHA := env_var_or_default("CI_COMMIT_SHA", `git rev-parse HEAD 2>/dev/null || echo "unknown"`)

# Set version based on CI environment
VERSION := if env_var_or_default("CI_COMMIT_TAG", "") != "" {
    replace(env_var("CI_COMMIT_TAG"), "v", "")
} else {
    "dev-" + `echo "${CI_COMMIT_SHA:-$(git rev-parse HEAD 2>/dev/null)}" | cut -c -10 2>/dev/null || echo "unknown"`
}

BUILD_DATE := `date -u +%Y-%m-%d`

LDFLAGS := "-s -w -extldflags '-static' -X codefloe.com/pat-s/backporter/shared/version.Version=" + VERSION + " -X codefloe.com/pat-s/backporter/shared/version.BuildDate=" + BUILD_DATE

### Recipes

## General

default:
    @just --list

version:
    @echo "{{VERSION}}"

fmt:
    gofumpt -w .
    golangci-lint run --fix
    gci write --skip-vendor --skip-generated -s standard -s default -s "prefix(codefloe.com/pat-s/backporter)" --custom-order .

generate:
    go generate ./...

install-dev-deps:
    @hash gofumpt > /dev/null 2>&1 || go install mvdan.cc/gofumpt@latest
    @hash golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    @hash gci > /dev/null 2>&1 || go install github.com/daixiang0/gci@latest
    @hash mockery > /dev/null 2>&1 || go install github.com/vektra/mockery/v2@latest

## Test

test: test-unit test-integration

test-unit:
    CGO_ENABLED=1 go test -race -cover -coverprofile coverage.out -timeout 60s -tags 'test' ./...

test-integration:
    CGO_ENABLED=1 go test -race -timeout 120s -tags 'integration test' ./...

## Lint

lint: install-dev-deps
    golangci-lint run

lint-fix: install-dev-deps
    golangci-lint run --fix

## Build

build:
    CGO_ENABLED=0 GOOS={{ TARGETOS }} GOARCH={{ TARGETARCH }} go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/backporter{{ BIN_SUFFIX }} ./cmd/backporter

build-dev:
    CGO_ENABLED=1 GOOS={{ TARGETOS }} GOARCH={{ TARGETARCH }} go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/backporter{{ BIN_SUFFIX }} ./cmd/backporter
    sudo mv ./dist/backporter /usr/local/bin/backporter
    sudo chmod +x /usr/local/bin/backporter

build-all: ## Create binaries for all platforms
    # Linux
    GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/linux_amd64/backporter       ./cmd/backporter
    GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/linux_arm64/backporter       ./cmd/backporter
    # macOS
    GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/darwin_amd64/backporter      ./cmd/backporter
    GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/darwin_arm64/backporter      ./cmd/backporter
    # Windows
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '{{ LDFLAGS }}' -o {{ DIST_DIR }}/windows_amd64/backporter.exe ./cmd/backporter
    # Create archives
    tar -czf {{ DIST_DIR }}/backporter_linux_amd64.tar.gz   -C {{ DIST_DIR }}/linux_amd64   backporter
    tar -czf {{ DIST_DIR }}/backporter_linux_arm64.tar.gz   -C {{ DIST_DIR }}/linux_arm64   backporter
    tar -czf {{ DIST_DIR }}/backporter_darwin_amd64.tar.gz  -C {{ DIST_DIR }}/darwin_amd64  backporter
    tar -czf {{ DIST_DIR }}/backporter_darwin_arm64.tar.gz  -C {{ DIST_DIR }}/darwin_arm64  backporter

clean:
    rm -rf {{ DIST_DIR }}

## Container images

image-alpine TAG REGISTRY='ghcr.io' IMAGE='pat-s/backporter':
    docker buildx build --platform linux/amd64,linux/arm64 \
        --build-arg VERSION={{ VERSION }} \
        -t {{ REGISTRY }}/{{ IMAGE }}:{{ TAG }} \
        -f images/Containerfile.alpine --push .

## Development

run *ARGS:
    go run ./cmd/backporter {{ ARGS }}

mockery:
    mockery
