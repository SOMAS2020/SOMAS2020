package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared" 

	"github.com/sajari/regression"
)

/********************/
/***  Desperation   */
/********************/

func (c *client) desperateForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         bestROIForageType(c.forageHistory),
		Contribution: c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Forage][Decision]: Desperate | Decision %v", forageDecision)
	return forageDecision
}

/********************/
/***    Foraging    */
/********************/

func (c *client) randomForage() shared.ForageDecision {
	// Up to 10% of our current resources
	forageContribution := shared.Resources(0.1*rand.Float64()) * c.gameState().ClientInfo.Resources

	var forageType shared.ForageType
	if rand.Float64() < 0.5 {
		forageType = shared.DeerForageType
	} else {
		forageType = shared.FishForageType
	}

	c.Logf("[Forage][Decision]:Random")
	// TODO Add fish foraging when it's done
	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}
}

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
	return shared.Resources(-r.Coeff(1) / 2*r.Coeff(2))
}

func (c *client) normalForage() shared.ForageDecision {
	forageType := bestROIForageType(c.forageHistory)
	if forageType == shared.ForageType(-1) {
		// Here we have found that the best idea is to not forage
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}
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
	c.Logf(
		"[Forage][Decision]: Decision %v | Expected ROI %v",
		forageDecision, expectedOutcome.ROI())

	return forageDecision
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.randomForageTurns {
		return c.randomForage(), nil
	} else if c.emotionalState() == Desperate {
		return c.desperateForage(), nil
	} else {
		return c.normalForage(), nil
	}
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, revenue shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		contribution: forageDecision.Contribution,
		revenue:      revenue,
		turn:         c.gameState().Turn,
	})

	c.Logf(
		"[Forage][Update]: ForageType %v | Profit %v | Contribution %v | Actual ROI %v",
		forageDecision.Type,
		revenue-forageDecision.Contribution,
		forageDecision.Contribution,
		(revenue/forageDecision.Contribution)-1,
	)
}