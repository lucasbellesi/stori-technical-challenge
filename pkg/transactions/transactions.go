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

	for i, record := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		dateStr := strings.TrimSpace(record[1])
		transactionStr := strings.TrimSpace(record[2])

		// Convert date format from "7/15" to "2006-07-15" using the current year
		currentYear := time.Now().Year()
		date, err := time.Parse("1/2", dateStr)
		if err != nil {
			return 0, nil, err
		}
		date = date.AddDate(currentYear-date.Year(), 0, 0)

		amount, err := strconv.ParseFloat(transactionStr, 64)
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
		numCredits := 0.0
		totalDebits := 0.0
		numDebits := 0.0
		for _, amount := range amounts {
			if amount > 0 {
				totalCredits += amount
				numCredits++
			} else {
				totalDebits += amount
				numDebits++
			}
		}
		avgCredit := totalCredits / numCredits
		avgDebit := totalDebits / numDebits
		summary[month] = map[string]float64{
			"num_transactions": numTransactions,
			"avg_credit":       avgCredit,
			"avg_debit":        avgDebit,
		}
	}

	return totalBalance, summary, nil
}
