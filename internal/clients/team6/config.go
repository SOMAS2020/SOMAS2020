package team6

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// Config configures our island's initial state
type Config struct {
	friendship Friendship
}

var configInfo = Config{
	friendship: Friendship{
		shared.Team1: 50,
		shared.Team2: 50,
		shared.Team3: 50,
		shared.Team4: 50,
		shared.Team5: 50,
	},
}
