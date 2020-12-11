package foraging

func sum(a []float64) float64 {
	result := 0.0
	for _, v := range a {
		result += v
	}
	return result
}
