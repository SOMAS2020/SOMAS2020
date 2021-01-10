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

type forecastHistory map[uint]forecastInfo                                   // stores history of past disasters
type receivedForecastHistory map[uint]shared.ReceivedDisasterPredictionsDict // stores history of received disasters

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

func analyseDisasterHistory(dh disasterHistory, fh forecastHistory) {
	forecastTurns := uintsAsFloats(fh.sortKeys())
	prevTurn := 0.0
	for turn := range dh {
		indexes, _ := floats.Find([]int{}, func(x float64) bool {
			return (x <= float64(turn)) && (x > prevTurn)
		}, forecastTurns, -1)
		prevTurn = float64(turn)
		fmt.Printf("Indexes: %v", indexes)
	}
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
