GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
TARGET=websrv

all: test build

.PHONY: test benchmark run image build

test:
	$(GOTEST) -v -race -cover ./...

benchmark:
	$(GOTEST) -race -bench=. ./...

build:
	$(GOBUILD) -o $(TARGET) ./cmd

install: build
	cp $(TARGET) $(GOPATH)/bin/

clean:
	rm $(TARGET)
	$(GOCLEAN)

uninstall:
	if [ -a $(GOPATH)/bin/$(TARGET) ]; then rm $(GOPATH)/bin/$(TARGET); fi;

image:
	docker build --rm -t antonyho-websrv-demo .

run:
	$(GORUN) cmd/main.go -p 8080