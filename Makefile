APP=e2u.io/amzimg3
PWD=$(shell pwd)
BUILD_DIR=$(PWD)/builds

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}
	
.PHONY: build
build:
	go build -o ${BUILD_DIR}/amzimg3
	

