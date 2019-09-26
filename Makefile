# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
GOLINTER=golangci-lint
BINARY_NAME=go-calendar

all: lint test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

lint:
	$(GOLINTER) run -v

test:
	$(GOTEST) -v ./...

generate:
	$(GOGENERATE) -x ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
