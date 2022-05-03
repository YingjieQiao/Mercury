go build -v -o Mercury .

./Mercury 8081 &
./Mercury 8082 &
./Mercury 8083 &
./Mercury 8084 &
./Mercury 8085 &

go run Client/Client.go 3

echo "killing some nodes ... "

kill -9 $(lsof -ti:8082)
kill -9 $(lsof -ti:8083)
kill -9 $(lsof -ti:8084)
kill -9 $(lsof -ti:8085)

echo "send requests again after some nodes are killed"

go run Client/Client.go 3

echo "sleep 5 seconds ... "

sleep 5

echo "reviving the nodes ... "

./Mercury 8082 &
./Mercury 8083 &
./Mercury 8084 &
./Mercury 8085 &

echo "send requests again after nodes are revived"

go run Client/Client.go 3

echo "killing all nodes and exit "

kill -9 $(lsof -ti:8081)
kill -9 $(lsof -ti:8082)
kill -9 $(lsof -ti:8083)
kill -9 $(lsof -ti:8084)
kill -9 $(lsof -ti:8085)