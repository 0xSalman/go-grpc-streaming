help:
	@echo "Usage: {options} make [target ...]"
	@echo
	@echo "Commands:"
	@echo "  install         Install required dependencies"
	@echo "  gen-rpc         Generate sources from proto files"
	@echo "  run-client      Start client"
	@echo "  run-server      Start server"
	@echo "  test            Run tests"
	@echo
	@echo "  help            Show available commands"
	@echo
	@echo "Examples:"
	@echo "  # Getting started"
	@echo "  make install gen-rpc"
	@echo "  make run-server run-client"
	@echo

install:
	@ echo "Download required dependencies"
	@ go get -u -v github.com/golang/protobuf/protoc-gen-go
	@ go build ./...
	@ echo "Finished downloading required dependencies"

gen-rpc:
	@ echo "Generating code from protos"
	@ protoc -I ./proto ./proto/number.proto --go_out=plugins=grpc:./proto
	@ echo "Finished generating code from protos"

run-server:
	@ echo "Starting server"
	@ go run server/server.go

run-client:
	@ echo "Starting client"
	@ go run client/client.go

test:
	@ echo "Running tests"
	@ go test -v ./...
	@ echo "Finished tests"


