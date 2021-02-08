package team5

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

// stores all information pertaining to a disaster
type disasterInfo struct {
	report  disasters.DisasterReport
	effects disasters.DisasterEffects
	season  uint
}

// stores learnt PDFs of various quantities using KDE (kernel density estimate) models
type disasterModel struct {
	period, magnitude, x, y kdeModel
	support                 uint // number of observations
}

type statistic int

const (
	mean statistic = iota
	variance
)

type disasterStatistics map[statistic]disasterInfo

func (d *disasterModel) updateModel(dR disasters.DisasterReport, period uint) {
	d.period.updateModel([]float64{float64(period)})
	d.magnitude.updateModel([]float64{dR.Magnitude})
	d.x.updateModel([]float64{dR.X})
	d.y.updateModel([]float64{dR.Y})

	d.support++ // one more observation
}

type disasterHistory map[uint]disasterInfo

// effects contain abs magnitude, prop. mag relative to other islands and CP mitigated mag.
func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.Logf("CRITICAL: Received notification of disaster: %v", dR.Display())
	period := c.getTurn() - c.disasterHistory.getLastDisasterTurn()
	c.disasterHistory[c.getTurn()] = disasterInfo{
		report:  dR,
		effects: effects,
		season:  c.getSeason(),
	}
	c.disasterModel.updateModel(dR, period)

	updatedPerf, err := c.evaluateForecastingPerformance()
	if err != nil {
		c.Logf("Encountered an error when evaluating forecasting performance: %v", err)
	} else {
		c.clientsForecastSkill = updatedPerf
		c.Logf("Updated our record of other agents' forecasting skill: %+v", updatedPerf)

		// now, update forecasting reputation of other teams based on their performance in forecasting last disaster
		for cID, perfMap := range updatedPerf {
			valSum := 0.0
			for _, val := range perfMap {
				valSum += val
			}
			meanPerf := valSum / float64(len(perfMap))                                     // this len is always > 0
			c.opinions[cID].updateOpinion(forecastingBasis, meanPerf*c.changeOpinion(0.4)) // 0.4 to control size of update
		}
	}
}

func (d disasterHistory) sortKeys() []uint {
	keys := make([]int, 0)
	for k := range d {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	finalKeys := make([]uint, 0)
	for _, k := range keys {
		finalKeys = append(finalKeys, uint(k))
	}
	return finalKeys
}

func (d disasterHistory) getLastDisasterTurn() uint {
	sortedTurns := d.sortKeys()
	l := len(sortedTurns)
	if l > 0 {
		return sortedTurns[l-1]
	}
	return 0
}

func (d disasterHistory) getPastDisasterPeriods() []float64 {
	periods := []float64{} // use float so we can use stat.Variance() later
	prevTurn := float64(startTurn)
	for _, turn := range d.sortKeys() { // TODO: find instances where assumption of ordered map keys is relied upon
		periods = append(periods, float64(turn)-prevTurn) // period = no. turns between successive disasters
		prevTurn = float64(turn)
	}
	return periods
}
