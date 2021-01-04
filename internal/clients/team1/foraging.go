package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"

	"github.com/sajari/regression"
)

// ForageOutcome is how much we put in to a foraging trip and how much we got
// back
type ForageOutcome struct {
	turn         uint
	participant  shared.ClientID
	contribution shared.Resources
	revenue      shared.Resources
}

func (o ForageOutcome) ROI() float64 {
	if o.contribution == 0 {
		return 0
	}

	return float64((o.revenue / o.contribution) - 1)
}

func (o ForageOutcome) profit() shared.Resources {
	return o.revenue - o.contribution
}

type ForageHistory map[shared.ForageType][]ForageOutcome

// A forageDecider gives a ForageDecision as well as the expected reward if that decision is taken
type forageDecider = func(client) (shared.ForageDecision, shared.Resources)

func randomDecider(c client) (shared.ForageDecision, shared.Resources) {
	// Up to 10% of our current resources
	forageContribution := shared.Resources(0.2*rand.Float64()) * c.gameState().ClientInfo.Resources
	var forageType shared.ForageType
	if rand.Float64() < 0.5 {
		forageType = shared.DeerForageType
	} else {
		forageType = shared.FishForageType
	}

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}, shared.Resources(-1)
}

func roiDecider(c client) (shared.ForageDecision, shared.Resources) {
	forageType := bestROIForageType(c.forageHistory)
	if forageType == shared.ForageType(-1) {
		// Here we have found that the best idea is to not forage
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}, 0
	}
	expectedOutcome := bestROIOutcome(c.forageHistory[forageType])

	contribution := expectedOutcome.contribution
	// Cap to 20% of resources
	contribution = shared.Resources(math.Min(
		float64(contribution),
		c.config.forageContributionCapPercent*float64(c.gameState().ClientInfo.Resources),
	))
	// Add some noise
	contribution += shared.Resources(math.Min(
		rand.Float64(),
		c.config.forageContributionNoisePercent*float64(c.gameState().ClientInfo.Resources),
	))

	forageDecision := shared.ForageDecision{
		Type:         forageType,
		Contribution: contribution,
	}
	return forageDecision, expectedOutcome.revenue
}

func regressionDecider(c client) (shared.ForageDecision, shared.Resources) {
	bestReward := shared.Resources(0)
	var decision shared.ForageDecision

	for forageType := range c.forageHistory {
		r := outcomeRegression(c.forageHistory[forageType])
		contribution := c.regressionOptimalContribution(r)
		if contribution > c.gameState().ClientInfo.Resources {
			contribution = c.gameState().ClientInfo.Resources
		}
		expectedRewardF, _ := r.Predict([]float64{float64(contribution)})
		expectedReward := shared.Resources(expectedRewardF)

		if expectedReward > bestReward {
			bestReward = expectedReward
			decision = shared.ForageDecision{
				Type:         forageType,
				Contribution: contribution,
			}
		}
	}

	return decision, bestReward
}

const flipScale = 0.3

// flipDecider does the opposite of what the mass did the previous turn. It
// forages deer, with an amount inversely proportional to the sum of contributed
// resources
func flipDecider(c client) (shared.ForageDecision, shared.Resources) {
	deerHistory := c.forageHistory[shared.DeerForageType]
	totalContributionLastTurn := shared.Resources(0)
	totalHuntersLastTurn := 0
	totalRevenueLastTurn := shared.Resources(0)
	for _, outcome := range deerHistory {
		if outcome.turn == c.gameState().Turn-1 {
			totalContributionLastTurn += outcome.contribution
			totalRevenueLastTurn += outcome.revenue
			totalHuntersLastTurn += 1
		}
	}

	if totalContributionLastTurn == shared.Resources(0) {
		// Big contribution
		contribution := 0.3 * c.gameState().ClientInfo.Resources
		return shared.ForageDecision{
			Contribution: 0.3 * c.gameState().ClientInfo.Resources,
			Type:         shared.DeerForageType,
		}, 1.2 * contribution // arbitrarily expect a 20% ROI
	}

	// Proxy for population
	totalROI := totalRevenueLastTurn / totalContributionLastTurn
	averageContribution := totalContributionLastTurn / shared.Resources(totalHuntersLastTurn)

	contribution := flipScale * totalROI * averageContribution
	contribution = shared.Resources(math.Min(
		float64(0.2 * c.gameState().ClientInfo.Resources),
		float64(contribution),
	))

	return shared.ForageDecision{
		Contribution: contribution,
		Type:         shared.DeerForageType,
	}, contribution * totalROI
}

func outcomeRegression(history []ForageOutcome) regression.Regression {
	r := new(regression.Regression)
	r.SetObserved("Team1 reward")
	r.SetVar(0, "Team1 contribution")

	for _, outcome := range history {
		r.Train(regression.DataPoint(float64(outcome.revenue), []float64{float64(outcome.contribution)}))
	}

	r.AddCross(regression.PowCross(0, 2))
	r.Run()

	return *r
}

func (c *client) regressionOptimalContribution(r regression.Regression) shared.Resources {
	if r.Coeff(2) > 0 {
		// Curves up, mo money is mo money
		return 0.4 * c.gameState().ClientInfo.Resources
	}

	// Curves down, mo money is not always mo money
	// c.Logf("%v x^2 + %v x + %v", r.Coeff(2), r.Coeff(1), r.Coeff(0))
	// c.Logf("%v / %v = %v ", -r.Coeff(1), 2*r.Coeff(2), shared.Resources(-r.Coeff(1)/(2*r.Coeff(2))))
	return shared.Resources(-r.Coeff(1) / (2 * r.Coeff(2)))
}

func desperateDecider(c client) (shared.ForageDecision, shared.Resources) {
	forageType := bestROIForageType(c.forageHistory)
	if forageType == shared.ForageType(-1) {
		forageType = shared.DeerForageType
	}

	contribution := c.gameState().ClientInfo.Resources

	r := outcomeRegression(c.forageHistory[forageType])
	expectedRewardF, _ := r.Predict([]float64{float64(contribution)})

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: contribution,
	}, shared.Resources(expectedRewardF)
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	decision, expectedReward := c.config.forageDecider(*c)
	c.expectedForageReward = expectedReward

	return decision, nil
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, revenue shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		contribution: forageDecision.Contribution,
		revenue:      revenue,
		turn:         c.gameState().Turn,
		participant:  c.GetID(),
	})

	c.Logf(
		"[Forage result]: %v(%05.3f) | Expectation: %+05.3f | Reward: %+05.3f | Error: %.0f%%",
		forageDecision.Type,
		forageDecision.Contribution,
		c.expectedForageReward,
		revenue,
		((c.expectedForageReward-revenue)/revenue)*100,
	)
}

/************************/
/***  Forage Helpers  ***/
/************************/

// bestROIForageType finds the ForageType that has the best average ROI. If all
// forageTypes have negative returns then it returns shared.ForageType(-1)
func bestROIForageType(forageHistory ForageHistory) shared.ForageType {
	bestForageType := shared.ForageType(-1)
	bestROI := 0.0

	for forageType, outcomes := range forageHistory {
		ROIsum := 0.0
		for _, outcome := range outcomes {
			if outcome.contribution != 0 {
				ROIsum += outcome.ROI()
			}
		}

		averageROI := ROIsum / float64(len(outcomes))

		if averageROI > bestROI {
			bestROI = averageROI
			bestForageType = forageType
		}
	}

	return bestForageType
}

func bestROIOutcome(history []ForageOutcome) ForageOutcome {
	bestOutcome := ForageOutcome{}

	for _, outcome := range history {
		if outcome.profit() >= 10 && outcome.contribution != 0 && outcome.ROI() > bestOutcome.ROI() {
			bestOutcome = outcome
		}
	}
	return bestOutcome
}
