package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
)

type mockClientEcho struct {
	common.Client
	id   shared.ClientID
	echo string
}

func (c *mockClientEcho) GetID() shared.ClientID {
	return c.id
}

func (c *mockClientEcho) Echo(s string) string {
	return c.echo
}

// TestGetEcho also exercises getClientInfoFromRegisteredClients
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
				clientMap: map[shared.ClientID]common.Client{
					shared.Team1: mClient,
				},
			}

			got := server.getEcho(tc.input)
			testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}
