package main

import (
	"fmt"
	"math"
	"os"
	"time"

	//"github.com/go-audio/wav"
	"github.com/go-audio/wav"
	"github.com/vobie/snac-spectrum/utils"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	// "github.com/vobie/snac-spectrum/utils" // Adjust this import to match your utils package
)

func plotAutocorrelation(correlation []float64, title string) *plot.Plot {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Shift"
	p.Y.Label.Text = "Autocorrelation Value"
	p.Y.Max = 1
	p.Y.Min = -1

	data := make(plotter.XYs, len(correlation))
	for i := range correlation {
		data[i] = plotter.XY{X: float64(i), Y: correlation[i]}
	}

	line, err := plotter.NewLine(data)

	if err != nil {
		panic(err)
	}
	line.LineStyle.Width = 2
	p.Add(line)

	for i := 0; i < len(correlation); i++ {
		line.XYs[i].X = float64(i)
	}

	return p
}

func main() {
	squareSignal := []int{math.MaxInt, math.MaxInt, math.MinInt, math.MinInt, math.MaxInt, math.MaxInt, math.MinInt, math.MinInt, math.MaxInt, math.MaxInt}

	// fmt.Printf("[naive] Zeroes autocorrelation:")
	// fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 3))
	// fmt.Printf("[naive] Ones autocorrelation:")
	// fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}), 3))
	fmt.Printf("[naive] Square(2samples) autocorrelation: \n")
	fmt.Println(utils.NormalizeArray(NaiveAutocorrelation(utils.MakeBuffer(squareSignal))))

	// fmt.Printf("[optimized] Zeroes autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})))
	// fmt.Printf("[optimized] Ones autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1})))
	fmt.Printf("[optimized] Square(2samples) autocorrelation: \n")
	fmt.Println(utils.NormalizeArray(OptimizedAutocorrelation(utils.MakeBuffer(squareSignal))))

	// WAV stuff
	file, err := os.Open("uneven_8192.wav")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	fullBuffer, err := decoder.FullPCMBuffer()
	if err != nil {
		fmt.Println("Error reading full PCM buffer:", err)
		return
	}

	slicedBuffer := utils.SliceBuffer(fullBuffer, 8192) //About 40 full cycles of A440

	fmt.Println()
	fmt.Println()
	fmt.Printf("***** A440 TEST ***** \n")
	fmt.Printf("Buffer size: %d \n", slicedBuffer.NumFrames())

	start := time.Now()
	naive := NaiveAutocorrelationNorm(slicedBuffer)[:8100] // Skip very small shifts, error too big. TODO investigate what happens here.
	fmt.Printf("Naive: %v\n", time.Since(start))

	start = time.Now()
	opti := OptimizedAutocorrelationNorm(slicedBuffer) // ATTN: FFT WAY slower if frame size not divisible by 2
	fmt.Printf("Optimized: %v\n", time.Since(start))

	naivePlot := plotAutocorrelation(naive, "Naive autocorrelation")
	naivePlot.Save(4*vg.Inch, 4*vg.Inch, "naive.png")

	optiPlot := plotAutocorrelation(opti, "Optimized autocorrelation")
	optiPlot.Save(4*vg.Inch, 4*vg.Inch, "optimized.png")
}
