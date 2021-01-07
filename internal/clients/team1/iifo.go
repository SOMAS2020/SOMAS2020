package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type disaster struct {
	meanDisaster        disasters.DisasterReport
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

	for id, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead {
			shareTo = append(shareTo, id)
		}
	}

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

func (c *client) DisasterNotification(disaster disasters.DisasterReport, effect disasters.DisasterEffects) {
	if disaster.Magnitude != 0 {
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
		// TODO: Improvement by keeping a track of a historgram of daysSinceDisaster -> howManyDisasterHappened.
		c.disasterInfo.estimatedDDay = c.gameState().Turn + uint(c.disasterInfo.meanDisasterTurn)
		c.disasterInfo.disasterTurnCounter = 0
	}
}

// MakeDisasterPrediction evaluates the mean of X, Y, Magnitude, Turn
// Confidence doesn't mean much.
// If there is nothing in Disaster
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	c.disasterInfo.disasterTurnCounter++
	currTurn := int(c.gameState().Turn)

	timeLeft := int(c.disasterInfo.estimatedDDay) - currTurn
	// TODO: If timeLeft is less than 0, increase confidence
	var confidence float64

	if c.disasterInfo.numberOfDisasters == 0 {
		if timeLeft < 0 {
			confidence = 0.1 * float64(currTurn)
		}
		disasterPrediction := shared.DisasterPrediction{
			CoordinateX: rand.Float64() * 10,
			CoordinateY: rand.Float64() * 10,
			Magnitude:   rand.Float64(),
			Confidence:  0 + confidence,
			TimeLeft:    timeLeft,
		}
		c.Logf("[DISASTER INFO] is empty. Creating random prediction: %v", disasterPrediction)
		return shared.DisasterPredictionInfo{
			PredictionMade: disasterPrediction,
			TeamsOfferedTo: c.aliveClients,
		}
	}

	if timeLeft < 0 {
		confidence = -float64(timeLeft) / math.Pow(float64(c.disasterInfo.meanDisasterTurn), 2)
	}

	// TODO: Calculate SD (confidence) for CoordinateX and CoordinateY and sum to confidence
	disasterPrediction := shared.DisasterPrediction{
		CoordinateX: c.disasterInfo.meanDisaster.X + rand.Float64(),
		CoordinateY: c.disasterInfo.meanDisaster.Y + rand.Float64(),
		Magnitude:   c.disasterInfo.meanDisaster.Magnitude,
		TimeLeft:    timeLeft,
		Confidence:  0.25 + confidence,
	}
	c.Logf("[DISASTER INFO] Creating prediction: %v", disasterPrediction)
	return shared.DisasterPredictionInfo{
		PredictionMade: disasterPrediction,
		TeamsOfferedTo: c.aliveClients,
	}
}

// TODO: Store all prediction. Based on their confidence, create a trustworthy map of other islands.
// Those who are more trustworthy (including us), we use their prediction.
// Trustworthy is a point/score system
// Calculated the distance from timeLeft and the currTurn where disaster happened, take into account
// confidence level,
func (c client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
}
