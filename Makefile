GOCMD=go1.12.5
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
TARGET=websrv
DEPCMD=dep

all: deps test build

.PHONY: dep deps-test deps test benchmark run image build

dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

deps-test: dep
	$(DEPCMD) ensure -v

deps: dep
	$(DEPCMD) ensure

test: deps
	$(GOTEST) -v -cover ./...

benchmark: deps
	$(GOTEST) -v -bench=. ./...

build: deps
	$(GOBUILD) -o $(TARGET) .

install: build
	cp $(TARGET) $(GOPATH)/bin/

clean:
	rm $(TARGET)
	$(GOCLEAN)

uninstall:
	if [ -a $(GOPATH)/bin/$(TARGET) ]; then rm $(GOPATH)/bin/$(TARGET); fi;

image:
	docker build --rm -t antonyho-websrv-demo .

run: deps
	$(GORUN) cmd/main.go -p 8080