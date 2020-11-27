package common

import (
	"reflect"
	"testing"
)

func TestForageSplitEven(t *testing.T) {
	cases := []struct {
		name                  string
		teamForageInvestments map[ClientID]uint
		totalPayoff           uint
		want                  map[ClientID]uint
	}{
		{
			name: "basic even split",
			teamForageInvestments: map[ClientID]uint{
				1: 123,
				2: 123,
			},
			totalPayoff: 420,
			want: map[ClientID]uint{
				1: 210,
				2: 210,
			},
		},
		{
			name: "not completely divisible",
			teamForageInvestments: map[ClientID]uint{
				3: 123,
				4: 123,
			},
			totalPayoff: 419,
			want: map[ClientID]uint{
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
		teamForageInvestments map[ClientID]uint
		totalPayoff           uint
		want                  map[ClientID]uint
	}{
		{
			name: "basic even split",
			teamForageInvestments: map[ClientID]uint{
				1: 123,
				2: 123,
			},
			totalPayoff: 420,
			want: map[ClientID]uint{
				1: 210,
				2: 210,
			},
		},
		{
			name: "1:2 ratio",
			teamForageInvestments: map[ClientID]uint{
				3: 1,
				4: 2,
			},
			totalPayoff: 30,
			want: map[ClientID]uint{
				3: 10,
				4: 20,
			},
		},
		{
			name: "not completely divisible, some resources discarded",
			teamForageInvestments: map[ClientID]uint{
				3: 1,
				4: 2,
			},
			totalPayoff: 31,
			want: map[ClientID]uint{
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
