package main

import (
	"fmt"

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
// func OptimizedAutocorrelation(fSamples []float64) []float64 {
// 	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
// 	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
// 	n := len(fSamples)
// 	if n > 0 && (n&(n-1)) != 0 {
// 		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
// 	}

// 	fftResult := fft.FFTReal(fSamples)

// 	// FIXME: Optimize this step
// 	powerSpectrum := make([]complex128, len(fftResult))
// 	for i, val := range fftResult {
// 		powerSpectrum[i] = complex(real(val)*real(val)+imag(val)*imag(val), 0)
// 	}

// 	autocorr := fft.IFFT(powerSpectrum)

// 	autocorrReal := make([]float64, len(autocorr))
// 	for i, val := range autocorr {
// 		autocorrReal[i] = real(val)
// 	}

// 	return autocorrReal
// }

func OptimizedAutocorrelationNorm(fSamples []float64) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := len(fSamples)
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

	frequencySpectrum := fft.FFTReal(fSamples)
	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
	autocorr := fft.IFFT(powerSpectrum)
	averagePower := utils.AveragePowerFromSpectrum(powerSpectrum)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val) / averagePower
	}

	return autocorrReal
}

func OptimizedBiasedAutocorrelation(fSamples []float64) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := len(fSamples)
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

	frequencySpectrum := fft.FFTReal(fSamples)
	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
	autocorr := fft.IFFT(powerSpectrum)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val)
	}

	return autocorrReal
}

func OptimizedAutocorrelation(fSamples []float64) []float64 {
	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
	n := len(fSamples)
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
	}

	frequencySpectrum := fft.FFTReal(fSamples)
	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
	autocorr := fft.IFFT(powerSpectrum)

	autocorrReal := make([]float64, len(autocorr))
	for i, val := range autocorr {
		autocorrReal[i] = real(val)
	}

	normalized := SnacNormalize(autocorrReal, fSamples)

	return normalized
}

func SnacNormalize(biased []float64, inputBuffer []float64) []float64 {

	// will get more logic for noise filter later
	rzero := biased[0]
	frameSize := len(inputBuffer)
	seek := frameSize / 2

	// Modified SNAC unbiasing "triangle"
	normalized := make([]float64, frameSize)
	normalized[0] = 1

	// Computed sequentially - see katjaas site
	normIntegral := 2 * rzero
	for n := 1; n < seek; n++ {
		s1 := inputBuffer[n-1]
		s2 := inputBuffer[frameSize-n]
		normIntegral -= (s1 * s1) + (s2 * s2)

		normalized[n] = biased[n] / (normIntegral * 0.5)
		println(normalized[n])
	}

	return normalized
}

// Failed experiment - The weird effects are from using the power spectrum to calculate autocorrelation, NOT from the average power calc
// func OptimizedAutocorrelationNorm2(buf *audio.IntBuffer) []float64 {
// 	// ATTN: This may be an issue for SNAC when choosing window size, resampling may be needed
// 	// Could implement automatic resampling for convenience here, but always warn as time is of the essence
// 	n := buf.NumFrames()
// 	if n > 0 && (n&(n-1)) != 0 {
// 		fmt.Printf("FFT WARNING: Frame size (%d) not power of 2\n", n)
// 	}

// 	frequencySpectrum := fft.FFTReal(utils.BufferToFloat64(buf))
// 	powerSpectrum := utils.PowerSpectrum(frequencySpectrum)
// 	autocorr := fft.IFFT(powerSpectrum)
// 	averagePowers := utils.CumulativeAveragePowerPerSample(buf) //1st nan

// 	fmt.Print(averagePowers)

// 	autocorrReal := make([]float64, len(autocorr))
// 	for i, val := range autocorr {
// 		autocorrReal[i] = real(val) / averagePowers[i]
// 	}

// 	return autocorrReal
// }
