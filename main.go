package main

import (
	"fmt"
	"math"
	"os"

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

	data := make(plotter.XYs, 2000)
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
	fmt.Printf("[naive] Square(2samples) autocorrelation:")
	fmt.Println(utils.NormalizeArray(NaiveAutocorrelation(utils.MakeBuffer(squareSignal), 3)))

	// fmt.Printf("[optimized] Zeroes autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})))
	// fmt.Printf("[optimized] Ones autocorrelation:")
	// fmt.Println(OptimizedAutocorrelation(utils.MakeBuffer([]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1})))
	fmt.Printf("[optimized] Square(2samples) autocorrelation:")
	fmt.Println(utils.NormalizeArray(OptimizedAutocorrelation(utils.MakeBuffer(squareSignal))))

	// WAV stuff
	file, err := os.Open("440hz.wav")
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

	// Plots
	naive := utils.NormalizeArray(NaiveAutocorrelation(fullBuffer, 200))
	opti := utils.NormalizeArray(OptimizedAutocorrelation(fullBuffer)[:200])
	fmt.Printf("[naive] 440hz autocorrelation:")
	fmt.Println(naive)
	fmt.Printf("[optimized] 440hz autocorrelation:")
	fmt.Println(opti)

	naivePlot := plotAutocorrelation(naive, "Naive autocorrelation")
	naivePlot.Save(4*vg.Inch, 4*vg.Inch, "naive.png")

	optiPlot := plotAutocorrelation(opti, "Optimized autocorrelation")
	optiPlot.Save(4*vg.Inch, 4*vg.Inch, "optimized.png")
}
