package team1

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestSortClientByID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	clients := sortByOpinion{
		opinionOnTeam{clientID: shared.Team2, opinion: 5},
		opinionOnTeam{clientID: shared.Team4, opinion: -5},
		opinionOnTeam{clientID: shared.Team6, opinion: 2},
		opinionOnTeam{clientID: shared.Team1, opinion: -1},
		opinionOnTeam{clientID: shared.Team3, opinion: -10},
		opinionOnTeam{clientID: shared.Team5, opinion: 0},
	}
	rand.Shuffle(len(clients), func(i, j int) { clients[i], clients[j] = clients[j], clients[i] })

	want := sortByOpinion{
		opinionOnTeam{clientID: shared.Team2, opinion: 5},
		opinionOnTeam{clientID: shared.Team6, opinion: 2},
		opinionOnTeam{clientID: shared.Team5, opinion: 0},
		opinionOnTeam{clientID: shared.Team1, opinion: -1},
		opinionOnTeam{clientID: shared.Team4, opinion: -5},
		opinionOnTeam{clientID: shared.Team3, opinion: -10},
	}

	sort.Sort(sortByOpinion(clients))
	if !reflect.DeepEqual(want, clients) {
		t.Errorf("want '%v' got '%v'", want, clients)
	}
}
