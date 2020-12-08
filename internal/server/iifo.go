package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/action"

// runIIFO : IIFO makes recommendations about the optimal (and fairest) contributions this term.
// to mitigate common risk dilemma.
func (s *SOMASServer) runIIFO() ([]action.Action, error) {
	s.logf("start runIIFO")
	defer s.logf("finish runIIFO")
	// TODO:- IIFO team
	return nil, nil
}
