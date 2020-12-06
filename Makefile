GOCMD=go
GOBUILD=$(GOCMD) build ./cmd/main.go
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=simple_service
LINTER=golangci-lint

all: test build

test: $(GOTEST) ./... -v

build: $(GOBUILD) -o $(BINARY_NAME) -v

lint: $(LINTER) run
run:
	$(GORUN) ./cmd/main.go