package indicator

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"strconv"
	"testing"
)

type TestData struct {
	ClosePrices     []float64
	HighPrices      []float64
	LowPrices       []float64
	ExpectedLines   []float64
	ExpectedSqueeze []int
	Nums            []float64
}

func readCSV(filename string) (*TestData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var testData TestData
	for i, row := range data {
		if i == 0 {
			continue
		}
		closeVal, _ := strconv.ParseFloat(row[4], 64)
		highVal, _ := strconv.ParseFloat(row[2], 64)
		lowVal, _ := strconv.ParseFloat(row[3], 64)
		lineVal, _ := strconv.ParseFloat(row[5], 64)
		squeezeVal, _ := strconv.Atoi(row[6])

		testData.ClosePrices = append(testData.ClosePrices, closeVal)
		testData.HighPrices = append(testData.HighPrices, highVal)
		testData.LowPrices = append(testData.LowPrices, lowVal)
		testData.ExpectedLines = append(testData.ExpectedLines, lineVal)
		testData.ExpectedSqueeze = append(testData.ExpectedSqueeze, squeezeVal)
		testData.Nums = append(testData.Nums, float64(lineVal))
	}
	return &testData, nil
}

func isNaN(value float64) bool {
	return math.IsNaN(value) || value == 0
}

func runSqueezeTestWithLine(t *testing.T, filename string, bbLength, kcLength int, bbMult, kcMult float64, useTrueRange bool, minVolatility float64) {
	testData, err := readCSV(filename)
	if err != nil {
		t.Fatalf("Error reading CSV file: %v", err)
	}

	indicator := NewSqueezeIndicator(bbLength, kcLength, bbMult, kcMult, useTrueRange, minVolatility)
	sqzOn := indicator.CalculateSqueeze(testData.ClosePrices, testData.HighPrices, testData.LowPrices, testData.ExpectedLines)
	checkSqueeze := !isNaN(testData.ExpectedLines[0])

	for i := range sqzOn {
		if i < indicator.KCLength-1 {
			continue
		}

		if isNaN(testData.ExpectedLines[i]) {
			continue
		}

		if !isNaN(testData.ExpectedLines[i]) {
			assert.InEpsilon(t, testData.ExpectedLines[i], testData.Nums[i], 1e-6, "Line mismatch at index %d", i)
		}

		if !checkSqueeze {
			expectedSqueeze := testData.ExpectedSqueeze[i] == 1
			assert.Equal(t, expectedSqueeze, sqzOn[i], "Squeeze mismatch at index %d", i)
		}
	}
}

func TestSqueezeWithLineCSV(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		bbLength      int
		kcLength      int
		bbMult        float64
		kcMult        float64
		useTrueRange  bool
		minVolatility float64
	}{
		{
			name:          "BINANCE_1000BONKUSDT",
			filename:      "files/BINANCE_1000BONKUSDT.P, 1 (1).csv",
			bbLength:      20,
			bbMult:        1.5,
			kcLength:      20,
			kcMult:        1.2,
			useTrueRange:  true,
			minVolatility: 0.0001,
		},
		{
			name:          "BITSTAMP_BTCUSD",
			filename:      "files/BITSTAMP_BTCUSD.csv",
			bbLength:      10,
			kcLength:      10,
			bbMult:        1.5,
			kcMult:        1.2,
			useTrueRange:  true,
			minVolatility: 0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSqueezeTestWithLine(t, tt.filename, tt.bbLength, tt.kcLength, tt.bbMult, tt.kcMult, tt.useTrueRange, tt.minVolatility)
		})
	}
}
