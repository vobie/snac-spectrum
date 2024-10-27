package utils

import (
	"github.com/go-audio/audio"
)

func BufferToFloat64(buf *audio.IntBuffer) []float64 {
	samples := buf.Data
	fSamples := make([]float64, len(samples))
	for i, sample := range buf.Data {
		fSamples[i] = float64(sample)
	}
	return fSamples
}

func MakeBuffer(samples []int) *audio.IntBuffer {
	return &audio.IntBuffer{
		Format: &audio.Format{SampleRate: 44000, NumChannels: 1},
		Data:   samples,
	}
}

// NormalizeArray normalizes a float64 slice to the range [0, 1].
func NormalizeArray(arr []float64) []float64 {
	if len(arr) == 0 {
		return arr // Return empty array as is
	}

	// Find min and max values
	min := arr[0]
	max := arr[0]
	for _, value := range arr {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	// Handle case where max == min to avoid division by zero
	if max == min {
		return make([]float64, len(arr)) // Return an array of zeros
	}

	// Normalize the array
	normalized := make([]float64, len(arr))
	for i, value := range arr {
		normalized[i] = (value - min) / (max - min)
	}

	return normalized
}

func SliceBuffer(existingBuffer *audio.IntBuffer, N int) *audio.IntBuffer {
	if N > len(existingBuffer.Data) {
		N = len(existingBuffer.Data)
	}

	newBuffer := &audio.IntBuffer{
		Data:   existingBuffer.Data[:N],
		Format: existingBuffer.Format,
	}

	return newBuffer
}

func PowerSpectrum(frequencySpectrum []complex128) []complex128 {
	powerSpectrum := make([]complex128, len(frequencySpectrum))
	for i, val := range frequencySpectrum {
		powerSpectrum[i] = complex(real(val)*real(val)+imag(val)*imag(val), 0)
	}
	return powerSpectrum
}

/*
Produces a series where cumulativeTotalPower[n] is the cumulative power in the input signal from input[0:n]
*/
func CumulativeTotalPower(buf *audio.IntBuffer) []float64 {
	cumulativePower := make([]float64, buf.NumFrames())
	data := buf.Data
	prev := float64(0)
	for i, sample := range data {
		power := prev + float64(sample)*float64(sample)
		cumulativePower[i] = power
		prev = power
	}
	return cumulativePower
}

func CumulativeAveragePowerPerSample(buf *audio.IntBuffer) []float64 {
	cumulativePowerPS := make([]float64, buf.NumFrames())
	data := buf.Data
	prev := float64(0)
	for i, sample := range data {
		power := prev + float64(sample)*float64(sample)
		cumulativePowerPS[i] = power
		prev = power
	}
	for i, _ := range cumulativePowerPS {
		if cumulativePowerPS[i] == 0 {
			cumulativePowerPS[i] = 1
		} else {
			cumulativePowerPS[i] /= float64(i)
		}
	}
	return cumulativePowerPS
}

func TotalPower(buf *audio.IntBuffer) float64 {
	var totalPower float64
	for _, sample := range buf.Data {
		totalPower += float64(sample * sample)
	}
	return totalPower
}

/*
Unsure if appropriately named. Normalizes autocorrelation computed via power spectrum properly to (-1,1)
*/
func AveragePowerFromSpectrum(powerSpectrum []complex128) float64 {
	return TotalPowerFromSpectrum(powerSpectrum) / float64(len(powerSpectrum))
}

/*
Unsure if appropriately named
*/
func TotalPowerFromSpectrum(powerSpectrum []complex128) float64 {
	var totalPower float64
	totalPower += real(powerSpectrum[0])

	for i := 1; i < len(powerSpectrum)/2; i++ {
		totalPower += 2 * real(powerSpectrum[i])
	}

	if len(powerSpectrum)%2 == 0 {
		totalPower += real(powerSpectrum[len(powerSpectrum)/2])
	}

	return totalPower
}
