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

https://en.wikipedia.org/wiki/Fast_Fourier_transform#FFT_algorithms_specialized_for_real_or_symmetric_data
Evaluate if such an algorithm could be used here, or maybe already is. FFTReal() kind of suggests this, but double check
Benchmark each of "real->half size complex", "Sorensen, 1987 or similar" and the "removing redundant parts" algortihms. Ask an expert (katjaas?)
*/
func OptimizedAutocorrelation(buf *audio.IntBuffer) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := buf.NumFrames()
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

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

func OptimizedAutocorrelationNorm(buf *audio.IntBuffer) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := buf.NumFrames()
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

	frequencySpectrum := fft.FFTReal(utils.BufferToFloat64(buf))
	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
	autocorr := fft.IFFT(powerSpectrum)
	averagePower := utils.AveragePowerFromSpectrum(powerSpectrum)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val) / averagePower
	}

	return autocorrReal
}

func OptimizedAutocorrelationNorm2(buf *audio.IntBuffer) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := buf.NumFrames()
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

	frequencySpectrum := fft.FFTReal(utils.BufferToFloat64(buf))
	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
	autocorr := fft.IFFT(powerSpectrum)
	averagePowers := utils.CumulativeAveragePowerPerSample(buf) //1st nan

	fmt.Print(averagePowers)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val) / averagePowers[i]
	}

	return autocorrReal
}
