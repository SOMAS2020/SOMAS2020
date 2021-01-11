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

func (c client) randomForage() shared.ForageDecision {
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
	}
}

func (c client) roiForage() shared.ForageDecision {
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
	return forageDecision
}

func (c client) regressionForage() (shared.ForageDecision, error) {
	bestReward := shared.Resources(0)
	var decision shared.ForageDecision

	for forageType := range c.forageHistory {
		var expectedReward shared.Resources
		var contribution shared.Resources
		// Regression throws an error when the size of array is less than 2
		if len(c.forageHistory[forageType]) < 2 {
			if forageType == shared.DeerForageType {
				c.Logf("[Forage decision] Ha! jokes. Flipping instead")
				contribution = c.flipForage().Contribution
			} else {
				c.Logf("[Forage decision] Ha! jokes. random instead")
				contribution = shared.Resources(0.1*rand.Float64()) * c.gameState().ClientInfo.Resources
			}
			decision = shared.ForageDecision{
				Type:         forageType,
				Contribution: contribution,
			}
		} else {
			r, err := outcomeRegression(c.forageHistory[forageType])
			if err != nil {
				return shared.ForageDecision{}, err
			}
			contribution := c.regressionOptimalContribution(r)
			if contribution > c.gameState().ClientInfo.Resources {
				contribution = c.gameState().ClientInfo.Resources
			}

			expectedRewardF, err := r.Predict([]float64{float64(contribution)})
			if err != nil {
				return shared.ForageDecision{}, err
			}
			expectedReward = shared.Resources(expectedRewardF)

			if expectedReward > bestReward {
				bestReward = expectedReward
				decision = shared.ForageDecision{
					Type:         forageType,
					Contribution: contribution,
				}
			}
		}
	}

	return decision, nil
}

// flipDecider does the opposite of what the mass did the previous turn. It
// forages deer, with an amount inversely proportional to the sum of contributed
// resources
func (c client) flipForage() shared.ForageDecision {
	resources := c.gameState().ClientInfo.Resources
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

	if c.forageType == shared.FishForageType {
		return shared.ForageDecision{
			Contribution: shared.Resources(0.1*rand.Float64()) * c.gameState().ClientInfo.Resources,
			Type:         shared.FishForageType,
		}
	}

	if c.forageType == shared.DeerForageType && (totalContributionLastTurn == shared.Resources(0) || totalHuntersLastTurn == 0) {
		contribution := shared.Resources(math.Min(
			float64(shared.Resources(
				c.config.forageContributionCapPercent)*resources,
			),
			float64(c.config.soloDeerHuntContribution),
		))
		// Big contribution
		return shared.ForageDecision{
			Contribution: contribution,
			Type:         shared.DeerForageType,
		}
	}

	// Proxy for population
	totalROI := totalRevenueLastTurn / totalContributionLastTurn
	averageContribution := totalContributionLastTurn / shared.Resources(totalHuntersLastTurn)
	contribution := shared.Resources(c.config.flipForageScale) * totalROI * averageContribution
	contribution += shared.Resources(c.config.forageContributionNoisePercent) * resources
	contribution = shared.Resources(math.Min(
		float64(shared.Resources(
			c.config.forageContributionCapPercent)*resources,
		),
		float64(contribution),
	))
	c.Logf("[Forage decision] flipping results: %v", contribution)
	return shared.ForageDecision{
		Contribution: contribution,
		Type:         shared.DeerForageType,
	}
}

func (c client) constantForage(resourcePercent float64) shared.ForageDecision {
	decision := shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: 0.2 * c.gameState().ClientInfo.Resources,
	}

	c.Logf("[Forage decision]: constant (%v)", decision)
	return decision
}

func outcomeRegression(history []ForageOutcome) (regression.Regression, error) {
	r := new(regression.Regression)
	r.SetObserved("Team1 reward")
	r.SetVar(0, "Team1 contribution")

	for _, outcome := range history {
		r.Train(regression.DataPoint(float64(outcome.revenue), []float64{float64(outcome.contribution)}))
	}

	r.AddCross(regression.PowCross(0, 2))
	err := r.Run()

	return *r, err
}

func (c *client) regressionOptimalContribution(r regression.Regression) shared.Resources {
	if r.Coeff(2) >= 0 {
		// Curves up, mo money is mo money
		return 0.4 * c.gameState().ClientInfo.Resources
	}
	// r.Coeff(2) < 0. Curves down, mo money is not always mo money
	return shared.Resources(-r.Coeff(1) / (2 * r.Coeff(2)))
}

func (c client) desperateForage() shared.ForageDecision {
	forageType := bestROIForageType(c.forageHistory)
	if forageType == shared.ForageType(-1) {
		forageType = shared.DeerForageType
	}

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: c.gameState().ClientInfo.Resources,
	}
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.randomForageTurns {
		c.Logf("[Forage decision]: random")
		return c.randomForage(), nil
	} else if c.emotionalState() == Desperate {
		c.Logf("[Forage decision]: desperate")
		return c.desperateForage(), nil
	} else {
		c.Logf("[Forage decision]: flip")
		return c.flipForage(), nil
	}
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, revenue shared.Resources, numberCaught uint) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		contribution: forageDecision.Contribution,
		revenue:      revenue,
		turn:         c.gameState().Turn,
		participant:  c.GetID(),
	})

	notEnoughMoney := revenue < 2*c.gameConfig().CostOfLiving
	if notEnoughMoney || revenue/forageDecision.Contribution < 1 {
		switch c.forageType {
		case shared.DeerForageType:
			c.forageType = shared.FishForageType
		case shared.FishForageType:
			c.forageType = shared.DeerForageType
		}
	}

	if revenue == 0 {
		c.Logf(
			"[Forage result]: %v(%05.3f) | Reward: 0 ",
			forageDecision.Type,
			forageDecision.Contribution,
		)
	} else {
		c.Logf(
			"[Forage result]: %v(%05.3f) | Reward: %+05.3f | Error: %.0f%%",
			forageDecision.Type,
			forageDecision.Contribution,
			revenue,
			((c.expectedForageReward-revenue)/revenue)*100,
		)
	}
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

		if len(outcomes) == 0 {
			panic("Empty history for forageType in team1's forageHistory")
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
