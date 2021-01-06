package config

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// SelectivelyVisibleFloat64 represents a wrapped float64 whose value is valid only if the Valid flag is set to true
type SelectivelyVisibleFloat64 struct {
	Value float64
	Valid bool
}

func getSelectivelyVisibleFloat64(value float64, valid bool) SelectivelyVisibleFloat64 {
	var res float64
	if valid {
		res = value
	}
	return SelectivelyVisibleFloat64{
		Value: res,
		Valid: valid,
	}
}

// SelectivelyVisibleResources represents a wrapped Resources whose value is valid only if the Valid flag is set to true
type SelectivelyVisibleResources struct {
	Value shared.Resources
	Valid bool
}

type SelectivelyVisibleUint struct {
	Value uint
	Valid bool
}

type SelectivelyVisibleBool struct {
	Value bool
	Valid bool
}

func getSelectivelyVisibleResources(value shared.Resources, valid bool) SelectivelyVisibleResources {
	var res shared.Resources
	if valid {
		res = value
	}
	return SelectivelyVisibleResources{
		Value: res,
		Valid: valid,
	}
}
func getSelectivelyVisibleUint(value uint, valid bool) SelectivelyVisibleUint {
	var val uint
	if valid {
		val = value
	}
	return SelectivelyVisibleUint{
		Value: val,
		Valid: valid,
	}
}
func getSelectivelyVisibleBool(value bool, valid bool) SelectivelyVisibleBool {
	var val bool
	if valid {
		val = value
	}
	return SelectivelyVisibleBool{
		Value: val,
		Valid: valid,
	}
}
