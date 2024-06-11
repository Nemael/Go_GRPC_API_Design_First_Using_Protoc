generate_grpc_code:
	protoc \
	--go_out=Books \
	--go_opt=paths=source_relative \
	--go-grpc_out=Books \
	--go-grpc_opt=paths=source_relative \
	--plugin=/home/nemael/go/bin/protoc-gen-go \
	--plugin=/home/nemael/go/bin/protoc-gen-go-grpc \
	books.proto