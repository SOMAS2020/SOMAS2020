package shared

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestSortClientByID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	clients := []ClientID{
		Team2,
		Team4,
		Team6,
		Team1,
		Team3,
		Team5,
	}
	rand.Shuffle(len(clients), func(i, j int) { clients[i], clients[j] = clients[j], clients[i] })

	want := []ClientID{
		Team1,
		Team2,
		Team3,
		Team4,
		Team5,
		Team6,
	}

	sort.Sort(SortClientByID(clients))
	if !reflect.DeepEqual(want, clients) {
		t.Errorf("want '%v' got '%v'", want, clients)
	}
}
