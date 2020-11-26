package forage

import "testing"

func TestForagePayoffSquare(t *testing.T) {
	cases := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "basic 0",
			input: 0,
			want:  0,
		},
		{
			name:  "basic 1",
			input: 1,
			want:  1,
		},
		{
			name:  "42",
			input: 42,
			want:  42 * 42,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := foragePayoffSquare(tc.input)
			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestForagePayoffCube(t *testing.T) {
	cases := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "basic 0",
			input: 0,
			want:  0,
		},
		{
			name:  "basic 1",
			input: 1,
			want:  1,
		},
		{
			name:  "42",
			input: 42,
			want:  42 * 42 * 42,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := foragePayoffCube(tc.input)
			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}
