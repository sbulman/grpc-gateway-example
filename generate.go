// +build tools

package tools

//go:generate protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=. --go_opt=module=github.com/sbulman --go-grpc_out=. --go-grpc_opt=module=github.com/sbulman todo_service.proto
