package transactions

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

func ProcessTransactions(filePath string) (float64, map[string]map[string]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, nil, err
	}

	transactions := make(map[string][]float64)
	totalBalance := 0.0

	for _, record := range records {
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return 0, nil, err
		}

		amountStr := strings.TrimSpace(record[1])
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return 0, nil, err
		}

		month := date.Format("2006-01")
		transactions[month] = append(transactions[month], amount)
		totalBalance += amount
	}

	summary := make(map[string]map[string]float64)
	for month, amounts := range transactions {
		numTransactions := float64(len(amounts))
		totalCredits := 0.0
		totalDebits := 0.0
		for _, amount := range amounts {
			if amount > 0 {
				totalCredits += amount
			} else {
				totalDebits += amount
			}
		}
		avgCredit := totalCredits / numTransactions
		avgDebit := totalDebits / numTransactions
		summary[month] = map[string]float64{
			"num_transactions": numTransactions,
			"avg_credit":       avgCredit,
			"avg_debit":        avgDebit,
		}
	}

	return totalBalance, summary, nil
}
