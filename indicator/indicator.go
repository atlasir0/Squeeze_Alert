package indicator

import (
	"math"
)

type SqueezeIndicator struct {
	BBLength      int
	BBMult        float64
	KCLength      int
	KCMult        float64
	UseTrueRange  bool
	MinVolatility float64
}

func NewSqueezeIndicator(bbLength, kcLength int, bbMult, kcMult float64, useTrueRange bool, minVolatility float64) *SqueezeIndicator {
	return &SqueezeIndicator{
		BBLength:      bbLength,
		BBMult:        bbMult,
		KCLength:      kcLength,
		KCMult:        kcMult,
		UseTrueRange:  useTrueRange,
		MinVolatility: minVolatility,
	}
}

func (si *SqueezeIndicator) Calculate(close, high, low []float64) ([]float64, []bool) {
	n := len(close)
	if n == 0 {
		return []float64{}, []bool{}
	}

	basis := sma(close, si.BBLength)
	dev := multiply(stdev(close, si.BBLength), si.BBMult)
	upperBB := add(basis, dev)
	lowerBB := subtract(basis, dev)

	ma := sma(close, si.KCLength)
	var range_ []float64
	if si.UseTrueRange {
		range_ = trueRange(high, low, close)
	} else {
		range_ = subtract(high, low)
	}
	rangema := sma(range_, si.KCLength)
	upperKC := add(ma, multiply(rangema, si.KCMult))
	lowerKC := subtract(ma, multiply(rangema, si.KCMult))

	atrValues := atr(high, low, close, 14)
	isVolatile := make([]bool, n)
	for i := range atrValues {
		isVolatile[i] = atrValues[i] > si.MinVolatility
	}

	val := make([]float64, n)
	for i := si.KCLength - 1; i < n; i++ {
		highestHigh := max(high[i-(si.KCLength-1) : i+1])
		lowestLow := min(low[i-(si.KCLength-1) : i+1])

		avgHL := (highestHigh + lowestLow) / 2

		smaClose := average(close[i-(si.KCLength-1) : i+1])

		regressionInput := make([]float64, si.KCLength)
		for j := 0; j < si.KCLength; j++ {
			regressionInput[j] = close[i-(si.KCLength-1)+j] - (avgHL+smaClose)/2
		}

		val[i] = linearRegression(regressionInput, si.KCLength)
	}

	sqzOn := make([]bool, n)
	for i := range close {
		if i < si.KCLength-1 {
			continue
		}

		sqzOn[i] = lowerBB[i] > lowerKC[i] && upperBB[i] < upperKC[i] && isVolatile[i]
	}

	return val, sqzOn
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

func atr(high, low, close []float64, period int) []float64 {
	tr := trueRange(high, low, close)
	atr := make([]float64, len(tr))
	sum := 0.0
	for i := 0; i < len(tr); i++ {
		sum += tr[i]
		if i >= period {
			sum -= tr[i-period]
		}
		if i >= period-1 {
			atr[i] = sum / float64(period)
		}
	}
	return atr
}

func linearRegression(data []float64, period int) float64 {
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
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

func max(data []float64) float64 {
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func min(data []float64) float64 {
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}
