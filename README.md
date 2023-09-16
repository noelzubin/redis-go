# Redis in go

### To run server
``` sh
go run server/server.go
``` 
starts server on port 6379.

### To run client
``` sh
go run client/client.go localhost:6379
```
Connects to default server

### Build 
``` sh
make build
```
binaries are generated in the `./bin` folder

### Commands implemented
```
PING
GET <key>
SET <key> <value> [EX_seconds]
DEL <key> [...<key>]
Expire <key> <EX_seconds>
Keys
ZAdd <setName> [<score> <value>] [...]
ZRange <setName> <start> <stope> [WITHSCORES]
```

# TODO
[ ] Add cron to delete expired keys
[ ] Allow sorted set to use float scores. Currently uses u64 
[ ] Benchmarking tests
