package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat"
)

type forecastInfo struct {
	epiX       shared.Coordinate // x co-ord of disaster epicentre
	epiY       shared.Coordinate // y ""
	mag        shared.Magnitude
	turn       uint
	confidence float64
}

type forecastHistory map[uint]forecastInfo                                   // stores history of past disasters
type receivedForecastHistory map[uint]shared.ReceivedDisasterPredictionsDict // stores history of received disasters

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {

	spatialMagPred := c.estimateSpatialAndMag() // estimate for x,y coords and magnitude
	estPeriod, periodConf := c.estimateDisasterPeriod()

	lastDisasterTurn := c.getLastDisasterTurn()

	prediction := shared.DisasterPrediction{
		CoordinateX: spatialMagPred.epiX,
		CoordinateY: spatialMagPred.epiY,
		Magnitude:   spatialMagPred.mag,
		TimeLeft:    int(lastDisasterTurn + estPeriod - c.getTurn()),
	}
	pBias := c.config.periodConfidenceBias
	if math.Abs(pBias) > 1 {
		c.Logf("WARNING: Invalid period confidence bias value of %v. Setting default value of 0.5.", pBias)
		pBias = 0.5 // assign default if out of range
	}
	prediction.Confidence = periodConf*pBias + (1-pBias)*spatialMagPred.confidence
	trustedIslandIDs := []shared.ClientID{}
	trustThresh := c.config.forecastTrustTreshold
	for id := range c.getTrustedTeams(trustThresh, false, forecastingBasis) {
		trustedIslandIDs = append(trustedIslandIDs, id)
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslandIDs,
	}
	c.lastDisasterPrediction = prediction
	// update forecast history
	c.forecastHistory[c.getTurn()] = forecastInfo{
		epiX:       prediction.CoordinateX,
		epiY:       prediction.CoordinateY,
		mag:        prediction.Magnitude,
		turn:       uint(prediction.TimeLeft) + c.getTurn(),
		confidence: prediction.Confidence,
	}
	return predictionInfo
}

func (c client) getLastDisasterTurn() uint {
	sortedTurns := c.disasterHistory.sortKeys()
	l := len(sortedTurns)
	if l > 0 {
		return sortedTurns[l-1]
	}
	return 0
}

// provides estimate of *when* next disaster will occur and associated conf
func (c *client) estimateDisasterPeriod() (period uint, conf float64) {

	if len(c.disasterHistory) == 0 {
		return 0, 0 // we can't make any predictions with no disaster history!
	}
	periods := []float64{} // use float so we can use stat.Variance() later
	periodSum := 0.0       // to offset this from average
	prevTurn := float64(startTurn)
	for _, turn := range c.disasterHistory.sortKeys() { // TODO: find instances where assumption of ordered map keys is relied upon
		periods = append(periods, float64(turn)-prevTurn) // period = no. turns between successive disasters
		periodSum += periods[len(periods)-1]
		prevTurn = float64(turn)
	}
	c.Logf("Periods final: %v", periods)
	if len(periods) == 1 {
		return uint(periods[0]), 50.0 // if we only have one past observation. Best we can do is estimate that period again.
	}
	// if we have more than 1 observation
	v := stat.Variance(periods, nil)

	meanPeriod := periodSum / float64(len(periods))
	varThresh := meanPeriod
	varianceRatio := math.Min(v/varThresh, 1.0) // should be between 0 (min var) and 1 (max var)

	conf = (1 - varianceRatio) * 100
	// if not consistent, return mean period we've seen so far
	return uint(meanPeriod), conf
}

// gets confidence of x,y coord and magnitude estimates
func (c *client) estimateSpatialAndMag() forecastInfo {
	sumX, sumY, sumMag := 0.0, 0.0, 0.0

	for _, dInfo := range c.disasterHistory {
		sumX += dInfo.report.X
		sumY += dInfo.report.Y
		sumMag += dInfo.report.Y
	}
	n := float64(len(c.disasterHistory))
	historicalInfo := forecastInfo{
		epiX: sumX / n,
		epiY: sumY / n,
		mag:  sumMag / n,
		turn: uint(n), // this will be updated by period forecast
	}
	totalDisaster := forecastInfo{}
	sqDiff := func(a, b float64) float64 { return math.Pow(a-b, 2) }
	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, d := range c.forecastHistory {
		totalDisaster.epiX += sqDiff(d.epiX, historicalInfo.epiX)
		totalDisaster.epiY += sqDiff(d.epiY, historicalInfo.epiY)
		totalDisaster.mag += sqDiff(d.mag, historicalInfo.mag)
	}

	// TODO: find a better method of calculating confidence
	// Find the sum of the variances and the average variance
	variance := (totalDisaster.epiX + totalDisaster.epiY + totalDisaster.mag) / float64(len(c.forecastHistory))
	variance = math.Min(c.config.maxForecastVariance, variance)

	historicalInfo.confidence = c.config.maxForecastVariance - variance
	return historicalInfo
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence

	sumX, sumY, sumMag, sumConf := 0.0, 0.0, 0.0, 0.0
	sumTime := 0

	c.updateForecastingReputations(receivedPredictions) // update our perceptions of other teams

	//c.lastDisasterForecast.Confidence *= 1.3 // inflate confidence of our prediction above others
	receivedPredictions[ourClientID] = shared.ReceivedDisasterPredictionInfo{PredictionMade: c.lastDisasterPrediction, SharedFrom: ourClientID}

	//TODO: decide whether our prediction should be included in this history or not
	c.receivedForecastHistory[c.getTurn()] = receivedPredictions // update rxForecastsHistory

	// weight predictions by their confidence and our assessment of their forecasting reputation
	for rxTeam, pred := range receivedPredictions {
		rep := float64(c.opinions[rxTeam].getForecastingRep()) + 1 // our notion of another island's forecasting reputation
		sumX += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateX * rep
		sumY += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateY * rep
		sumMag += pred.PredictionMade.Confidence * pred.PredictionMade.Magnitude * rep
		sumTime += int(pred.PredictionMade.Confidence) * pred.PredictionMade.TimeLeft * int(rep)
		sumConf += pred.PredictionMade.Confidence * rep
	}

	sumConf = math.Max(sumConf, 1) // guard against div by zero error below
	// Finally get the final prediction generated by considering predictions from all islands that we have available
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: sumX / sumConf,
		CoordinateY: sumY / sumConf,
		Magnitude:   sumMag / sumConf,
		TimeLeft:    int((float64(sumTime) / sumConf) + 0.5),     // +0.5 for rounding
		Confidence:  sumConf / float64(len(receivedPredictions)), // this len will always be >= 1
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}

func (c *client) updateForecastingReputations(receivedPredictions shared.ReceivedDisasterPredictionsDict) {

	for team, predInfo := range receivedPredictions {
		// if teams make predictions with conf > 50% before first disaster, downgrade their rep by 75%
		if len(c.disasterHistory) == 0 {
			if predInfo.PredictionMade.Confidence > 50 {
				c.opinions[team].updateOpinion(forecastingBasis, -0.75)
			}
		}
		// decrease trust in teams who are overly confident
		if predInfo.PredictionMade.Confidence > 98 {
			c.opinions[team].updateOpinion(forecastingBasis, -0.3)
		}
		// TODO: add more sophisticated opinion forming
	}

}
