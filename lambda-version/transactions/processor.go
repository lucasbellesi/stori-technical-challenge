package transactions

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"stori-technical-challenge/lambda-version/email"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ProcessTransactionsFromS3(bucket, key string) (float64, map[string]email.Summary, float64, float64, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, nil, 0, 0, fmt.Errorf("error getting object from S3: %w", err)
	}
	defer obj.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(obj.Body)

	reader := csv.NewReader(buf)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, nil, 0, 0, fmt.Errorf("error reading CSV: %w", err)
	}

	transactions := parseTransactions(records)
	totalBalance, totalCredit, totalDebit, numCredits, numDebits := calculateTotals(transactions)
	summary := calculateSummary(transactions)

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

func parseTransactions(records [][]string) map[string][]float64 {
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

func calculateTotals(transactions map[string][]float64) (float64, float64, float64, int, int) {
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

func calculateSummary(transactions map[string][]float64) map[string]email.Summary {
	summary := make(map[string]email.Summary)

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

		summary[month] = email.Summary{
			NumTransactions: numTransactions,
			AvgCredit:       avgCredit,
			AvgDebit:        avgDebit,
		}
	}

	return summary
}
