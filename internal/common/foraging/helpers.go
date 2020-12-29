package foraging

import (
	"encoding/json"
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ForagingReport holds information about the result of a foraging session
type ForagingReport struct {
	InputResources     shared.Resources // combined input resources
	NumberParticipants uint             // number of participants in foraging session
	NumberCaught       uint             // number of deer/fish/... caught
	TotalUtility       shared.Resources // total return of foraging session before distribution
	CatchSizes         []float64        // sizes/weights of individual deer/fish/... caught
	Turn               uint             // turn in which this report was generated. Should be populated by caller
}

func getTotalInput(contribs map[shared.ClientID]shared.Resources) shared.Resources {
	i := shared.Resources(0.0)
	for _, x := range contribs {
		i += x
	}
	return i
}

func compileForagingReport(contribs map[shared.ClientID]shared.Resources, forageReturns []shared.Resources) ForagingReport {
	fR := ForagingReport{InputResources: getTotalInput(contribs), NumberParticipants: uint(len(contribs))}
	for _, r := range forageReturns {
		fR.TotalUtility += r
		if r > 0.0 {
			fR.CatchSizes = append(fR.CatchSizes, float64(r)) // store deer weights for post-hunt analysis
		}
	}
	fR.NumberCaught = uint(len(fR.CatchSizes))
	return fR
}

// utilityTier gets the discrete utility tier (i.e. max number of deer/fish) for given scalar resource input
func utilityTier(input shared.Resources, maxNumberPerHunt uint, decay float64) uint {
	inputF := float64(input)
	sum := 0.0
	for i := uint(0); i < maxNumberPerHunt; i++ {
		sum += math.Pow(decay, float64(i))
		if inputF < sum {
			return i
		}
	}
	return maxNumberPerHunt
}

// Display returns a JSON string of a foraging report
func (f ForagingReport) Display() string {
	out, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(out)
}

// Copy returns a deep copy of the ClientInfo.
func (f ForagingReport) Copy() ForagingReport {
	ret := f
	return ret
}
