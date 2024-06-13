package transactions

import (
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
		return 0, nil, 0, 0, err
	}

	transactions := p.parseTransactions(records)
	totalBalance, totalCredit, totalDebit, numCredits, numDebits := p.calculateTotals(transactions)
	summary := p.calculateSummary(transactions)

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

func (p *Processor) calculateTotals(transactions map[string][]float64) (float64, float64, float64, int, int) {
	totalBalance := 0.0
	totalCredit := 0.0
	totalDebit := 0.0
	numCredits := 0
	numDebits := 0

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

	return totalBalance, totalCredit, totalDebit, numCredits, numDebits
}

func (p *Processor) calculateSummary(transactions map[string][]float64) map[string]Summary {
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

	return summary
}
