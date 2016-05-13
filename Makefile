APP=e2u.io/amzimg3
PWD=$(shell pwd)
BUILD_DIR=$(PWD)/builds

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}

.PHONY: build
build:
	go build -o ${BUILD_DIR}/amzimg3


.PHONY: build-docker
build-docker:
	docker build -t ${REPOSITORY} .

.PHONY: run
run:
	go run main.go -allow etc/allow_sources.txt -data /tmp/001 -logfile /tmp/amzimg3.log
