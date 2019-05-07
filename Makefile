default: test

test:
	go test ./...

bench:
	go test ./... -run=NONE -bench=. -benchmem

.PHONY: default test

# proto ---------------------------------------------------------------

proto: proto.go
proto.go: messages.pb.go

.PHONY: proto proto.go

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. --proto_path=.:$$GOPATH/src $<
