package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) getTurn() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Turn
	}
	return 0
}

func (c *client) getSeason() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Season
	}
	return 0
}

func (c *client) getTurnLength(role shared.Role) uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameConfig().IIGOClientConfig.IIGOTermLengths[role]
	}
	return 0
}

func buildHistoryInfo(pairs []rules.VariableValuePair) (retInfo judgeHistoryInfo, ok bool) {
	resourceOK := 0
	taxOK := 0
	allocationOK := 0
	for _, val := range pairs {
		switch val.VariableName {
		case rules.IslandActualPrivateResources:
			if len(val.Values) > 0 {
				retInfo.Resources.expected = shared.Resources(val.Values[0])
				resourceOK++
			}
		case rules.IslandReportedPrivateResources:
			if len(val.Values) > 0 {
				retInfo.Resources.actual = shared.Resources(val.Values[0])
				resourceOK++
			}
		case rules.ExpectedTaxContribution:
			if len(val.Values) > 0 {
				retInfo.Taxation.expected = shared.Resources(val.Values[0])
				taxOK++
			}
		case rules.IslandTaxContribution:
			if len(val.Values) > 0 {
				retInfo.Taxation.actual = shared.Resources(val.Values[0])
				taxOK++
			}
		case rules.ExpectedAllocation:
			if len(val.Values) > 0 {
				retInfo.Allocation.expected = shared.Resources(val.Values[0])
				allocationOK++
			}
		case rules.IslandAllocation:
			if len(val.Values) > 0 {
				retInfo.Allocation.actual = shared.Resources(val.Values[0])
				allocationOK++
			}
		}
	}

	ok = resourceOK == 2 && taxOK == 2 && allocationOK == 2

	return retInfo, ok
}

func boolToFloat(input bool) float64 {
	if input {
		return 1
	}
	return 0
}
