package server

import "github.com/SOMAS2020/SOMAS2020/internal/common"

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() ([]common.Action, error) {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	// TODO:- IITO team
	return nil, nil
}
