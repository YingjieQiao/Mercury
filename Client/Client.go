package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type PushEvent struct {
	Key   string
	Value string
}

type ClientPushResp struct {
	Success []bool
}

type ClientGetResp struct {
	Values []string
}

func main() {
	for {
		fmt.Printf("Enter 1: Scalability Test\n" +
			"Enter 2: Vector Clock/Correctness Test\n" +
			"Enter any other key to run Default KV Test\n")

		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		if char == '1' {
			fmt.Println("Scalability Test")
			// TODO
			break
		} else if char == '2' {
			fmt.Println("Vector Clock/Correctness Test")
			// TODO
			break
		} else {
			fmt.Println("Default KV Test")
			defaultTest()
			break
		}
	}
}

func defaultTest() {
	client, err := rpc.DialHTTP("tcp", ":8081") // connect to the node
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	reply := ClientPushResp{}
	reply2 := ClientGetResp{}

	args := PushEvent{"Bruce", "Banner"}
	err = client.Call("Server.PushValue", &args, &reply)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	fmt.Printf("Push value response: %v\n", reply)

	err = client.Call("Server.GetValue", "Bruce", &reply2)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	fmt.Printf("Get value response: Bruce %v\n", reply2)
}
