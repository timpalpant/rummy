SRC_DIR = .

all: proto server clients

proto:
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		deck/*.proto
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		*.proto

server:
	go build rummy/gameserver/gamed

clients:
	go build rummy/clients/cli
