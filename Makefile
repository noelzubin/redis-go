gen_mocks:
	mockery --dir=store --name Store

build:
	go build -o bin/server server/server.go
	go build -o bin/client client/client.go