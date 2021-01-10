package team5

import (
	"testing"

	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestCommonPoolResourceRequest(t *testing.T) {
	// test submit request to president when we are poor & current CP has nothing
	var turn, season uint = 2, 1
	var currentCP shared.Resources
	currentCP = 2
	c.cpAllocationHistory[1] = 20
	c.cpAllocationHistory[2] = 20
	currentTier := imperialStudent
	reqAmount := c.calculateRequestToPresident(turn, season, currentTier, currentCP)
	w := c.config.imperialThreshold
	if roundTo(float64(w), 1) != roundTo(float64(reqAmount), 1) {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, roundTo(float64(reqAmount), 2))
	}

	// test submit request to president when we are poor and current CP has a lot
	currentCP = 1000
	reqAmount = c.calculateRequestToPresident(turn, season, currentTier, currentCP)
	w = currentCP / 6
	if w != reqAmount {
		t.Errorf("Not requesting proper # of resources to president. Want %v, got %v", w, reqAmount)
	}

	// test submit request to president when we are normal and have histories
	currentTier = middleClass
	reqAmount = c.calculateRequestToPresident(turn, season, currentTier, currentCP)
	currentTier = middleClass
	w = c.cpAllocationHistory[2] + 10
	if w != reqAmount {
		t.Errorf("Not requesting proper # of resources to president. Want %v, got %v", w, reqAmount)
	}

	// test submit request to president in turn 1, season 1
	turn = 1
	reqAmount = c.calculateRequestToPresident(turn, season, currentTier, currentCP)
	w = c.config.imperialThreshold
	if w != reqAmount {
		t.Errorf("Not requesting proper # of resources to president. Want %v, got %v", w, reqAmount)
	}
}

func TestRequestAllocation(t *testing.T) {
	// test request allocation when we are poor & we don't have any allocation from President
	//	there is nothing in commonpool
	currentTier := imperialStudent
	var currentCP shared.Resources
	currentCP = 1
	allocationAmount := shared.Resources(0)
	allocationMade := true
	allocation := c.calculateAllocationFromCP(currentTier, currentCP, allocationMade, allocationAmount)

	w := c.config.imperialThreshold
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when we are poor & we have allocation from President
	// 	there is nothing in commonpool
	allocationAmount = 50
	w = 50
	allocation = c.calculateAllocationFromCP(currentTier, currentCP, allocationMade, allocationAmount)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when there is a lot of resource in the common pool and we are poor
	currentCP = 10000
	w = currentCP / 6
	allocation = c.calculateAllocationFromCP(currentTier, currentCP, allocationMade, allocationAmount)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when we are in normal state and there's allocation
	currentTier = middleClass
	allocationAmount = 50
	w = 50
	allocation = c.calculateAllocationFromCP(currentTier, currentCP, allocationMade, allocationAmount)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}
}

// for both tax and other contribution
func TestGetCommonPoolContribution(t *testing.T) {
	// day 1 season 1 case
	var turn, season uint = 1, 1
	currentTier := middleClass
	taxAmount := shared.Resources(20)
	var w shared.Resources

	tax := calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution := c.calculateCPContribution(turn, season, currentTier)
	total := tax + contribution
	w = 20
	if w != total {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}

	//when we are poor
	currentTier = imperialStudent
	turn = 2
	tax = calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season, currentTier)
	total = tax + contribution
	w = 0
	if w != total {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, total)
	}

	//when we are normal and use histories for contribution
	//when the cashflow of CP is positive
	currentTier = middleClass
	turn = 1
	season = 2
	c.cpResourceHistory[1] = 20
	c.cpResourceHistory[2] = 40
	//c.Logf("differen2: %v", c.cpResourceHistory[1])
	tax = calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season, currentTier)
	total = tax + contribution
	w = 20/6 + 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}

	//when the cashflow of CP is negative
	turn = 2
	c.cpResourceHistory[1] = 60
	tax = calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season, currentTier)
	total = tax + contribution
	w = 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}
}

func TestCalculateDisasterContributionCP(t *testing.T) {
	ourLocationInfo := disasters.IslandLocationInfo{
		ID: shared.Team5,
		X:  6,
		Y:  0,
	}
	geography := disasters.ArchipelagoGeography{
		Islands: make(map[shared.ClientID]disasters.IslandLocationInfo),
	}
	geography.Islands[shared.Team5] = ourLocationInfo
	//case 1 where there is no forecast initially
	currentTurn := uint(1)
	currentResource := shared.Resources(1000)
	lastDisaster := uint(0)
	contribution := c.calculateDisasterContributionCP(currentTurn, currentResource, geography, lastDisaster)
	w := shared.Resources(0)
	if w != contribution {
		t.Errorf("1. Not generating proper # of resources to mitigate disaster. Want %v, got %v", w, contribution)
	}

	// case 2 when there is forecast. But it's in 9 days so we don't contribute
	forecast := forecastInfo{
		epiX:       shared.Coordinate(3),
		epiY:       shared.Coordinate(4),
		mag:        10,
		period:     10,
		confidence: 0.6,
	}
	c.forecastHistory[1] = forecast
	contribution = c.calculateDisasterContributionCP(currentTurn, currentResource, geography, lastDisaster)
	if w != contribution {
		t.Errorf("2. Not generating proper # of resources to mitigate disaster. Want %v, got %v", w, contribution)
	}

	// case 3 when there is forecast, and it's in 1 days so we contribute 20% of the calculated damage
	// and we are quite rich
	currentTurn = 9
	c.forecastHistory[8] = forecast
	distance := math.Sqrt(math.Pow(ourLocationInfo.X-forecast.epiX, 2) + math.Pow(ourLocationInfo.Y-forecast.epiY, 2))
	c.Logf("test distance: %v forecast mag: %v", distance, forecast.mag)
	w = shared.Resources((float64(forecast.mag) * 100 / distance) * 0.2)
	contribution = c.calculateDisasterContributionCP(currentTurn, currentResource, geography, lastDisaster)
	if w != contribution {
		t.Errorf("3. Not generating proper # of resources to mitigate disaster. Want %v, got %v", w, contribution)
	}

	// case 4 when there is forecast, it's in 1 days and we are not that rich
	currentResource = 3
	w = currentResource / 2
	contribution = c.calculateDisasterContributionCP(currentTurn, currentResource, geography, lastDisaster)
	if w != contribution {
		t.Errorf("4. Not generating proper # of resources to mitigate disaster. Want %v, got %v", w, contribution)
	}

	// case 5 when there is forecast, it's today and we are rich
	c.forecastHistory[9] = forecast
	currentResource = 1000
	currentTurn = 10
	lastDisaster = 0
	w = shared.Resources((float64(forecast.mag) * 100 / distance) * 0.8)
	contribution = c.calculateDisasterContributionCP(currentTurn, currentResource, geography, lastDisaster)
	if w != contribution {
		t.Errorf("5. Not generating proper # of resources to mitigate disaster. Want %v, got %v", w, contribution)
	}
}
