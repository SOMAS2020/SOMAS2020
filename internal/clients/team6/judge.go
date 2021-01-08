package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	*client
}

func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	pardons := map[int][]bool{}
	maxSanctionTime := int(j.client.ServerReadHandle.GetGameConfig().IIGOClientConfig.SanctionCacheDepth - 1)

	for timeStep, sanctions := range currentSanctions {
		for who, sanction := range sanctions {
			if timeStep == maxSanctionTime && j.client.friendship[sanction.ClientID] == j.client.clientConfig.maxFriendship {
				// we can pardon certain islands having maximum friendship with us
				pardons[timeStep][who] = true
			}
		}
	}

	return pardons
}
