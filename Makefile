run-book: 
	go run ./book/cmd/main.go

proto-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto-gen:
	protoc \
		-I=idls/proto/v1 --go_out=. --go-grpc_out=. \
    idls/proto/v1/*.proto
