package team4

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

type trust struct {
	trustMap map[shared.ClientID]float64
}

func (t *trust) GetTrustMatrix() []float64 {
	trustMatrix := make([]float64, len(t.trustMap))
	for clientID, trustValue := range t.trustMap {
		index := int(clientID)
		if index < len(trustMatrix) {
			trustMatrix[index] = trustValue
		}
	}
	return trustMatrix
}

func (t *trust) GetClientTrust(clientID shared.ClientID) float64 {

	if clientTrust, ok := t.trustMap[clientID]; ok {
		return clientTrust
	}
	return 0
}

func (t *trust) ChangeClientTrust(clientID shared.ClientID, diff float64) {
	// diff is percentage to change trust for clientID ie. diff in range [-1,1]
	if _, ok := t.trustMap[clientID]; ok {
		t.trustMap[clientID] = t.trustMap[clientID] * (1 + diff)
		t.normalise()
	}
}

func (t *trust) SetClientTrust(clientID shared.ClientID, newValue float64) {
	// diff is percentage to change trust for clientID ie. diff in range [-1,1]
	if _, ok := t.trustMap[clientID]; ok {
		t.trustMap[clientID] = newValue
		t.normalise()
	}
}

func (t *trust) totalTrustSum() float64 {
	totalTrust := 0.0
	for _, trust := range t.trustMap {
		totalTrust += trust
	}
	return totalTrust
}

func (t *trust) expectedTrustSum() float64 {
	return 0.5 * float64(len(t.trustMap))
}

//normalise ensures that trust values are always in range [0,1]
func (t *trust) normalise() {
	// it ensures that the general trust sums to 0.5 * number of clients
	if len(t.trustMap) > 0 && t.totalTrustSum() > 0 {
		normaliseCoef := t.expectedTrustSum() / t.totalTrustSum()
		for clientID, trust := range t.trustMap {
			t.trustMap[clientID] = trust * normaliseCoef
		}
	}
}

func (t *trust) initialise() {
	for _, clientID := range shared.TeamIDs {
		t.trustMap[clientID] = 0.5
	}
	t.normalise()
}

//Return a list of clients above a trust threshold
func (t *trust) trustedClients(threshold float64) []shared.ClientID {
	var lst []shared.ClientID
	for client, val := range t.trustMap {
		if val > threshold {
			lst = append(lst, client)
		}
	}
	return lst
}
