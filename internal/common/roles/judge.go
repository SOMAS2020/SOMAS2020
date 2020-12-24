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
	PayPresident() shared.Resources
	InspectHistory() (map[shared.ClientID]EvaluationReturn, error)
	DeclareSpeakerPerformance() (result bool, didRole bool)
	DeclarePresidentPerformance() (result bool, didRole bool)
}
