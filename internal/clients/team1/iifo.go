package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type disaster struct {
	meanDisaster        disasters.DisasterReport
	allDisaster         []disasters.DisasterReport
	disasterTurnCounter uint
	numberOfDisasters   uint
	meanDisasterTurn    float64
	estimatedDDay       uint
}

/**************************/
/*** 	  Foraging	 	***/
/**************************/

func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID

	shareTo = c.aliveClients()

	lastDecisionTurn := -1
	var lastDecision shared.ForageDecision
	var lastRevenue shared.Resources

	for forageType, outcomes := range c.forageHistory {
		for _, outcome := range outcomes {
			if int(outcome.turn) > lastDecisionTurn {
				lastDecisionTurn = int(outcome.turn)
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: outcome.contribution,
				}
				lastRevenue = outcome.revenue
			}
		}
	}

	if lastDecisionTurn < 0 {
		shareTo = []shared.ClientID{}
	}

	forageInfo := shared.ForageShareInfo{
		ShareTo:          shareTo,
		ResourceObtained: lastRevenue,
		DecisionMade:     lastDecision,
	}

	c.Logf("Sharing forage info: %v", forageInfo)
	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfos []shared.ForageShareInfo) {
	for _, forageInfo := range forageInfos {
		c.forageHistory[forageInfo.DecisionMade.Type] =
			append(
				c.forageHistory[forageInfo.DecisionMade.Type],
				ForageOutcome{
					participant:  forageInfo.SharedFrom,
					contribution: forageInfo.DecisionMade.Contribution,
					revenue:      forageInfo.ResourceObtained,
					turn:         c.gameState().Turn - 1,
				},
			)
	}
}

/******************************/
/*** 		 Disasters 		  */
/******************************/

// findConfidence is only called when a disaster has happened. Therefore len(disasterHistory) > 0
func (c client) findConfidence() float64 {
	disasterHistory := c.disasterInfo.allDisaster
	meanDisaster := c.disasterInfo.meanDisaster
	totalDisaster := disasters.DisasterReport{}
	for _, disaster := range disasterHistory {
		totalDisaster.X += math.Pow(disaster.X-meanDisaster.X, 2)
		totalDisaster.Y += math.Pow(disaster.Y-meanDisaster.Y, 2)
		totalDisaster.Magnitude += math.Pow(disaster.Magnitude-meanDisaster.Magnitude, 2)
	}

	disasterHistorySize := float64(len(disasterHistory))
	if disasterHistorySize == 0 {
		return 0
	}
	sqrtDisasterHistory := math.Sqrt(disasterHistorySize)
	xSD := math.Sqrt(totalDisaster.X / disasterHistorySize)
	ySD := math.Sqrt(totalDisaster.Y / disasterHistorySize)
	magSD := math.Sqrt(totalDisaster.Magnitude / disasterHistorySize)

	// 1.645 is Z value for 90% Confidence Interval
	// See link for more maths: https://www.mathsisfun.com/data/confidence-interval.html
	// Formula: meanX +- Z( sdX / sqrt(n) )
	// Z( sdX / sqrt(n) ) is the upper and lower bounds and therefore can be transformed
	// to a confidence percentage by (meanX - Z(sdX / sqrt(n))) / meanX
	confidenceIntervalX, confidenceIntervalY, confidenceIntervalM := 0.0, 0.0, 0.0
	if meanDisaster.X != 0 {
		confidenceIntervalX = 1 - (1.645 * xSD / (sqrtDisasterHistory * meanDisaster.X))
	}
	if meanDisaster.Y != 0 {
		confidenceIntervalY = 1 - (1.645 * ySD / (sqrtDisasterHistory * meanDisaster.Y))
	}
	if meanDisaster.Magnitude != 0 {
		confidenceIntervalM = 1 - (1.645 * magSD / (sqrtDisasterHistory * meanDisaster.Magnitude))
	}

	// Return average
	return (confidenceIntervalX + confidenceIntervalY + confidenceIntervalM) / 3
}

func (c *client) DisasterNotification(disaster disasters.DisasterReport, effect disasters.DisasterEffects) {
	turnCounter := c.disasterInfo.disasterTurnCounter
	if disaster.Magnitude != 0 {
		c.disasterInfo.allDisaster = append(c.disasterInfo.allDisaster, disaster)
		if c.disasterInfo.numberOfDisasters == 0 {
			c.disasterInfo.meanDisaster = disaster
			c.disasterInfo.numberOfDisasters++
			c.disasterInfo.meanDisasterTurn = float64(c.gameState().Turn)
		} else {
			numOfDisasters := c.disasterInfo.numberOfDisasters
			denominator := float64(c.disasterInfo.numberOfDisasters + 1)
			c.disasterInfo.meanDisaster.X = (c.disasterInfo.meanDisaster.X*float64(numOfDisasters) + disaster.X) / denominator
			c.disasterInfo.meanDisaster.Y = (c.disasterInfo.meanDisaster.Y*float64(numOfDisasters) + disaster.Y) / denominator
			c.disasterInfo.meanDisaster.Magnitude = (c.disasterInfo.meanDisaster.Magnitude*float64(numOfDisasters) + disaster.Magnitude) / denominator
			c.disasterInfo.meanDisasterTurn = (c.disasterInfo.meanDisasterTurn*float64(numOfDisasters) + float64(c.disasterInfo.disasterTurnCounter)) / denominator
		}
		// TODO: Improve this estimation by keeping a track of a historgram of daysSinceDisaster -> howManyDisasterHappened.
		c.disasterInfo.estimatedDDay = c.gameState().Turn + uint(c.disasterInfo.meanDisasterTurn)
		c.disasterInfo.disasterTurnCounter = 0
	}

	for id, team := range c.othersDisasterPrediction {
		timeDistance := team.PredictionMade.TimeLeft
		trustValue := float64(turnCounter) * team.PredictionMade.Confidence
		if timeDistance == 0 {
			c.trustTeams[id] += trustValue
		} else {
			// Add one to avoid divide by 1 case.
			c.trustTeams[id] += trustValue / (float64(timeDistance) + 1)
		}
	}
}

// MakeDisasterPrediction evaluates the mean of X, Y, Magnitude, Turn
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	c.disasterInfo.disasterTurnCounter++
	currTurn := c.gameState().Turn
	timeLeft := c.disasterInfo.estimatedDDay - currTurn
	confidence := 0.0
	if c.disasterInfo.estimatedDDay < currTurn {
		timeLeft = 0
	}

	if c.disasterInfo.numberOfDisasters == 0 {
		disasterPrediction := shared.DisasterPrediction{
			CoordinateX: rand.Float64() * 10,
			CoordinateY: rand.Float64() * 10,
			Magnitude:   rand.Float64(),
			Confidence:  confidence,
			TimeLeft:    timeLeft,
		}
		return shared.DisasterPredictionInfo{
			PredictionMade: disasterPrediction,
			TeamsOfferedTo: c.aliveClients(),
		}
	}

	if len(c.disasterInfo.allDisaster) > 1 {
		// TODO: Cache confidence level
		confidence = c.findConfidence()
	}
	disasterPrediction := shared.DisasterPrediction{
		CoordinateX: c.disasterInfo.meanDisaster.X,
		CoordinateY: c.disasterInfo.meanDisaster.Y,
		Magnitude:   c.disasterInfo.meanDisaster.Magnitude,
		TimeLeft:    timeLeft,
		// TODO: Add timeLeft to confidence level
		Confidence: confidence,
	}

	// Store own disasterPrediction for evaluation in DisasterNotification
	c.othersDisasterPrediction[c.GetID()] = shared.ReceivedDisasterPredictionInfo{
		PredictionMade: disasterPrediction,
		SharedFrom:     c.GetID(),
	}

	return shared.DisasterPredictionInfo{
		PredictionMade: disasterPrediction,
		TeamsOfferedTo: c.aliveClients(),
	}
}

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for id, predictions := range receivedPredictions {
		if predictions.PredictionMade.TimeLeft+1 != c.othersDisasterPrediction[id].PredictionMade.TimeLeft {
			c.trustTeams[id]--
		}
	}

	for id, predictions := range receivedPredictions {
		c.othersDisasterPrediction[id] = predictions
	}
}
