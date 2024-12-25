package indicator

import (
	"math"
)

type SqueezeIndicator struct {
	Length       int
	Mult         float64
	LengthKC     int
	MultKC       float64
	UseTrueRange bool
}

func NewSqueezeIndicator(length, lengthKC int, mult, multKC float64, useTrueRange bool) *SqueezeIndicator {
	return &SqueezeIndicator{
		Length:       length,
		Mult:         mult,
		LengthKC:     lengthKC,
		MultKC:       multKC,
		UseTrueRange: useTrueRange,
	}
}

func (si *SqueezeIndicator) Calculate(close, high, low []float64) ([]float64, []bool, []bool) {
	n := len(close)
	val := make([]float64, n)
	sqzOn := make([]bool, n)
	sqzOff := make([]bool, n)

	basis := sma(close, si.Length)
	dev := multiply(stdev(close, si.Length), si.Mult)
	upperBB := add(basis, dev)
	lowerBB := subtract(basis, dev)

	ma := sma(close, si.LengthKC)
	var range_ []float64
	if si.UseTrueRange {
		range_ = trueRange(high, low, close)
	} else {
		range_ = subtract(high, low)
	}
	rangema := sma(range_, si.LengthKC)
	upperKC := add(ma, multiply(rangema, si.MultKC))
	lowerKC := subtract(ma, multiply(rangema, si.MultKC))

	for i := si.LengthKC - 1; i < n; i++ {
		val[i] = linearRegression(getWindow(close, i, si.LengthKC), si.LengthKC)

		sqzOn[i] = lowerBB[i] > lowerKC[i] && upperBB[i] < upperKC[i]
		sqzOff[i] = lowerBB[i] < lowerKC[i] && upperBB[i] > upperKC[i]
	}

	return val, sqzOn, sqzOff
}


func sma(data []float64, period int) []float64 {
	result := make([]float64, len(data))
	sum := 0.0

	for i := 0; i < len(data); i++ {
		sum += data[i]
		if i >= period {
			sum -= data[i-period]
		}
		if i >= period-1 {
			result[i] = sum / float64(period)
		}
	}
	return result
}

func stdev(data []float64, period int) []float64 {
	result := make([]float64, len(data))

	for i := period - 1; i < len(data); i++ {
		window := data[i-period+1 : i+1]
		mean := average(window)
		sumSquares := 0.0
		for _, v := range window {
			sumSquares += (v - mean) * (v - mean)
		}
		result[i] = math.Sqrt(sumSquares / float64(period))
	}
	return result
}

func trueRange(high, low, close []float64) []float64 {
	result := make([]float64, len(high))
	result[0] = high[0] - low[0]

	for i := 1; i < len(high); i++ {
		hl := high[i] - low[i]
		hc := math.Abs(high[i] - close[i-1])
		lc := math.Abs(low[i] - close[i-1])
		result[i] = math.Max(hl, math.Max(hc, lc))
	}
	return result
}

func linearRegression(data []float64, period int) float64 {
	n := len(data)
	if n < period {
		return 0
	}

	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i := 0; i < period; i++ {
		x := float64(i)
		y := data[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (float64(period)*sumXY - sumX*sumY) / (float64(period)*sumX2 - sumX*sumX)
	return slope
}

func add(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

func subtract(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

func multiply(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] * b
	}
	return result
}

func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func highest(data []float64) float64 {
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func lowest(data []float64) float64 {
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func getWindow(data []float64, currentIndex, period int) []float64 {
	start := currentIndex - period + 1
	if start < 0 {
		start = 0
	}
	return data[start : currentIndex+1]
}
