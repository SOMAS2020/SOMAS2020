package fileutils

import (
	"os"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
)

func TestPathExists(t *testing.T) {
	cases := []struct {
		name           string
		osStatErr      error
		logContainsSeq []string
		want           bool
	}{
		{
			name:      "exists",
			osStatErr: nil,
			want:      true,
		},
		{
			name:      "surely does not exists",
			osStatErr: os.ErrNotExist,
			want:      false,
		},
		{
			name:      "ambiguous, log and default to false",
			osStatErr: errors.Errorf("12345"),
			want:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stat := func(filename string) (os.FileInfo, error) {
				return nil, tc.osStatErr
			}
			got := pathExists("x", stat)
			if tc.want != got {
				t.Errorf("want %v got %v", tc.want, got)
			}
		})
	}
}

func TestRemovePathIfExists(t *testing.T) {
	cases := []struct {
		name         string
		pathExists   bool
		removeCalled bool
		removeErr    error
		want         error
	}{
		{
			name:         "does not exist",
			pathExists:   false,
			removeCalled: false,
			removeErr:    nil,
			want:         nil,
		},
		{
			name:         "exists and successfully removed",
			pathExists:   true,
			removeCalled: true,
			removeErr:    nil,
			want:         nil,
		},
		{
			name:         "exists and failed to removed",
			pathExists:   true,
			removeCalled: true,
			removeErr:    errors.Errorf("Error!"),
			want:         errors.Errorf("Error!"),
		},
	}

	for _, tc := range cases {
		gotRemoveCalled := false
		remove := func(name string) error {
			gotRemoveCalled = true
			return tc.removeErr
		}
		pathExists := func(path string) bool {
			return tc.pathExists
		}
		got := removePathIfExists("foobar", remove, pathExists)
		testutils.CompareTestErrors(tc.want, got, t)
		if tc.removeCalled != gotRemoveCalled {
			t.Errorf("Remove called: want '%v' got '%v'", tc.removeCalled, gotRemoveCalled)
		}
	}
}
