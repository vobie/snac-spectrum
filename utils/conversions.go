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

func CumulativePower(buf *audio.IntBuffer) []float64 {
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

func TotalPower(buf *audio.IntBuffer) float64 {
	var totalPower float64
	for _, sample := range buf.Data {
		totalPower += float64(sample * sample)
	}
	return totalPower
}
