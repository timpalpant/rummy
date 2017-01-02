SRC_DIR = .

all: proto server clients

proto:
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		deck/*.proto
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:${SRC_DIR} \
		*.proto
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		*.proto

server:
	go build rummy/gameserver/gamed

clients:
	go build rummy/clients/cli
