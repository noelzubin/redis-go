gen_mocks:
	mockery --dir=store --name Store --output store/mocks
	mockery --dir=set --name IStringSet --output set/mocks

build:
	go build -o bin/server server/server.go
	go build -o bin/client client/client.go