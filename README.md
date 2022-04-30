# Mercury

Mercury is a distributed key-value store implementing consistent hashing with a ring structure and fault tolerance.

## Run locally

```bash
go build -v -o Mercury .

./Mercury 8081

./Mercury 8082

./Mercury 8083

./Mercury 8084

./Mercury 8085

go run Client/Client.go
```
