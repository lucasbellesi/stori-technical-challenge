package transactions

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

type Summary struct {
	NumTransactions int
	AvgCredit       float64
	AvgDebit        float64
}

func ProcessTransactions(filePath string) (float64, map[string]Summary, float64, float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, 0, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, nil, 0, 0, err
	}

	transactions := make(map[string][]float64)
	totalBalance := 0.0
	totalCredit := 0.0
	totalDebit := 0.0
	numCredits := 0
	numDebits := 0

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
			return 0, nil, 0, 0, err
		}
		date = date.AddDate(currentYear-date.Year(), 0, 0)

		amount, err := strconv.ParseFloat(transactionStr, 64)
		if err != nil {
			return 0, nil, 0, 0, err
		}

		month := date.Format("2006-01")
		transactions[month] = append(transactions[month], amount)
		totalBalance += amount

		if amount > 0 {
			totalCredit += amount
			numCredits++
		} else {
			totalDebit += amount
			numDebits++
		}
	}

	summary := make(map[string]Summary)
	for month, amounts := range transactions {
		numTransactions := len(amounts)
		totalCredits := 0.0
		totalDebits := 0.0
		numCredits := 0
		numDebits := 0
		for _, amount := range amounts {
			if amount > 0 {
				totalCredits += amount
				numCredits++
			} else {
				totalDebits += amount
				numDebits++
			}
		}
		avgCredit := 0.0
		if numCredits > 0 {
			avgCredit = totalCredits / float64(numCredits)
		}
		avgDebit := 0.0
		if numDebits > 0 {
			avgDebit = totalDebits / float64(numDebits)
		}
		summary[month] = Summary{
			NumTransactions: numTransactions,
			AvgCredit:       avgCredit,
			AvgDebit:        avgDebit,
		}
	}

	avgCredit := 0.0
	if numCredits > 0 {
		avgCredit = totalCredit / float64(numCredits)
	}
	avgDebit := 0.0
	if numDebits > 0 {
		avgDebit = totalDebit / float64(numDebits)
	}

	return totalBalance, summary, avgDebit, avgCredit, nil
}
