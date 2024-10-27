package main

import (
	"fmt"

	"github.com/vobie/snac-spectrum/utils"
)

func main() {
	fmt.Printf("Zeroes autocorrelation:")
	fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 3))
	fmt.Printf("Ones autocorrelation:")
	fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}), 3))
	fmt.Printf("Square(2samples) autocorrelation:")
	fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{1, 1, -1, -1, 1, 1, -1, -1, 1, 1}), 3))

	// file, err := os.Open("./440hz.wav")
	// if err != nil {
	// 	log.Fatalf("Could not read audio file %s", err)
	// }
	// defer file.Close()

	// decoder := wav.NewDecoder(file)
	// if !decoder.IsValidFile() {
	// 	log.Fatalf("invalid WAV file")
	// }

	// buffer, err := decoder.FullPCMBuffer()
	// if err != nil {
	// 	log.Fatalf("WAV decoding failed %s", err)
	// }

	// samples := buffer.Data
	// fSamples := make([]float64, len(samples))
	// for i, sample := range buffer.Data {
	// 	fSamples[i] = float64(sample)
	// }

	// autocorrelation := NaiveAutocorrelation(fSamples, 200)
	// for i, ac := range autocorrelation {
	// 	println(i, ac)
	// }
}
