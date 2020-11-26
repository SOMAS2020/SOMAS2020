package forage

import (
	"reflect"
	"testing"
)

func TestForageSplitEven(t *testing.T) {
	cases := []struct {
		name                  string
		teamForageInvestments map[int]int
		totalPayoff           int
		want                  map[int]int
	}{
		{
			name: "basic even split",
			teamForageInvestments: map[int]int{
				1: 123,
				2: 123,
			},
			totalPayoff: 420,
			want: map[int]int{
				1: 210,
				2: 210,
			},
		},
		{
			name: "not completely divisible",
			teamForageInvestments: map[int]int{
				3: 123,
				4: 123,
			},
			totalPayoff: 419,
			want: map[int]int{
				3: 209,
				4: 209,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := forageSplitEven(tc.teamForageInvestments, tc.totalPayoff)

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestForageSplitProportionate(t *testing.T) {
	cases := []struct {
		name                  string
		teamForageInvestments map[int]int
		totalPayoff           int
		want                  map[int]int
	}{
		{
			name: "basic even split",
			teamForageInvestments: map[int]int{
				1: 123,
				2: 123,
			},
			totalPayoff: 420,
			want: map[int]int{
				1: 210,
				2: 210,
			},
		},
		{
			name: "1:2 ratio",
			teamForageInvestments: map[int]int{
				3: 1,
				4: 2,
			},
			totalPayoff: 30,
			want: map[int]int{
				3: 10,
				4: 20,
			},
		},
		{
			name: "not completely divisible, some resources discarded",
			teamForageInvestments: map[int]int{
				3: 1,
				4: 2,
			},
			totalPayoff: 31,
			want: map[int]int{
				3: 10,
				4: 20,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := forageSplitProprotionate(tc.teamForageInvestments, tc.totalPayoff)

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}
