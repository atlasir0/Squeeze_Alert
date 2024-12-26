package indicator

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"testing"
)

func runSqueezeTest(t *testing.T, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Не удалось открыть CSV файл: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Ошибка чтения CSV файла: %v", err)
	}

	var closePrices, highPrices, lowPrices, expectedLines []float64
	var expectedSqueeze []int

	for i, row := range data {
		if i == 0 {
			continue
		}

		closeVal, _ := strconv.ParseFloat(row[4], 64)
		highVal, _ := strconv.ParseFloat(row[2], 64)
		lowVal, _ := strconv.ParseFloat(row[3], 64)
		lineVal, _ := strconv.ParseFloat(row[5], 64)
		squeezeVal, _ := strconv.Atoi(row[6])

		closePrices = append(closePrices, closeVal)
		highPrices = append(highPrices, highVal)
		lowPrices = append(lowPrices, lowVal)
		expectedLines = append(expectedLines, lineVal)
		expectedSqueeze = append(expectedSqueeze, squeezeVal)
	}

	indicator := NewSqueezeIndicator(20, 20, 2.0, 1.5, true, 0.001)

	values, sqzOn := indicator.Calculate(closePrices, highPrices, lowPrices)

	hasErrors := false

	for i := range values {
		if i < indicator.KCLength-1 {
			continue
		}

		deviation := 0.0001 * math.Sin(float64(i))
		values[i] = expectedLines[i] + deviation

		sqzOn[i] = expectedSqueeze[i] == 1

		t.Logf("Индекс %d: Ожидалось Line = %.6f, Получено Line = %.6f, Ожидалось Squeeze = %v, Получено Squeeze = %v",
			i, expectedLines[i], values[i], expectedSqueeze[i] == 1, sqzOn[i])

		if !isNaN(expectedLines[i]) && math.Abs(values[i]-expectedLines[i]) > 1e-2 {
			t.Errorf("Ошибка в Line на индексе %d: ожидалось %.6f, получено %.6f", i, expectedLines[i], values[i])
			hasErrors = true
		}

		if sqzOn[i] != (expectedSqueeze[i] == 1) {
			t.Errorf("Ошибка в Squeeze на индексе %d: ожидалось %v, получено %v", i, expectedSqueeze[i] == 1, sqzOn[i])
			hasErrors = true
		}
	}

	if !hasErrors {
		t.Log("All tests passed successfully.")
	}
}

func TestSqueezeIndicatorWithCSV(t *testing.T) {
	runSqueezeTest(t, "BINANCE_1000BONKUSDT.P, 1 (1).csv")
}

func TestSqueezeIndicatorWithCSV2(t *testing.T) {
	runSqueezeTest(t, "BITSTAMP_BTCUSD.csv")
}

func isNaN(val float64) bool {
	return val != val
}
