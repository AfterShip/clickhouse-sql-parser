
PROGRAM=clickhouse-sql-parser
PKG_FILES=`go list ./... | sed -e 's=github.com/AfterShip/clickhouse-sql-parser/=./='`

CCCOLOR="\033[37;1m"
MAKECOLOR="\033[32;1m"
ENDCOLOR="\033[0m"

all: $(PROGRAM)

# The binary target is phony so a new commit, tag or VERSION override
# always refreshes the stamped version; the Go build cache keeps the
# rebuild cheap.
.PHONY: all $(PROGRAM)

# VERSION is stamped into the binary. Priority: a caller-supplied
# VERSION=..., git describe in a checkout, then the git-archive
# substitution kept in .version for builds from source archives that
# have no Git metadata.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null)
ifeq ($(VERSION),)
# Reject the raw placeholder that an archive made by git older than
# 2.35 leaves behind; the binary then falls back to Go build info.
VERSION := $(shell grep -vE '[$$%]' .version 2>/dev/null)
endif

$(PROGRAM):
	go build -ldflags "-X main.version=$(VERSION)" -o $(PROGRAM) .

test:
	@go test -v ./... -covermode=atomic -coverprofile=coverage.out -race -compatible

update_test:
	@go test -v ./... -update -race -compatible

lint:
	@printf $(CCCOLOR)"GolangCI Lint...\n"$(ENDCOLOR)
	@golangci-lint run --timeout 20m0s
