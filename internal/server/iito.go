package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/action"

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() ([]action.Action, error) {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	// TOOD:- IITO team
	return nil, nil
}
