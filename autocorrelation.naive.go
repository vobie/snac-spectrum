package main

import (
	"github.com/go-audio/audio"
	"github.com/vobie/snac-spectrum/utils"
)

// Naive algortihm: pointwise multiplication. For sanity checking more efficient power spectrum version.
func NaiveAutocorrelation(buf *audio.IntBuffer, maxShift int) []float64 {
	samples := utils.BufferToFloat64(buf)

	result := make([]float64, maxShift)
	slen := len(samples)

	for shift := 0; shift < maxShift; shift++ {
		var sum float64
		for i := 0; i < slen-maxShift; i++ {
			sum += samples[i] * samples[i+shift]
		}
		result[shift] = sum //or slen
	}

	return result
}
