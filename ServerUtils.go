package main



// get value from local DB
func (s *Server) GetReplicatesValue(key *string, getReplicaResp *GetReplicaResp) error {
	value, _ := s.Database.Get(*key)
	getReplicaResp.Value = value
	return nil
}

// put value to local DB
func (s *Server) PutReplica(pushEvent *PushEvent, putReplicaResp *PutReplicaResp) error {
	s.Database.Put(pushEvent.Key, pushEvent.Value)
	putReplicaResp.Success = true
	return nil
}


