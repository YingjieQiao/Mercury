package main

import (

	"time"
	"strconv"
	"net/rpc"
	"log"
)

func CreateServer(port uint64) *Server {
	server := Server{ID: int(time.Now().UnixNano()), 
		Database: *NewDatabase(), 
		nodeManager: NewManager(port, 5, 5*time.Second),  //Peiyuan: Default: 5 virtual node per physical node. 5 seconds before expire
		selfPort: port, 
		numReplica: 3, //Peiyuan: Default: replicate on 5 nodes
		nodeSet: []uint64{8081, 8082,8083, 8084, 8085,}, //Peiyuan: defualt node map
	}
	
	return &server
}

func (s *Server) DiscoverNodes() {
	//add self
	s.nodeManager.UpdateNode(NodeInfo{
		Alive: true,
		Port: s.selfPort,  
		Version: time.Now().UnixNano(),
	})

	go func(){
		for {
			s.sendHeatBeart()
			time.Sleep(2 * time.Second)
		}
	}()

}

func (s *Server) sendHeatBeart() {

	nodeInfo := NodeInfo{
		Alive: true,
		Port: s.selfPort,  
		Version: time.Now().UnixNano(),
	}

	heartBeatMessage := HeartBeatMessage{Info: nodeInfo}
	for _, port := range s.nodeSet {
		reply := HeartBeatReply{}
		client, err := rpc.DialHTTP("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			log.Printf("Node %d fail dial to %d", s.selfPort, port)
			log.Print("sendHeatBeart error: ", err)
			continue
		}
		err = client.Call("Server.ReceiveHeatBeat", &heartBeatMessage, &reply)
		if err != nil {
			log.Printf("Node %d fail to send heartbeat to %d", s.selfPort, port)
			log.Print("sendHeatBeart error: ", err)
			continue
		}
		s.nodeManager.ImportRing(reply.Ring)
		log.Printf("Node %d succeed in sending heartbeat to %d", s.selfPort, port)
	}
}

func (s *Server) ReceiveHeatBeat(heartBeatMessage *HeartBeatMessage,heartBeatReply *HeartBeatReply) error {
	s.nodeManager.UpdateNode(heartBeatMessage.Info)
	heartBeatReply.Ring = s.nodeManager.ExportRing()
	return nil
}

// called by client, get value from preferencelists and aggregate 
func (s *Server) GetValue(key *string, clientGetResp *ClientGetResp) error {
	// hash key
	key_uint64, _ := strconv.ParseUint(*key, 10, 64)

	preferenceList, _ := s.nodeManager.GetPreferenceList(key_uint64, s.numReplica)
	reply_list := []string{}
	for _, port := range preferenceList{
		client, err := rpc.DialHTTP("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			log.Fatal("Dialing: ", err)
		}

		reply := GetReplicaResp{}

		err = client.Call("Server.GetReplicatesValue", key, &reply)
		if err != nil {
			log.Fatal("Server.GetValue error:", err)
		}

		reply_list = append(reply_list, reply.Value)
	}
	clientGetResp.Values = reply_list
	return nil
}

//called by client, push value to all nodes in preference list
func (s *Server) PushValue(pushEvent *PushEvent, clientPushResp *ClientPushResp) error {
	

	// hash key
	key_uint64, _ := strconv.ParseUint(pushEvent.Key, 10, 64)

	preferenceList, _ := s.nodeManager.GetPreferenceList(key_uint64, s.numReplica)
	log.Printf("pereference list for key %s is %v", pushEvent.Key, preferenceList)
	reply_list := []bool{}
	for _, port := range preferenceList{
		client, err := rpc.DialHTTP("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			log.Fatal("Dialing: ", err)
		}

		reply := PutReplicaResp{}

		err = client.Call("Server.PutReplica", pushEvent, &reply)
		if err != nil {
			log.Fatal("Server.GetValue error:", err)
		}
		reply_list = append(reply_list, reply.Success)
	}

	clientPushResp.Success = reply_list

	return nil
}
