package transactions

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	reader CSVReader
}

func NewProcessor(reader CSVReader) *Processor {
	return &Processor{reader: reader}
}

func (p *Processor) ProcessTransactions(filePath string) (float64, map[string]Summary, float64, float64, error) {
	records, err := p.reader.Read(filePath)
	if err != nil {
		return 0, nil, 0, 0, fmt.Errorf("error reading transactions: %w", err)
	}

	transactions := p.parseTransactions(records)

	totalBalance, avgCredit, avgDebit, err := p.calculateTotalsAndAverages(transactions)
	if err != nil {
		return 0, nil, 0, 0, err
	}

	summary := p.generateSummary(transactions)

	return totalBalance, summary, avgDebit, avgCredit, nil
}

func (p *Processor) parseTransactions(records [][]string) map[string][]float64 {
	transactions := make(map[string][]float64)
	currentYear := time.Now().Year()

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		dateStr := strings.TrimSpace(record[1])
		transactionStr := strings.TrimSpace(record[2])
		date, err := time.Parse("1/2", dateStr)
		if err != nil {
			continue
		}
		date = date.AddDate(currentYear-date.Year(), 0, 0)

		amount, err := strconv.ParseFloat(transactionStr, 64)
		if err != nil {
			continue
		}

		month := date.Format("2006-01")
		transactions[month] = append(transactions[month], amount)
	}

	return transactions
}

func (p *Processor) calculateTotalsAndAverages(transactions map[string][]float64) (float64, float64, float64, error) {
	var totalBalance, totalCredit, totalDebit float64
	var numCredits, numDebits int

	for _, amounts := range transactions {
		for _, amount := range amounts {
			totalBalance += amount
			if amount > 0 {
				totalCredit += amount
				numCredits++
			} else {
				totalDebit += amount
				numDebits++
			}
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

	return totalBalance, avgCredit, avgDebit, nil
}

func (p *Processor) generateSummary(transactions map[string][]float64) map[string]Summary {
	summary := make(map[string]Summary)

	for date, amounts := range transactions {
		var totalCredit, totalDebit float64
		var numCredits, numDebits int

		for _, amount := range amounts {
			if amount > 0 {
				totalCredit += amount
				numCredits++
			} else {
				totalDebit += amount
				numDebits++
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

		summary[date] = Summary{
			NumTransactions: len(amounts),
			AvgCredit:       avgCredit,
			AvgDebit:        avgDebit,
		}
	}

	return summary
}
