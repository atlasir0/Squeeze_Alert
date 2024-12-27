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
	range_ := trueRange(high, low, close)
	rangema := sma(range_, si.KCLength)
	upperKC := add(ma, multiply(rangema, si.KCMult))
	lowerKC := subtract(ma, multiply(rangema, si.KCMult))

	atrValues := atr(high, low, close, si.KCLength)

	highestHighs := make([]float64, n)
	lowestLows := make([]float64, n)
	for i := 0; i < n; i++ {
		startIdx := int(math.Max(0, float64(i-si.KCLength+1)))
		highSlice := high[startIdx : i+1]
		lowSlice := low[startIdx : i+1]
		highestHighs[i] = max(highSlice)
		lowestLows[i] = min(lowSlice)
	}

	linregInput := make([]float64, n)
	for i := si.KCLength - 1; i < n; i++ {
		hlAvg := (highestHighs[i] + lowestLows[i]) / 2
		finalAvg := (hlAvg + ma[i]) / 2
		linregInput[i] = round(close[i]-finalAvg, 12)
	}

	val := make([]float64, n)
	for i := si.KCLength - 1; i < n; i++ {
		windowStart := int(math.Max(0, float64(i-si.KCLength+1)))
		window := linregInput[windowStart : i+1]
		val[i] = linearRegression(window, len(window))
	}

	sqzOn := make([]bool, n)
	for i := range close {
		if i < si.KCLength-1 {
			continue
		}

		for i := range close {
			if i < si.KCLength-1 {
				continue
			}

			isVolatile := atrValues[i] > si.MinVolatility

			lowerBBRounded := round(lowerBB[i], 12)
			upperBBRounded := round(upperBB[i], 12)
			lowerKCRounded := round(lowerKC[i], 12)
			upperKCRounded := round(upperKC[i], 12)
			bbInsideKC := (lowerBBRounded >= lowerKCRounded && upperBBRounded <= upperKCRounded)
			sqzOn[i] = bbInsideKC && isVolatile
		}

	}

	return val, sqzOn
}

func linearRegression(data []float64, period int) float64 {
	if len(data) < 2 {
		return 0
	}

	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	n := float64(len(data))

	for i := 0; i < len(data); i++ {
		x := float64(i)
		y := data[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	lastValue := intercept + slope*(n-1)
	return round(lastValue, 12)
}
func (si *SqueezeIndicator) CalculateSqueeze(close, high, low, lines []float64) []bool {
	n := len(close)
	if n == 0 {
		return []bool{}
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

	atrValues := atr(high, low, close, si.KCLength)

	sqzOn := make([]bool, n)
	for i := range close {
		if i < si.KCLength-1 {
			continue
		}
		isVolatile := atrValues[i] > si.MinVolatility

		lowerBBRounded := round(lowerBB[i], 12)
		upperBBRounded := round(upperBB[i], 12)
		lowerKCRounded := round(lowerKC[i], 12)
		upperKCRounded := round(upperKC[i], 12)

		sqzOn[i] = lowerBBRounded > lowerKCRounded && upperBBRounded < upperKCRounded && isVolatile

	}

	return sqzOn
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

func round(value float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(value*factor) / factor
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

func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
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

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}
