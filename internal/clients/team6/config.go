package team6

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// Config configures our island's initial state
type Config struct {
	friendshipWith FriendshipWith
}

var clientConfigInfo = Config{
	friendshipWith: FriendshipWith{
		shared.Team1: neutral,
		shared.Team2: neutral,
		shared.Team3: neutral,
		shared.Team4: neutral,
		shared.Team5: neutral,
	},
}
