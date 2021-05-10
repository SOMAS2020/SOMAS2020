package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team1"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team2"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team3"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team4"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team5"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team6"
	"strconv"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type ClientFactory func(shared.ClientID) baseclient.Client

// func DefaultClientConfig() map[shared.ClientID]ClientFactory {
// 	return map[shared.ClientID]ClientFactory{
// 		shared.Teams["Team1"]: team1.DefaultClient,
// 		shared.Teams["Team2"]: team1.DefaultClient,
// 		shared.Teams["Team3"]: team1.DefaultClient,
// 		shared.Teams["Team4"]: team1.DefaultClient,
// 		shared.Teams["Team5"]: team1.DefaultClient,
// 		shared.Teams["Team6"]: team1.DefaultClient,
// 		// shared.Team7: team1.DefaultClient,
// 		// shared.Team8: team1.DefaultClient,
// 		// shared.Team9: team1.DefaultClient,
// 		// shared.Team10: team1.DefaultClient,
// 		// shared.Team11: team1.DefaultClient,
// 		// shared.Team12: team1.DefaultClient,
// 	}
// }

func DefaultClientConfig() map[shared.ClientID]ClientFactory {
	clientMapping := make(map[shared.ClientID]ClientFactory)

	for i := 0; i < shared.TotalTeams; i++ {
		clientMapping[shared.Teams["Team"+strconv.Itoa(i+1)]] = team1.DefaultClient
	}

	return clientMapping
}
