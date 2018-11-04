.PHONY: proto
proto:
	@protoc --go_out=plugins=grpc:. ./api/pinger/pinger.proto

.PHONY: server
server:
	@go run ./cmd/server/server.go

.PHONY: client
client:
	@go run ./cmd/client/client.go

.PHONY: grpcurl
grpcurl:
	grpcurl -d '{"message":"hi"}' -plaintext :50051 pinger.Pinger.Ping