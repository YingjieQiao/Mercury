package main

import (

	"fmt"
	"log"
	"net/rpc"
)

type PushEvent struct {
	Key   string 
	Value string 
}
type ClientPushResp struct{
	Success []bool
}

type ClientGetResp struct{
	Values []string
}

func main() {
	client, err := rpc.DialHTTP("tcp", ":8081") // connect to the node
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	reply  :=  ClientPushResp{}
	reply2 := ClientGetResp{}

	args := PushEvent{"Hello", "There"}


	err = client.Call("Server.PushValue", &args, &reply)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	fmt.Printf("%v\n", reply) // should be true, pushed successfully

	err = client.Call("Server.GetValue", "Hello", &reply2)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	fmt.Printf("Hello %v\n", reply2) // should be "There" as value was pushed successfully

}
