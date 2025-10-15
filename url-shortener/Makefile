.PHONY: proto clean

proto:
	protoc --go_out=../ --go-grpc_out=../ api/grpc/entry.proto

clean:
	rm -f api/grpc/entry.pb.go api/grpc/entry_grpc.pb.go

run:
	go run .

lint:
	golangci-lint run
