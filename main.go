package main

import "fmt"

func main() {
	fmt.Printf("Zeroes autocorrelation:")
	fmt.Println(NaiveAutocorrelation([]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 3))
	fmt.Printf("Ones autocorrelation:")
	fmt.Println(NaiveAutocorrelation([]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 3))
	fmt.Printf("Square(2samples) autocorrelation:")
	fmt.Println(NaiveAutocorrelation([]float64{1, 1, -1, -1, 1, 1, -1, -1, 1, 1}, 3))
}
