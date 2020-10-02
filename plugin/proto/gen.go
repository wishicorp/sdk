//go:generate protoc -I ./ ./schema.proto ./consul.proto ./backend.proto ./storage.proto ./component.proto --go_out=plugins=grpc:.

package proto
