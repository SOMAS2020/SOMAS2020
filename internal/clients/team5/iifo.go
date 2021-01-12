package team5

import (
	"fmt"
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"
)

type forecastVariable int

const (
	x forecastVariable = iota
	y
	magnitude
	period
)

type forecastInfo struct {
	epiX       shared.Coordinate // x co-ord of disaster epicentre
	epiY       shared.Coordinate // y ""
	mag        shared.Magnitude
	period     uint
	confidence float64
}

type forecastHistory map[uint]forecastInfo                                       // stores history of past disasters
type receivedForecastHistory map[uint]shared.ReceivedDisasterPredictionsDict     // stores history of received disasters
type clientsForecastPerformance map[shared.ClientID]map[forecastVariable]float64 // stores current notion of prediction performance of other teams

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {

	lastDisasterTurn := c.disasterHistory.getLastDisasterTurn()

	fInfo, confMap, err := c.disasterModel.generateForecast(c.config)

	if err != nil {
		c.Logf("ERROR: unable to generate forecast. Encountered %v", err)
		// we can still proceed - fInfo will just be default with confidence zero
	} else {
		c.Logf("forecast: %+v, model support: %v, confidence map: %+v", fInfo, c.disasterModel.support, confMap)
	}

	prediction := shared.DisasterPrediction{
		CoordinateX: fInfo.epiX,
		CoordinateY: fInfo.epiY,
		Magnitude:   fInfo.mag,
		TimeLeft:    uint(lastDisasterTurn + fInfo.period - c.getTurn()),
	}

	trustedIslandIDs := []shared.ClientID{}
	trustThresh := c.config.forecastTrustTreshold
	for id := range c.getTrustedTeams(trustThresh, false, forecastingBasis) { // TODO: decide if this should be general or forecasting basis
		trustedIslandIDs = append(trustedIslandIDs, id)
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslandIDs,
	}
	c.lastDisasterPrediction = prediction
	// update forecast history
	c.forecastHistory[c.getTurn()] = fInfo
	return predictionInfo
}

func (d *disasterModel) generateForecast(conf clientConfig) (f forecastInfo, confMap map[forecastVariable]float64, err error) {
	nSamples := d.support

	if nSamples == 0 {
		return forecastInfo{}, map[forecastVariable]float64{}, errors.Errorf("Cannot generate forecast with no data")
	}
	magStats, errM := d.magnitude.getStatistics(nSamples)
	xStats, errX := d.x.getStatistics(nSamples)
	yStats, errY := d.y.getStatistics(nSamples)
	periodStats, errP := d.period.getStatistics(nSamples)

	for _, err := range []error{errM, errX, errY, errP} {
		if err != nil {
			return forecastInfo{}, map[forecastVariable]float64{}, errors.Errorf("Unable to generate forecast. First error encountered: %v", err)
		}
	}

	confidence, confMap := getWeightedConfidence(map[forecastVariable]modelStats{
		period:    periodStats,
		magnitude: magStats,
		x:         xStats,
		y:         yStats,
	}, conf)

	f = forecastInfo{
		epiX:       xStats.mean,
		epiY:       yStats.mean,
		mag:        magStats.mean,
		period:     uint(math.Round(periodStats.mean)),
		confidence: confidence,
	}
	return f, confMap, nil
}

// computes confidence combination of modelStats weighted by the perceived importance
// of each estimated quantity. For example, we may want to weight period confidence higher.
func getWeightedConfidence(paramStats map[forecastVariable]modelStats, config clientConfig) (float64, map[forecastVariable]float64) {

	weightsConf := config.forecastParamWeights

	weights := []float64{}
	confMap := map[forecastVariable]float64{}
	confidence := 0.0
	// note: these string keys should match those in config
	for param, stats := range paramStats {
		baseConf := stats.meanConfidence
		confidence += baseConf * weightsConf[param]
		confMap[param] = baseConf // store this for logging purposes
		weights = append(weights, weightsConf[param])
	}
	return confidence / floats.Sum(weights), confMap
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence

	if len(receivedPredictions) == 0 {
		c.Logf("[turn %v]: Nobody wanna share forecasts with us :((", c.getTurn())
		return
	}

	sumX, sumY, sumMag, sumConf := 0.0, 0.0, 0.0, 0.0
	sumTime := uint(0)

	c.updateForecastingReputations(receivedPredictions) // update our perceptions of other teams

	c.receivedForecastHistory[c.getTurn()] = receivedPredictions // update rxForecastsHistor

	//c.lastDisasterForecast.Confidence *= 1.3 // inflate confidence of our prediction above others
	receivedPredictions[c.GetID()] = shared.ReceivedDisasterPredictionInfo{PredictionMade: c.lastDisasterPrediction, SharedFrom: c.GetID()}

	// weight predictions by their confidence and our assessment of their forecasting reputation
	for rxTeam, pred := range receivedPredictions {
		rep := float64(c.opinions[rxTeam].getForecastingRep()) + 1 // our notion of another island's forecasting reputation
		sumX += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateX * rep
		sumY += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateY * rep
		sumMag += pred.PredictionMade.Confidence * pred.PredictionMade.Magnitude * rep
		sumTime += uint(pred.PredictionMade.Confidence * float64(pred.PredictionMade.TimeLeft) * rep)
		sumConf += pred.PredictionMade.Confidence * rep
	}

	sumConf = math.Max(sumConf, 1) // guard against div by zero error below
	// Finally get the final prediction generated by considering predictions from all islands that we have available
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: sumX / sumConf,
		CoordinateY: sumY / sumConf,
		Magnitude:   sumMag / sumConf,
		TimeLeft:    uint((float64(sumTime) / sumConf) + 0.5),    // +0.5 for rounding
		Confidence:  sumConf / float64(len(receivedPredictions)), // this len will always be >= 1
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}

func (c *client) updateForecastingReputations(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for team, predInfo := range receivedPredictions {
		// if teams make predictions with conf > 50% before first disaster, downgrade their rep by 75%
		if len(c.disasterHistory) == 0 {
			if predInfo.PredictionMade.Confidence > 50 {
				c.opinions[team].updateOpinion(forecastingBasis, c.changeOpinion(-0.75))
			}
		}
		// decrease trust in teams who are overly confident
		if predInfo.PredictionMade.Confidence > 98 {
			c.opinions[team].updateOpinion(forecastingBasis, c.changeOpinion(-0.3))
		}
		// note: more sophisticated updates happen in DisasterNotification()
	}
}

// This should be called after every disaster occurs. Evaluates historical performance of clients' forecasts
// compared to ours
func (c *client) evaluateForecastingPerformance() (map[shared.ClientID]map[forecastVariable]float64, error) {
	clientSkills := map[shared.ClientID]map[forecastVariable]float64{}

	if c.getTurn() != c.disasterHistory.getLastDisasterTurn() {
		return clientSkills, errors.Errorf("Turn of most recent disaster does not match current turn")
	}
	ourForecastErrors, errF := computeForecastingPerformance(c.disasterHistory, c.forecastHistory, c.config)
	ourErr, errP := c.aggregateForecastingError(ourForecastErrors, true)

	if errF != nil || errP != nil {
		return clientSkills, errors.Errorf("Encountered error while computing our forecasting performance: forecast err: %v, performance err: %v", errF, errP)
	}
	clientForecasts := map[shared.ClientID]forecastHistory{}
	clientErrors := map[shared.ClientID][]map[forecastVariable]float64{}

	// collect history of client forecasts
	for turn, forecastMap := range c.receivedForecastHistory {
		for client, predInfo := range forecastMap {
			clientForecasts[client] = forecastHistory{}
			clientForecasts[client][turn] = c.parsePredictionInfo(predInfo.PredictionMade)
		}
	}

	for cID, fHist := range clientForecasts {
		errorMap, err := computeForecastingPerformance(c.disasterHistory, fHist, c.config)

		if err != nil {
			return clientSkills, errors.Errorf("Encountered error while computing client's forecasting errors. Client ID: %v, received forecast history: %v", cID, fHist)
		}
		clientErrors[cID] = errorMap
	}

	for cID, clientErrMaps := range clientErrors {
		clientErr, err := c.aggregateForecastingError(clientErrMaps, true)
		if err != nil {
			return clientSkills, errors.Errorf("Encountered error while computing client's forecasting performance. Client ID: %v, received error maps: %v", cID, clientErrMaps)
		}
		clientSkills[cID] = map[forecastVariable]float64{} // initialise in memory
		for k, val := range clientErr {
			if ourErr[k] == 0 { // we're perfect at forecasting this variable - possible, but unlikely
				clientSkills[cID][k] = 0 // client skills are zero relative to us as we're infinitely good
			}
			clientSkills[cID][k] = absoluteCap(1-(val/ourErr[k]), 1) // 0 => on par with us. < 0 => worse than us (larger error). > 0 => better than us (smalle error).
		}
	}
	return clientSkills, nil
}

// function to perform exp-weighted average of historical forecast errors
func (c client) aggregateForecastingError(errorMaps []map[forecastVariable]float64, useExpWeighting bool) (map[forecastVariable]float64, error) {
	n := len(errorMaps)
	result := map[forecastVariable]float64{}

	if n == 0 {
		return result, errors.Errorf("Cannot create aggregated perfromance with no samples (len = 0)")
	}
	decay := 1.0
	if useExpWeighting {
		decay = c.config.forecastTemporalDecay
	}
	weights := make([]float64, n)
	for i, errorMap := range errorMaps {
		for k, val := range errorMap {
			w := math.Pow(decay, float64(n-i-1)) // exponential weighting (decay = 1 has no effect)
			result[k] += val * w
			weights[i] = w
		}
	}
	for k := range result {
		result[k] /= math.Max(floats.Sum(weights), 1) // sum of weights will be > 1 if len > 0 since decay^0 = 1
	}
	return result, nil
}

func (c client) parsePredictionInfo(p shared.DisasterPrediction) forecastInfo {
	period := c.getTurn() - c.disasterHistory.getLastDisasterTurn() + p.TimeLeft
	return forecastInfo{
		epiX:   p.CoordinateX,
		epiY:   p.CoordinateY,
		mag:    p.Magnitude,
		period: period,
	}
}

func computeForecastingPerformance(dh disasterHistory, fh forecastHistory, conf clientConfig) (forecastErrors []map[forecastVariable]float64, err error) {
	disasterPeriodHistory := dh.getPastDisasterPeriods()
	forecastTurns := uintsAsFloats(fh.sortKeys())
	prevTurn := 0.0
	forecastErrors = []map[forecastVariable]float64{}
	i := 0 // index variable
	for turn, report := range dh {
		indexes, _ := floats.Find([]int{}, func(x float64) bool {
			return (x <= float64(turn)) && (x > prevTurn)
		}, forecastTurns, -1)
		prevTurn = float64(turn)

		precedingForecasts := make([]forecastInfo, len(indexes))
		for i, index := range indexes {
			precedingForecasts[i] = fh[uint(forecastTurns[index])]
		}
		forecastErrors = append(forecastErrors, analyseForecastSkill(precedingForecasts, report, float64(disasterPeriodHistory[i]), conf))
		i++
	}
	return forecastErrors, nil
}

func analyseForecastSkill(forecasts []forecastInfo, disaster disasterInfo, disasterPeriod float64, conf clientConfig) (mseMap map[forecastVariable]float64) {
	mseMap = map[forecastVariable]float64{}
	pVals, xVals, yVals, magVals := []float64{}, []float64{}, []float64{}, []float64{}
	for _, f := range forecasts {
		pVals = append(pVals, float64(f.period))
		xVals = append(xVals, f.epiX)
		yVals = append(yVals, f.epiY)
		magVals = append(magVals, f.mag)
	}

	decay := conf.forecastTemporalDecay

	mseMap[period] = univariateWeightedMSE(pVals, disasterPeriod, decay)
	mseMap[magnitude] = univariateWeightedMSE(magVals, disaster.report.Magnitude, decay)
	mseMap[x] = univariateWeightedMSE(xVals, disaster.report.X, decay)
	mseMap[y] = univariateWeightedMSE(yVals, disaster.report.Y, decay)

	return mseMap
}

// exponential weighting to series of floats (sorted chronologically so latest is late). Usually,
// decay in [0; 1] so that latest values are weighted more.
func expWeighting(x []float64, decay float64) (weightedX, weights []float64) {
	for i, el := range x {
		w := math.Pow(decay, float64(i))
		weights = append(weights, w)
		weightedX = append(weightedX, w*el)
	}
	return weightedX, weights
}

// mean squared error for a exponentially weighted sample compared to a target value
func univariateWeightedMSE(sample []float64, target, expDecay float64) (MSE float64) {
	weightedSample, _ := expWeighting(sample, expDecay)
	for _, s := range weightedSample {
		MSE += math.Pow(s-target, 2)
	}
	return MSE / math.Max(float64(len(weightedSample)), 1)
}

func (f forecastVariable) String() string {
	strings := [...]string{"x", "y", "magnitude", "period"}
	if f >= 0 && int(f) < len(strings) {
		return strings[f]
	}
	return fmt.Sprintf("UNKNOWN forecast variable '%v'", int(f))
}

// GoString implements GoStringer
func (f forecastVariable) GoString() string {
	return f.String()
}

// MarshalText implements TextMarshaler
func (f forecastVariable) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(f.String())
}

// MarshalJSON implements RawMessage
func (f forecastVariable) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(f.String())
}

func (fh forecastHistory) sortKeys() []uint {
	keys := make([]int, 0)
	for k := range fh {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	finalKeys := make([]uint, 0)
	for _, k := range keys {
		finalKeys = append(finalKeys, uint(k))
	}
	return finalKeys
}
