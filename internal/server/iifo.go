package server

import "github.com/SOMAS2020/SOMAS2020/internal/common"

// runIIFO : IIFO makes recommendations about the optimal (and fairest) contributions this term.
// to mitigate common risk dilemma.
func (s *SOMASServer) runIIFO() ([]common.Action, error) {
	s.logf("start runIIFO")
	defer s.logf("finish runIIFO")
	// TODO:- IIFO team
	return nil, nil
}
