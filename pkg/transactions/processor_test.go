package transactions

import (
	"reflect"
	"testing"
)

func TestParseTransactions(t *testing.T) {
	reader := &MockCSVReader{
		records: [][]string{
			{"2023-11-01", "100.0"},
			{"2023-11-01", "-50.0"},
			{"2023-11-02", "200.0"},
		},
	}
	processor := NewProcessor(reader)

	transactions := processor.parseTransactions(reader.records)

	expected := map[string][]float64{
		"2023-11-01": {100.0, -50.0},
		"2023-11-02": {200.0},
	}

	if !reflect.DeepEqual(transactions, expected) {
		t.Errorf("expected %v, got %v", expected, transactions)
	}
}

func TestCalculateTotalsAndAverages(t *testing.T) {
	transactions := map[string][]float64{
		"2023-11-01": {100.0, -50.0},
		"2023-11-02": {200.0, -100.0},
	}

	processor := &Processor{}
	totalBalance, avgCredit, avgDebit, err := processor.calculateTotalsAndAverages(transactions)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if totalBalance != 150.0 {
		t.Errorf("expected total balance to be 150.0, got %f", totalBalance)
	}

	if avgCredit != 150.0 {
		t.Errorf("expected avg credit to be 150.0, got %f", avgCredit)
	}

	if avgDebit != -75.0 {
		t.Errorf("expected avg debit to be -75.0, got %f", avgDebit)
	}
}

func TestGenerateSummary(t *testing.T) {
	transactions := map[string][]float64{
		"2023-11-01": {100.0, -50.0},
		"2023-11-02": {200.0, -100.0},
	}

	processor := &Processor{}
	summary := processor.generateSummary(transactions)

	expected := map[string]Summary{
		"2023-11-01": {NumTransactions: 2, AvgCredit: 100.0, AvgDebit: -50.0},
		"2023-11-02": {NumTransactions: 2, AvgCredit: 200.0, AvgDebit: -100.0},
	}

	if !reflect.DeepEqual(summary, expected) {
		t.Errorf("expected %v, got %v", expected, summary)
	}
}

type MockCSVReader struct {
	records [][]string
}

func (m *MockCSVReader) Read(filePath string) ([][]string, error) {
	return m.records, nil
}
