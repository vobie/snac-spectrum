package main

import (
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
	fftResult := fft.FFTReal(utils.BufferToFloat64(buf))

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
