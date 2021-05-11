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
		2,
		4,
		6,
		1,
		3,
		5,
	}
	rand.Shuffle(len(clients), func(i, j int) { clients[i], clients[j] = clients[j], clients[i] })

	want := []ClientID{
		1,
		2,
		3,
		4,
		5,
		6,
	}

	sort.Sort(SortClientByID(clients))
	if !reflect.DeepEqual(want, clients) {
		t.Errorf("want '%v' got '%v'", want, clients)
	}
}
