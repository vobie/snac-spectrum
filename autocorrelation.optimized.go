package main

import (
	"fmt"

	"github.com/go-audio/audio"
	"github.com/mjibson/go-dsp/fft"
	"github.com/vobie/snac-spectrum/utils"
)

/*
Computes autocorrelation using Wiener-Khinchin theorem

1. FFT
2. Power spectrum (square coeff)
3. Inverse FFT
4. Discard complex component
*/
func OptimizedAutocorrelation(buf *audio.IntBuffer) []float64 {

	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := buf.NumFrames()
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size not power of 2\n")
	}

	// ATTN: FFT is WAY slower if frame size not divisible by 2
	fftResult := fft.FFTReal(utils.BufferToFloat64(buf))

	// FIXME: Optimize this step
	powerSpectrum := make([]complex128, len(fftResult))
	for i, val := range fftResult {
		powerSpectrum[i] = complex(real(val)*real(val)+imag(val)*imag(val), 0)
	}

	autocorr := fft.IFFT(powerSpectrum)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val)
	}

	return autocorrReal
}
