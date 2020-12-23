package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.

type EvaluationReturn struct {
	Rules       []rules.RuleMatrix
	Evaluations []bool
}
type Judge interface {
	PayPresident() (shared.Resources, error)
	InspectHistory() (map[shared.ClientID]EvaluationReturn, error)
	DeclareSpeakerPerformance() (BID int, result bool, SID shared.ClientID, checkRole bool, err error)
	DeclarePresidentPerformance() (RID int, result bool, PID shared.ClientID, checkRole bool, err error)
}
