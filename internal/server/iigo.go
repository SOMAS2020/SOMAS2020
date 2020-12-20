package server

import "github.com/SOMAS2020/SOMAS2020/internal/server/roles"

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	defer s.logf("finish runIIGO")
	// TODO:- IIGO team
	return roles.RunIIGO(&s.gameState)
}

func (s *SOMASServer) runIIGOEndOfTurn() error {
	s.logf("start runIIGOEndOfTurn")
	defer s.logf("finish runIIGOEndOfTurn")
	// TODO:- IIGO team
	return nil
}
