package server

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (bool, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")
	// TOOD:- env team
	return false, nil
}
