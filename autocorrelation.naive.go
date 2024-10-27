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
		result[shift] = sum
	}

	return result
}
