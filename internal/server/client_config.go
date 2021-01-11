package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team1"
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team2"
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3"
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team4"
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team5"
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team6"
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func DefaultClientMap() map[shared.ClientID]baseclient.Client {
	return map[shared.ClientID]baseclient.Client{
		shared.Team1: team1.DefaultClient(shared.Team1),
		shared.Team2: team2.DefaultClient(shared.Team2),
		shared.Team3: team3.DefaultClient(shared.Team3),
		shared.Team4: team4.DefaultClient(shared.Team4),
		shared.Team5: team5.DefaultClient(shared.Team5),
		shared.Team6: team6.DefaultClient(shared.Team6),
	}
}
