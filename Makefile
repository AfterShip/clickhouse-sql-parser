
PROGRAM=clickhouse-sql-parser
PKG_FILES=`go list ./... | sed -e 's=github.com/AfterShip/clickhouse-sql-parser/=./='`

CCCOLOR="\033[37;1m"
MAKECOLOR="\033[32;1m"
ENDCOLOR="\033[0m"

all: $(PROGRAM)

.PHONY: all

$(PROGRAM):
	go build -o $(PROGRAM) main.go

test:
	@go test -v ./... -covermode=atomic -coverprofile=coverage.out -race -compatible

update_test:
	@go test -v ./... -update -race -compatible

lint:
	@printf $(CCCOLOR)"GolangCI Lint...\n"$(ENDCOLOR)
	@golangci-lint run --timeout 20m0s
