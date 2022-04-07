package main





type Server struct {
	ID       int
	selfPort uint64
	Database Database
	nodeManager    *Manager
	numReplica int
	nodeSet []uint64
}
