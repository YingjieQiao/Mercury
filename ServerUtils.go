package main

// GetReplicatesValue get value from local DB
func (s *Server) GetReplicatesValue(getEvent *GetEvent, getReplicaResp *GetReplicaResp) error {
	value, _ := s.Database.Get(getEvent.Port, *getEvent.Key)
	getReplicaResp.Value = value
	return nil
}

// PutReplica put value to local DB
func (s *Server) PutReplica(pushEvent *PushEventVC, putReplicaResp *PutReplicaResp) error {
	s.Database.Put(pushEvent.Port, pushEvent.Key, pushEvent.Value)
	putReplicaResp.Success = true
	return nil
}
