package server

// vote for rule changes or election
func (s *SOMASServer) voting() (bool, error) {
	s.logf("start voting")
	defer s.logf("finish voting")
	// TODO:- voting team
	return false, nil
}
