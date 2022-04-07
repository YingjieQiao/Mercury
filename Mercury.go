package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

func main() {
	//if os.Getenv("PORT") == "" {
	//	os.Setenv("PORT", "1234")
	//}
	//port, _ := strconv.Atoi(os.Getenv("PORT"))

	args := os.Args[1:]
	port, _ := strconv.ParseInt(args[0], 10, 64)

	// create server
	server := CreateServer(uint64(port))
	rpc.Register(server)
	rpc.HandleHTTP()

	// initial nodes discovery
	server.DiscoverNodes()

	l, e := net.Listen("tcp", ":"+args[0])
	log.Printf("Listening %s \n", args[0])
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	http.Serve(l, nil)
}
