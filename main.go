package main

import (
	"errors"
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

func maxIndexRange(arr []float64, start int, end int) (int, error) {
	if start < 0 || end >= len(arr) || start > end {
		return 0, errors.New("invalid range")
	}

	max := arr[start]
	maxI := 0
	for i := start + 1; i <= end; i++ {
		if arr[i] > max {
			max = arr[i]
			maxI = i
		}
	}

	return maxI, nil
}

func main() {
	squareSignal := []int{math.MaxInt, math.MaxInt, math.MinInt, math.MinInt, math.MaxInt, math.MaxInt, math.MinInt, math.MinInt, math.MaxInt, math.MaxInt}

	// fmt.Printf("[naive] Zeroes autocorrelation:")
	// fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 3))
	// fmt.Printf("[naive] Ones autocorrelation:")
	// fmt.Println(NaiveAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}), 3))
	fmt.Printf("[naive] Square(2samples) autocorrelation: \n")
	fmt.Println(utils.NormalizeArray(NaiveAutocorrelation(utils.IntToFloat64(squareSignal))))

	// fmt.Printf("[optimized] Zeroes autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})))
	// fmt.Printf("[optimized] Ones autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1})))
	// fmt.Printf("[optimized] Square(2samples) autocorrelation: \n")
	// fmt.Println(utils.NormalizeArray(OptimizedAutocorrelation(utils.IntToFloat64(squareSignal))))

	// WAV stuff
	file, err := os.Open("C523.wav")
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

	slicedBuffer := utils.SliceBuffer(fullBuffer, 1024) //About 40 full cycles of A440
	slicedBuffer = utils.SnacPadBuffer(slicedBuffer)

	fmt.Println()
	fmt.Println()
	fmt.Printf("***** A440 TEST ***** \n")
	fmt.Printf("Buffer size: %d \n", slicedBuffer.NumFrames())

	start := time.Now()
	naive := NaiveAutocorrelationNorm(utils.BufferToFloat64(slicedBuffer)) // Skip very small shifts, error too big. TODO investigate what happens here.
	fmt.Printf("Naive: %v\n", time.Since(start))

	start = time.Now()
	opti := OptimizedAutocorrelation(utils.BufferToFloat64(slicedBuffer)) // ATTN: FFT WAY slower if frame size not divisible by 2
	fmt.Printf("Optimized: %v\n", time.Since(start))

	naivePlot := plotAutocorrelation(naive[:500], "Naive autocorrelation")
	naivePlot.Save(4*vg.Inch, 4*vg.Inch, "naive.png")

	optiPlot := plotAutocorrelation(opti, "Optimized autocorrelation")
	optiPlot.Save(4*vg.Inch, 4*vg.Inch, "optimized.png")

	triangle := SnacNormalizeTriangle(opti, utils.BufferToFloat64(slicedBuffer))

	trianglePlot := plotAutocorrelation(triangle, "SNAC triangle for autocorrelation")
	trianglePlot.Save(4*vg.Inch, 4*vg.Inch, "triangle.png")

	normdOpti := make([]float64, len(opti))
	for n := 0; n < len(triangle)/2; n++ {
		normdOpti[n] = opti[n] / triangle[n]
	}

	normdOptiPlot := plotAutocorrelation(normdOpti[:200], "SNAC normalized opti")
	normdOptiPlot.Save(4*vg.Inch, 4*vg.Inch, "normdTri.png")

	maxi, err := maxIndexRange(normdOpti, 75, 150)
	hz := float64(fullBuffer.Format.SampleRate) / float64(maxi)
	fmt.Printf("MAX: %d/%d = %dhz", maxi, fullBuffer.Format.SampleRate, hz)

}
