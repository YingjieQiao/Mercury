package main
// Peiyuan: Please put every message struct in this file
type NodeInfo struct {
	Alive bool
	Port uint64  // Peiyuan: just use port number as the node id. easier to hash, as consistent.go only accept unit64 type
	Version int64
}

type PushEvent struct {
	Key   string 
	Value string 
}

//Peiyuan: sent by client to coordinator
type ClientGetResp struct{
	Values []string
}
type ClientPushResp struct{
	Success []bool
}


//Peiyuan: sent by coordinator to nodes in preferece list
type GetReplicaResp struct{
	Value string
}
type PutReplicaResp struct{
	Success bool
}

type HeartBeatReply struct{
	Ring map[uint64]NodeInfo 
}

type HeartBeatMessage struct{
	Info NodeInfo
}