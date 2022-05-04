# Mercury

Mercury is a distributed key-value store implementing consistent hashing with a ring structure and fault tolerance.

## Run locally

1: Scalability Test

2: Vector Clock/Correctness Test

3: Fault Tolerance Test

```bash
go build -v -o Mercury .

./Mercury 8081

./Mercury 8082

./Mercury 8083

./Mercury 8084

./Mercury 8085

go run Client/Client.go 3   
```

## Test

### Scalability Test

```bash
go run Client/Client.go 1
```

### Consistency Test

```bash
go run Client/Client.go 2
```

### Fault Tolerance Test

```bash
./Client/faultTolerenceTest.sh  
```

