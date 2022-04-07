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
	var port int
	args := os.Args[1:]

	if os.Getenv("DOCKER") == "true" {
		if os.Getenv("PORT") == "" {
			os.Setenv("PORT", "1234")
		}
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	} else {
		port, _ = strconv.Atoi(args[0])
	}

	// create server
	server := CreateServer(uint64(port))
	rpc.Register(server)
	rpc.HandleHTTP()

	// initial nodes discovery
	server.DiscoverNodes()

	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	log.Printf("Listening %s \n", strconv.Itoa(port))
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	http.Serve(l, nil)
}
