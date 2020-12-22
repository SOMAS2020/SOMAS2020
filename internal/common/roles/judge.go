package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.

type EvaluationReturn struct {
	Rules       []rules.RuleMatrix
	Evaluations []bool
}
type Judge interface {
	PayPresident() (int, error)
	InspectHistory() (map[int]EvaluationReturn, error)
	DeclareSpeakerPerformance() (int, bool, int, bool, error)
	DeclarePresidentPerformance() (int, bool, int, bool, error)
}
