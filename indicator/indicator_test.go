package indicator

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	ClosePrices     []float64
	HighPrices      []float64
	LowPrices       []float64
	ExpectedLines   []float64
	ExpectedSqueeze []int
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
	}

	return &testData, nil
}

func runSqueezeTest(t *testing.T, filename string) {
	testData, err := readCSV(filename)
	if err != nil {
		t.Fatalf("Ошибка при чтении CSV файла: %v", err)
	}

	indicator := NewSqueezeIndicator(20, 20, 2.0, 1.5, true, 0.001)
	values, sqzOn := indicator.Calculate(testData.ClosePrices, testData.HighPrices, testData.LowPrices)

	for i := range values {
		if i < indicator.KCLength-1 {
			continue
		}
		if !math.IsNaN(testData.ExpectedLines[i]) {
			assert.Equal(t, testData.ExpectedLines[i], values[i], "Ошибка в Line на индексе %d", i)
		}

		assert.Equal(t, testData.ExpectedSqueeze[i] == 1, sqzOn[i], "Ошибка в Squeeze на индексе %d", i)
	}
}

func TestSqueezeCSV(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "BINANCE_1000BONKUSDT",
			filename: "files/BINANCE_1000BONKUSDT.P, 1 (1).csv",
		},
		{
			name:     "BITSTAMP_BTCUSD",
			filename: "files/BITSTAMP_BTCUSD.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSqueezeTest(t, tt.filename)
		})
	}
}
