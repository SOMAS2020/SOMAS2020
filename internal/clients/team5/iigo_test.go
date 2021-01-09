package team5

import (
	"math"
	"testing"

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
	if math.Abs(float64(w-reqAmount)) > 0.1 { //threshold is 0.1
		t.Errorf("Not requesting proper # of resources to president. Want %v, got %v", w, reqAmount)
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
	contribution := c.calculateCPContribution(turn, season)
	total := tax + contribution
	w = 20
	if w != total {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}

	//when we are poor
	currentTier = imperialStudent
	turn = 2
	tax = calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season)
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
	contribution = c.calculateCPContribution(turn, season)
	total = tax + contribution
	w = 20/6 + 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}

	//when the cashflow of CP is negative
	turn = 2
	c.cpResourceHistory[1] = 60
	tax = calculateTaxContribution(taxAmount, turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season)
	total = tax + contribution
	w = 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}
}
