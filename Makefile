
COMMIT_ID=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat VERSION)

NAME=ks

all: clean build

debug:
	@go run cmd/cli/main.go --debug

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
	@mkdir -p $(GOPATH)/bin
	@cp bin/$(NAME) $(GOPATH)/bin

.PHONY: all clean build install run debug
