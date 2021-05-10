package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team1"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team2"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team3"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team4"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team5"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team6"
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type ClientFactory func(shared.ClientID) baseclient.Client

func DefaultClientConfig() map[shared.ClientID]ClientFactory {
	return map[shared.ClientID]ClientFactory{
		shared.Team1: team1.DefaultClient,
		shared.Team2: team1.DefaultClient,
		shared.Team3: team1.DefaultClient,
		shared.Team4: team1.DefaultClient,
		shared.Team5: team1.DefaultClient,
		shared.Team6: team1.DefaultClient,
		shared.Team7: team1.DefaultClient,
		shared.Team8: team1.DefaultClient,
		// shared.Team9: team6.DefaultClient,
		// shared.Team10: team4.DefaultClient,
		// shared.Team11: team1.DefaultClient,
		// shared.Team12: team1.DefaultClient,
	}
}
