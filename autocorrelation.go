package main

// Naive algortihm: pointwise multiplication. For sanity checking more efficient power spectrum version.
func NaiveAutocorrelation(samples []float64, maxShift int) []float64 {
	result := make([]float64, maxShift)
	slen := len(samples)

	for shift := 0; shift < maxShift; shift++ {
		var sum float64
		for i := 0; i < slen-maxShift; i++ {
			sum += samples[i] * samples[i+shift]
		}
		result[shift] = sum
	}

	return result
}
