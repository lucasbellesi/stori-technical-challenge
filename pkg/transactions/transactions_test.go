package transactions_test

import (
	"os"
	"stori-technical-challenge/pkg/transactions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestCSV(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("Id,Date,Transaction\n")
	file.WriteString("0,7/15,+60.5\n")
	file.WriteString("1,7/28,-10.3\n")
	file.WriteString("2,8/2,-20.46\n")
	file.WriteString("3,8/13,+10\n")
}

func TestProcessTransactions(t *testing.T) {
	filePath := "test_transactions.csv"
	createTestCSV(filePath)
	defer os.Remove(filePath)

	totalBalance, summary, err := transactions.ProcessTransactions(filePath)
	assert.NoError(t, err, "Error processing transactions")
	assert.Equal(t, 39.74, totalBalance, "Total balance does not match")

	expectedSummary := map[string]map[string]float64{
		"2024-07": {
			"num_transactions": 2,
			"avg_credit":       60.5,
			"avg_debit":        -10.3,
		},
		"2024-08": {
			"num_transactions": 2,
			"avg_credit":       10.0,
			"avg_debit":        -20.46,
		},
	}

	assert.Equal(t, expectedSummary, summary, "Summary does not match")
}
