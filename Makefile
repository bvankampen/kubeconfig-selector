
COMMIT_ID=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)

NAME=cluster

all: clean build

run:
	@go run cmd/cli/main.go

clean:
	@echo ">> Cleaning..."
	@rm -rf bin

build: clean
	@echo ">> Building..."
	@echo "   Commit: $(COMMIT_ID)"
	@echo "   Version: $(VERSION)"
	@mkdir bin
	@go build -o bin/$(NAME) -ldflags "-X main.Version=$(VERSION) -X main.CommitId=$(COMMIT_ID)" ./cmd/...

install: clean build
	@echo ">> Installing $(NAME) in $(GOPATH)/bin..."
	@cp bin/$(NAME) $(GOPATH)/bin

.PHONY: all clean build install run
