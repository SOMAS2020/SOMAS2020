package server

import "github.com/SOMAS2020/SOMAS2020/internal/server/roles"

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	defer s.logf("finish runIIGO")
	// TODO:- IIGO team
	return roles.RunIIGO(&s.gameState)
}
