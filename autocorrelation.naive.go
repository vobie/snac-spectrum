package main

import (
	"github.com/go-audio/audio"
	"github.com/vobie/snac-spectrum/utils"
)

// Naive algortihm: pointwise multiplication. For sanity checking more efficient power spectrum version.
func NaiveAutocorrelation(buf *audio.IntBuffer) []float64 {
	n := buf.NumFrames()
	samples := utils.BufferToFloat64(buf)

	result := make([]float64, n)
	slen := len(samples)

	for shift := 0; shift < n; shift++ {
		var sum float64
		for i := 0; i < slen-shift; i++ {
			sum += samples[i] * samples[i+shift]
		}
		result[shift] = sum / float64(slen-shift) // Unbias: Divide by number of samples multiplied so that it's autocorrelation per sample investigated
		// Use NaiveAutocorrelationNorm for a proper normalization - it computes the power on the range investigated on each shift
	}
	return result
}

// Naive algortihm: pointwise multiplication. For sanity checking more efficient power spectrum version.
func NaiveAutocorrelationNorm(buf *audio.IntBuffer) []float64 {
	n := buf.NumFrames()
	samples := utils.BufferToFloat64(buf)

	result := make([]float64, n)
	slen := len(samples)
	cumulativePowers := utils.CumulativeTotalPower(buf)

	for shift := 0; shift < n; shift++ {
		power := cumulativePowers[n-shift-1]

		var sum float64
		for i := 0; i < slen-shift; i++ {
			sum += samples[i] * samples[i+shift]
		}

		if power != 0 {
			result[shift] = (sum / power)
		} else {
			result[shift] = 0
		}
	}

	return result
}
