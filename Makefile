gen:
	protoc -I/usr/local/include -I. \
      --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=:. --go-grpc_opt=paths=source_relative \
    proto/api.proto