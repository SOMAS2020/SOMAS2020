package config

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetSelectivelyVisibleFloat64(t *testing.T) {
	origVal := float64(69)
	cases := []struct {
		name  string
		valid bool
		want  SelectivelyVisibleFloat64
	}{
		{
			name:  "valid",
			valid: true,
			want: SelectivelyVisibleFloat64{
				Value: origVal,
				Valid: true,
			},
		},
		{
			name:  "not valid",
			valid: false,
			want: SelectivelyVisibleFloat64{
				Value: 0,
				Valid: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := getSelectivelyVisibleFloat64(origVal, tc.valid)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestGetSelectivelyVisibleResources(t *testing.T) {
	origVal := shared.Resources(69)
	cases := []struct {
		name  string
		valid bool
		want  SelectivelyVisibleResources
	}{
		{
			name:  "valid",
			valid: true,
			want: SelectivelyVisibleResources{
				Value: origVal,
				Valid: true,
			},
		},
		{
			name:  "not valid",
			valid: false,
			want: SelectivelyVisibleResources{
				Value: shared.Resources(0),
				Valid: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := getSelectivelyVisibleResources(origVal, tc.valid)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestGetSelectivelyVisibleInteger(t *testing.T) {
	origVal := 5
	cases := []struct {
		name  string
		valid bool
		want  SelectivelyVisibleInteger
	}{
		{
			name:  "valid",
			valid: true,
			want: SelectivelyVisibleInteger{
				Value: origVal,
				Valid: true,
			},
		},
		{
			name:  "not valid",
			valid: false,
			want: SelectivelyVisibleInteger{
				Value: 0,
				Valid: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := getSelectivelyVisibleInteger(origVal, tc.valid)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}
