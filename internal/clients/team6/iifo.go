package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

func (c client) checkIsStochastic() bool {
	isStochastic := false

	if len(c.disastersHistory) < 3 {
		// treats this as stochastic when data not enough to judge
		isStochastic = true
	}

	for i := 0; i < len(c.disastersHistory)-2; i++ {
		if c.disastersHistory[i+1].Turn-c.disastersHistory[i].Turn != c.disastersHistory[i+2].Turn-c.disastersHistory[i+1].Turn {
			isStochastic = true
		}
	}

	return isStochastic
}

func (c client) getPseudoPeriod(isStochastic bool) float64 {
	pseudoPeriod := float64(0)

	if len(c.disastersHistory) < 1 {
		return pseudoPeriod
	}

	if isStochastic {
		pseudoPeriod = float64(c.ServerReadHandle.GetGameState().Turn) / float64(len(c.disastersHistory))
	} else {
		pseudoPeriod = float64(c.disastersHistory[len(c.disastersHistory)-1].Turn) / float64(len(c.disastersHistory))
	}

	return pseudoPeriod
}

func (c client) getMeanMagnitude() float64 {
	meanMagnitudeFromHistory := float64(0)
	meanMagnitudeFromPredictions := float64(0)
	magnitudeSumOfHistory := float64(0)
	magnitudeSumOfPredictions := float64(0)

	if len(c.disastersHistory)+len(c.disasterPredictions) < 1 {
		return float64(0)
	}

	if len(c.disastersHistory) != 0 {
		for _, disaster := range c.disastersHistory {
			magnitudeSumOfHistory += disaster.Magnitude
		}

		meanMagnitudeFromHistory = magnitudeSumOfHistory / float64(len(c.disastersHistory))
	}

	if len(c.disasterPredictions) != 0 {
		for _, disaster := range c.disasterPredictions {
			magnitudeSumOfPredictions += disaster.Magnitude * (disaster.Confidence) / float64(100)
		}

		meanMagnitudeFromPredictions = magnitudeSumOfPredictions / float64(len(c.disasterPredictions))
	}

	return (meanMagnitudeFromHistory + meanMagnitudeFromPredictions) / float64(2)
}

func (c client) getTimeLeft(isStochastic bool, period float64) uint {
	timeLeft := uint(0)
	turnOfLastDisaster := uint(0)

	if !isStochastic {
		if len(c.disastersHistory) != 0 {
			turnOfLastDisaster = c.disastersHistory[len(c.disastersHistory)-1].Turn
		}

		timeLeft = uint(period) - (c.ServerReadHandle.GetGameState().Turn - turnOfLastDisaster)
	}

	return timeLeft
}

func (c client) determineConfidence(isStochastic bool, pseudoPeriod float64) float64 {
	confidence := float64(100)

	if isStochastic {
		confidence = confidence / pseudoPeriod
	}

	return confidence
}

func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	xMax := c.ServerReadHandle.GetGameState().Geography.XMax
	xMin := c.ServerReadHandle.GetGameState().Geography.XMin
	yMax := c.ServerReadHandle.GetGameState().Geography.YMax
	yMin := c.ServerReadHandle.GetGameState().Geography.YMin

	// two random variables generated from uniform distribution
	predictedX := distuv.Uniform{Min: xMin, Max: xMax}.Rand()
	predictedY := distuv.Uniform{Min: yMin, Max: yMax}.Rand()

	// check if if these parameters are exposed
	isIsStochasticExposed := c.ServerReadHandle.GetGameConfig().DisasterConfig.StochasticDisasters.Valid
	isPeriodExposed := c.ServerReadHandle.GetGameConfig().DisasterConfig.DisasterPeriod.Valid
	isStochastic := c.ServerReadHandle.GetGameConfig().DisasterConfig.StochasticDisasters.Value
	period := float64(c.ServerReadHandle.GetGameConfig().DisasterConfig.DisasterPeriod.Value)

	if !isPeriodExposed && !isIsStochasticExposed {
		isStochastic = c.checkIsStochastic()
		period = c.getPseudoPeriod(isStochastic)
	} else if !isIsStochasticExposed {
		isStochastic = c.checkIsStochastic()
	} else if !isPeriodExposed {
		period = c.getPseudoPeriod(isStochastic)
	}

	prediction := shared.DisasterPrediction{
		CoordinateX: predictedX,
		CoordinateY: predictedY,
		Magnitude:   c.getMeanMagnitude(),
	}
	teamsOfferingTo := []shared.ClientID{}

	if period != 0 {
		prediction.TimeLeft = c.getTimeLeft(isStochastic, period)
		prediction.Confidence = c.determineConfidence(isStochastic, period)
		teamsOfferingTo = shared.TeamIDs[:]
	}

	c.disasterPredictions[c.GetID()] = prediction

	return shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: teamsOfferingTo,
	}
}

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for _, prediction := range receivedPredictions {
		c.disasterPredictions[prediction.SharedFrom] = shared.DisasterPrediction{
			CoordinateX: c.trustRank[prediction.SharedFrom] * prediction.PredictionMade.CoordinateX,
			CoordinateY: c.trustRank[prediction.SharedFrom] * prediction.PredictionMade.CoordinateY,
			Magnitude:   c.trustRank[prediction.SharedFrom] * prediction.PredictionMade.Magnitude,
			TimeLeft:    uint(c.trustRank[prediction.SharedFrom]) * prediction.PredictionMade.TimeLeft,
			Confidence:  prediction.PredictionMade.Confidence,
		}
	}

	c.Logf("Final Prediction: [%v]", c.disasterPredictions[c.GetID()])
}

func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID // containing agents our agent wish to share informationwith

	for id, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status != shared.Dead {
			shareTo = append(shareTo, id)
		}
	}

	var lastDecision shared.ForageDecision
	var lastForageOut shared.Resources

	for forageType, results := range c.forageHistory {
		for _, result := range results {
			if uint(result.turn) == c.ServerReadHandle.GetGameState().Turn-1 {
				lastForageOut = result.forageReturn
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: result.forageIn,
				}
			}
		}
	}

	forageInfo := shared.ForageShareInfo{
		DecisionMade:     lastDecision,
		ResourceObtained: lastForageOut,
		ShareTo:          shareTo,
	}

	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfo []shared.ForageShareInfo) {
	for _, val := range forageInfo {
		c.forageHistory[val.DecisionMade.Type] =
			append(
				c.forageHistory[val.DecisionMade.Type],
				ForageResults{
					forageIn:     val.DecisionMade.Contribution,
					forageReturn: val.ResourceObtained,
					turn:         c.ServerReadHandle.GetGameState().Turn,
				},
			)
	}
}
