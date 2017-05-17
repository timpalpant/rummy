SRC_DIR = .

all: proto server clients

proto:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:${SRC_DIR} \
		deck/*.proto
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:${SRC_DIR} \
		*.proto
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		*.proto

server:
	go build github.com/timpalpant/rummy/gameserver/gamed

clients:
	go build github.com/timpalpant/rummy/clients/cli
