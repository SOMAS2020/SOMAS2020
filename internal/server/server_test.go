package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
)

type mockClientEcho struct {
	baseclient.Client
	id   shared.ClientID
	echo string
}

func (c *mockClientEcho) GetID() shared.ClientID {
	return c.id
}

func (c *mockClientEcho) Echo(s string) string {
	return c.echo
}

func TestGetEcho(t *testing.T) {
	cases := []struct {
		name  string
		input string
		reply string
		want  error
	}{
		{
			name:  "basic ok",
			input: "42",
			reply: "42",
			want:  nil,
		},
		{
			name:  "wrong reply",
			input: "42",
			reply: "43",
			want:  errors.Errorf("Echo error: want '42' got '43' from Team1"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mClient := &mockClientEcho{
				id:   shared.Team1,
				echo: tc.reply,
			}
			server := &SOMASServer{
				clientMap: map[shared.ClientID]baseclient.Client{
					shared.Team1: mClient,
				},
			}

			got := server.getEcho(tc.input)
			testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}

type initTestClient struct {
	baseclient.Client
	initialiseCalled bool
}

func (c *initTestClient) Initialise(baseclient.ServerReadHandle) {
	c.initialiseCalled = true
}

func TestSOMASServerFactoryInitialisesClients(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{}

	// We need to initialise a map of *initTestClient (not baseclient.Client)
	// because we need access to initialiseCalled. We can then convert this to a
	// map of baseClient.Client, just to pass to createSOMASServer.
	clientPtrsMap := map[shared.ClientID]*initTestClient{
		shared.Team1: {Client: baseclient.NewClient(shared.Team1)},
		shared.Team2: {Client: baseclient.NewClient(shared.Team2)},
		shared.Team3: {Client: baseclient.NewClient(shared.Team3)},
	}

	clientMap := map[shared.ClientID]baseclient.Client{}
	for k, v := range clientPtrsMap {
		clientMap[k] = v
	}

	_, err := createSOMASServer(clientInfos, clientMap, config.Config{})
	if err != nil {
		t.Error(err)
	}

	for clientID, client := range clientPtrsMap {
		if !client.initialiseCalled {
			t.Errorf("Initialise not called for %v", clientID)
		}
	}
}

func lstHasUniqueClientIDs(lst []shared.ClientID) bool {
	// set to contain what we've seen so far
	var s map[shared.ClientID]interface{}
	for _, id := range lst {
		if _, ok := s[id]; ok {
			return false
		}
		s[id] = nil
	}
	return true
}

func TestGetNRandClientIDsUniqueIfPossible(t *testing.T) {
	iterations := 10
	cases := []struct {
		name      string
		input     []shared.ClientID
		retLength int
	}{
		{
			name:      "RandomAssign 2 entries",
			input:     []shared.ClientID{shared.ClientID(0), shared.ClientID(1)},
			retLength: 3,
		},
		{
			name:      "RandomAssign 3 entries",
			input:     []shared.ClientID{shared.ClientID(0), shared.ClientID(1), shared.ClientID(2)},
			retLength: 3,
		},
		{
			name: "RandomAssign several entries",
			input: []shared.ClientID{
				shared.ClientID(0),
				shared.ClientID(1),
				shared.ClientID(2),
				shared.ClientID(3),
				shared.ClientID(4),
				shared.ClientID(5)},
			retLength: 3,
		},
		{
			name: "RandomAssign long return",
			input: []shared.ClientID{
				shared.ClientID(0),
				shared.ClientID(1),
				shared.ClientID(2),
				shared.ClientID(3),
				shared.ClientID(4),
				shared.ClientID(5)},
			retLength: 5,
		},
		{
			name:      "RandomAssign one entry, but long return lsit",
			input:     []shared.ClientID{shared.ClientID(0), shared.ClientID(1)},
			retLength: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.input) < tc.retLength {
				lst, _ := getNRandClientIDsUniqueIfPossible(tc.input, tc.retLength) // Only check for crash
				if len(lst) != tc.retLength {
					t.Errorf("%v - Return list length %v, different from expected length %v", tc.name, len(lst), tc.retLength)
				}
			} else {
				for i := 0; i < iterations; i++ { // As its using random numbers. Run each test several times to minimise probability
					lst, err := getNRandClientIDsUniqueIfPossible(tc.input, tc.retLength)
					if len(lst) != tc.retLength {
						t.Errorf("%v - Return list length %v, different from expected length %v", tc.name, len(lst), tc.retLength)
					}
					if err != nil && !lstHasUniqueClientIDs(lst) {
						t.Errorf("%v - Elements in the return list %v should all be unique but weren't", tc.name, lst)
					}
				}
			}
		})

	}

}
