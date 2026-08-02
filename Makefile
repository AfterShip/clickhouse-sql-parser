
PROGRAM=clickhouse-sql-parser
PKG_FILES=`go list ./... | sed -e 's=github.com/AfterShip/clickhouse-sql-parser/=./='`

CCCOLOR="\033[37;1m"
MAKECOLOR="\033[32;1m"
ENDCOLOR="\033[0m"

all: $(PROGRAM)

.PHONY: all

VERSION=$(shell git describe --tags --always --dirty 2>/dev/null)

$(PROGRAM):
	go build -ldflags "-X main.version=$(VERSION)" -o $(PROGRAM) .

test:
	@go test -v ./... -covermode=atomic -coverprofile=coverage.out -race -compatible

update_test:
	@go test -v ./... -update -race -compatible

lint:
	@printf $(CCCOLOR)"GolangCI Lint...\n"$(ENDCOLOR)
	@golangci-lint run --timeout 20m0s
