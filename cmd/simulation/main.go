package main

import (
	"os"

	"github.com/monochromegane/kaburaya"
)

type generator struct {
}

func (g generator) Generate(i int) [][]float64 {
	if i == 1 {
		return [][]float64{
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
			[]float64{10.0, 20.0, 30.0},
		}
	}
	if i > 25 && i < 88 {
		return [][]float64{
			[]float64{1.0, 2.0, 3.0},
			[]float64{5.0, 2.0, 3.0},
		}
	}

	if i > 110 && i < 140 {
		return [][]float64{
			[]float64{10.0, 2.0, 3.0},
			[]float64{5.0, 20.0, 30.0},
			[]float64{5.0, 5.0, 3.0, 2.0, 1.0},
		}
	}
	return [][]float64{[]float64{}}
}

func main() {
	simulator := kaburaya.Simulator{
		InitialNumWorker: 1,
		Resource:         100.0,
		Generator:        generator{},
		// Controller:       &kaburaya.FixController{},
		// Controller: &kaburaya.SimpleController{100.0},
		// Controller: &kaburaya.PController{Target: 90.0, K: 0.1},
		// Controller: &kaburaya.DynamicController{Span: 3, Controller: &kaburaya.PController{Target: 90.0, K: 0.1}},
		// Controller: &kaburaya.DynamicController{Span: 3, Controller: &kaburaya.PIController{Target: 0.0, Kp: 0.1, Ki: 0.01}},
		// OK edition
		Controller: &kaburaya.DynamicController{Span: 3, Controller: &kaburaya.PIDShortController{Target: 0.0, Kp: 0.1, Ki: 0.05, Kd: 0.1, Span: 100}},
		// Controller: &kaburaya.DynamicController{Span: 3, Controller: &kaburaya.PIDShortController{Target: 0.0, Kp: 0.1, Ki: 0.5, Kd: 0.5, Span: 10}},
		// Controller: &kaburaya.StabilityController{Span: 3, SD: 1.0, Controller: &kaburaya.PController{Target: 90.0, K: 0.1}},
		// Controller: &kaburaya.RateController{Controller: &kaburaya.PController{Target: 10.0, K: 0.05}},
		// Controller: &kaburaya.DynamicRateController{Span: 5, Controller: &kaburaya.RateController{Controller: &kaburaya.PController{Target: 10.0, K: 0.1}}},
		// Controller: &kaburaya.PIController{Target: 90.0, Kp: 0.1, Ki: 0.1},
	}
	results := simulator.Run(160)

	file, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, result := range results {
		file.WriteString(result.String())
	}

}
