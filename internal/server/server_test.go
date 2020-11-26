package server

import (
	"fmt"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
)

// a factory function that allows us to inject custom stuff to facilitate testing
func somasServerTester(clients []common.Client, gameState common.GameState) Server {
	return &SOMASServer{
		clients:   clients,
		gameState: gameState,
	}
}

// extend as required
type mockClient struct {
	common.Client
	reply string
	id    int
}

func (c *mockClient) Echo(s string) string {
	return c.reply
}

func (c *mockClient) GetID() int {
	return c.id
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
			want:  fmt.Errorf("Echo error: want '42' got '43' from client 1"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mClient := &mockClient{
				id:    1,
				reply: tc.reply,
			}
			clients := []common.Client{mClient}
			server := somasServerTester(clients, common.GameState{})
			got := server.GetEcho(tc.input)
			testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}
