# Makefile mínimo
# uso: make release 1.0.1

VERSION := $(word 2,$(MAKECMDGOALS))
PREFIX  ?= v
REMOTE  ?= origin

# Build info para ldflags
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
VERSION_INFO := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go build flags
LDFLAGS := -s -w -X main.version=$(VERSION_INFO) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

.PHONY: build
build:
	@CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o harborctl ./cmd/harborctl

.PHONY: release
release:
	@test -n "$(VERSION)" || (echo "uso: make release <versao>"; exit 1)
	git tag -a "$(PREFIX)$(VERSION)" -m "Release $(PREFIX)$(VERSION)"
	git push "$(REMOTE)" "$(PREFIX)$(VERSION)"

# ignora o argumento extra (ex.: 1.0.1) como alvo
%:
	@:
