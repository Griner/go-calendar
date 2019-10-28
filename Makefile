# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
GOLINTER=golangci-lint
BINARY_NAME=go-calendar
BINARY_UNIX=$(BINARY_NAME)

all: lint test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

build_unix:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

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

docker: build_unix
	docker build -t go-calendar -f scripts/docker/go-calendar.Dockerfile .
	docker build -t go-calendar-server -f scripts/docker/go-calendar-server.Dockerfile .
	docker build -t go-calendar-mqworker -f scripts/docker/go-calendar-mqworker.Dockerfile .
	docker build -t go-calendar-notifier -f scripts/docker/go-calendar-notifier.Dockerfile .

