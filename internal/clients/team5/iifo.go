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

type forecastHistory map[uint]forecastInfo // stores history of past disasters

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {

	meanDisaster := c.getHistoricalForecast()
	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.epiX,
		CoordinateY: meanDisaster.epiY,
		Magnitude:   meanDisaster.mag,
		TimeLeft:    int(meanDisaster.turn - c.getTurn()),
	}

	prediction.Confidence = c.determineForecastConfidence()
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
		epiX: prediction.CoordinateX,
		epiY: prediction.CoordinateY,
		mag:  prediction.Magnitude,
		turn: uint(prediction.TimeLeft) + c.getTurn(),
	}
	return predictionInfo
}

// averages observations over history to get 'mean' disaster
func (c client) getHistoricalForecast() forecastInfo {
	sumX, sumY, sumMag := 0.0, 0.0, 0.0

	for _, dInfo := range c.disasterHistory {
		sumX += dInfo.report.X
		sumY += dInfo.report.Y
		sumMag += dInfo.report.Y
	}
	n := float64(len(c.forecastHistory))
	period, conf := c.analyseDisasterPeriod()

	meanDisaster := forecastInfo{
		epiX:       sumX / n,
		epiY:       sumY / n,
		mag:        sumMag / n,
		turn:       c.getLastDisasterTurn() + period,
		confidence: conf,
	}
	return meanDisaster
}

func (c client) getLastDisasterTurn() uint {
	sortedTurns := c.disasterHistory.sortKeys()
	l := len(sortedTurns)
	if l > 0 {
		return sortedTurns[l-1]
	}
	return 0
}

func (c *client) analyseDisasterPeriod() (period uint, conf float64) {
	if len(c.disasterHistory) == 0 {
		return 0, 0 // we can't make any predictions with no disaster history!
	}
	periods := []float64{0} // use float so we can use stat.Variance() later
	periodSum := 0.0
	for turn := range c.disasterHistory {
		periods = append(periods, float64(turn)-periods[len(periods)-1]) // period = no. turns between successive disasters
		periodSum += periods[len(periods)-1]
	}
	periods = periods[1:] // remove leading 0
	if len(periods) == 1 {
		return uint(periods[0]), 50.0 // if we only have one past observation. Best we can do is estimate that period again.
	}
	// if we have more than 1 observation
	v := stat.Variance(periods, nil)

	meanPeriod := periodSum / float64(len(periods))
	varThresh := meanPeriod / 2
	varianceRatio := math.Max(v/varThresh, 1.0) // should be between 0 (min var) and 1 (max var)
	conf = (1 - varianceRatio) * 100
	// if not consistent, return mean period we've seen so far
	return uint(meanPeriod), conf
}

func (c *client) determineForecastConfidence() float64 {
	totalDisaster := forecastInfo{}
	sqDiff := func(x, meanX float64) float64 { return math.Pow(x-meanX, 2) }
	meanInfo := c.getHistoricalForecast()
	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, d := range c.forecastHistory {
		totalDisaster.epiX += sqDiff(d.epiX, meanInfo.epiX)
		totalDisaster.epiY += sqDiff(d.epiY, meanInfo.epiY)
		totalDisaster.mag += sqDiff(d.mag, meanInfo.mag)
	}

	// TODO: find a better method of calculating confidence
	// Find the sum of the variances and the average variance
	variance := (totalDisaster.epiX + totalDisaster.epiY + totalDisaster.mag) / float64(len(c.forecastHistory))
	variance = math.Min(c.config.maxForecastVariance, variance)

	return c.config.maxForecastVariance - variance
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence

	sumX, sumY, sumMag, sumConf := 0.0, 0.0, 0.0, 0.0
	sumTime := 0

	//c.lastDisasterForecast.Confidence *= 1.3 // inflate confidence of our prediction above others
	receivedPredictions[ourClientID] = shared.ReceivedDisasterPredictionInfo{PredictionMade: c.lastDisasterPrediction, SharedFrom: ourClientID}

	// weight predictions by their confidence and our assessment of their forecasting reputation
	for rxTeam, pred := range receivedPredictions {
		rep := float64(c.opinions[rxTeam].getForecastingRep()) + 1 // our notion of another island's forecasting reputation
		sumX += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateX * rep
		sumY += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateY * rep
		sumMag += pred.PredictionMade.Confidence * pred.PredictionMade.Magnitude * rep
		sumTime += int(pred.PredictionMade.Confidence) * pred.PredictionMade.TimeLeft * int(rep)
		sumConf += pred.PredictionMade.Confidence * rep
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: sumX / sumConf,
		CoordinateY: sumY / sumConf,
		Magnitude:   sumMag / sumConf,
		TimeLeft:    int((float64(sumTime) / sumConf) + 0.5), // +0.5 for rounding
		Confidence:  sumConf / float64(len(receivedPredictions)),
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}
