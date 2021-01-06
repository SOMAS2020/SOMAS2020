package team5

import (
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
	reqAmount := c.calculateRequestAllocation(turn, season, currentTier, currentCP)
	w := c.config.imperialThreshold
	if w != reqAmount {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, reqAmount)
	}

	// test submit request to president when we are poor and current CP has a lot
	currentCP = 1000
	reqAmount = c.calculateRequestAllocation(turn, season, currentTier, currentCP)
	w = currentCP / 6
	if w != reqAmount {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, reqAmount)
	}

	// test submit request to president when we are normal and have histories
	currentTier = middleClass
	reqAmount = c.calculateRequestAllocation(turn, season, currentTier, currentCP)
	currentTier = middleClass
	w = c.cpAllocationHistory[2] + 10
	if w != reqAmount {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, reqAmount)
	}

	// test submit request to president in turn 1, season 1
	turn = 1
	reqAmount = c.calculateRequestAllocation(turn, season, currentTier, currentCP)
	w = c.config.imperialThreshold
	if w != reqAmount {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, reqAmount)
	}
}

func TestRequestAllocation(t *testing.T) {
	// test request allocation when we are poor & we don't have any allocation from President
	//	there is nothing in commonpool
	currentTier := imperialStudent
	var currentCP shared.Resources
	currentCP = 1
	c.allocation = 0
	allocation := c.calculateRequestCommonPool(currentTier, currentCP)

	w := c.config.imperialThreshold
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when we are poor & we have allocation from President
	// 	there is nothing in commonpool
	c.allocation = 50
	w = 50
	allocation = c.calculateRequestCommonPool(currentTier, currentCP)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when there is a lot of resource in the common pool and we are poor
	currentCP = 10000
	w = currentCP / 6
	allocation = c.calculateRequestCommonPool(currentTier, currentCP)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}

	//test request allocation when we are in normal state and there's allocation
	currentTier = middleClass
	c.allocation = 50
	w = 50
	allocation = c.calculateRequestCommonPool(currentTier, currentCP)
	if w != allocation {
		t.Errorf("Not taking proper # of resources from cp. Want %v, got %v", w, allocation)
	}
}

// for both tax and other contribution
func TestGetCommonPoolContribution(t *testing.T) {
	// day 1 season 1 case
	var turn, season uint = 1, 1
	currentTier := imperialStudent
	c.taxAmount = 20
	var w shared.Resources

	tax := c.calculateTaxContribution(turn, season, currentTier)
	contribution := c.calculateCPContribution(turn, season)
	total := tax + contribution
	w = 10
	if w != total {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, total)
	}

	//when we are poor
	turn = 2
	tax = c.calculateTaxContribution(turn, season, currentTier)
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
	tax = c.calculateTaxContribution(turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season)
	total = tax + contribution
	w = 20/6 + 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}

	//when the cashflow of CP is negative
	turn = 2
	c.cpResourceHistory[1] = 60
	tax = c.calculateTaxContribution(turn, season, currentTier)
	contribution = c.calculateCPContribution(turn, season)
	total = tax + contribution
	w = 20
	if int(w) != int(total) {
		t.Errorf("Not generating proper # of resources to contribute. Want %v, got %v", w, tax)
	}
}
