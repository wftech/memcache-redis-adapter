CWD=$(shell pwd)
APP_USER=${USER}
CONTAINER_NAME=memcache-redis-adapter
BIN_NAME=memcache-redis-adapter

default: binary

binary: build clean
	docker run -v $(CWD):/app -it --rm $(CONTAINER_NAME) \
		go build -o $(BIN_NAME)

runshell: build
	docker run -v $(CWD):/app -it --rm $(CONTAINER_NAME) /bin/bash

clean:
	rm -f $(BIN_NAME)

upgrade-dependencies: build
	docker run -v $(CWD):/app -it --rm $(CONTAINER_NAME) go get -u

build:
	docker build --build-arg=APP_USER=$(APP_USER) -t $(CONTAINER_NAME) .

image:
	docker build -t $(BIN_NAME):latest -f Dockerfile.image .
