GOPATH ?= $(HOME)/go

export GOPATH

DOCKER_IMAGE ?= lilchomper/kafka-test
DOCKER_TAG ?= dev

all: kafka-test-image

kafka-test-image:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

build/kafka-test: main.go
	mkdir -p build
	go build -o build/kafka-test main.go

clean:
	rm -rf build

local:
	rm -f build/kafka-test
	$(MAKE) build/kafka-test

.PHONY: all local clean kafka-test kafka-test-image
