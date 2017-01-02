SRC_DIR = .
PROTOS := $(shell find $(SRC_DIR) -name '*.proto')

all: proto

proto:
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		deck/*.proto
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		*.proto
