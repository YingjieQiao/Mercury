package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"
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

func roundRobin() []string {
	var nodeList []string
	envNodeList := os.Getenv("NODE_LIST")

	if envNodeList == "" {
		log.Println("Using default node list")
		nodeList := []string{":8081", ":8082", ":8083", ":8084", ":8085"}
		return nodeList
	} else {
		arr := strings.FieldsFunc(envNodeList, func(r rune) bool {
			return r == ','
		})
		for i, node := range arr {
			nodeList[i] = node
		}
	}

	return nodeList
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
			scalabilityTest()
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

func scalabilityTest() {
	var wg sync.WaitGroup
	n := 350 // 2*n write + read
	writeFails := 0
	readFails := 0
	start := time.Now()

	for i := 1; i <= n; i++ {
		nodeList := roundRobin()
		rand.Seed(time.Now().UnixNano())
		node := nodeList[rand.Intn(len(nodeList))]

		client, err := rpc.DialHTTP("tcp", node)
		if err != nil {
			log.Fatal("Dialing:", err)
		}
		log.Printf("Successfulling connected to node %s", node)
		args := PushEvent{"Bruce", "Banner"}
		wg.Add(1)
		wg.Add(1)

		go func() {
			defer wg.Done()

			writeReply := ClientPushResp{}
			err = client.Call("Server.PushValue", &args, &writeReply)
			if err != nil {
				writeFails++
			}
		}()

		go func() {
			defer wg.Done()

			readReply := ClientGetResp{}
			err = client.Call("Server.GetValue", "Bruce", &readReply)
			if err != nil {
				readFails++
			}
		}()
	}

	wg.Wait()

	end := time.Now()
	interval := end.Sub(start).Seconds()
	log.Printf("%d write fails for %d operations, %d write fails for for %d operations. \n",
		writeFails, n, readFails, n)
	log.Printf("Throughout: %f RPS", float64(n*2)/float64(interval))
}

func defaultTest() {
	//nodeList := roundRobin()
	//rand.Seed(time.Now().Unix())
	//node := nodeList[rand.Intn(len(nodeList))]
	//
	//client, err := rpc.DialHTTP("tcp", node)
	//if err != nil {
	//	log.Fatal("Dialing:", err)
	//}
	//log.Printf("Successfulling connected to node %s", node)

	client, err := rpc.DialHTTP("tcp", ":8081")
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
