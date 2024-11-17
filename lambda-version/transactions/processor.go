package transactions

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Summary struct {
	NumTransactions int
	AvgCredit       float64
	AvgDebit        float64
}

// S3Client defines an interface for S3 operations.
type S3Client interface {
	GetObject(bucket, key string) (*s3.GetObjectOutput, error)
}

type DefaultS3Client struct{}

func (c *DefaultS3Client) GetObject(bucket, key string) (*s3.GetObjectOutput, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	return svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

// ProcessTransactionsFromS3 processes a CSV file from S3.
func ProcessTransactionsFromS3(bucket, key string, client S3Client) (float64, map[string]Summary, float64, float64, error) {
	obj, err := client.GetObject(bucket, key)
	if err != nil {
		return 0, nil, 0, 0, fmt.Errorf("error getting object from S3: %w", err)
	}
	defer obj.Body.Close()

	// Convert obj.Body (io.ReadCloser) into a buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(obj.Body); err != nil {
		return 0, nil, 0, 0, fmt.Errorf("error reading object body: %w", err)
	}

	// Create a bytes.Reader from the buffer
	reader := bytes.NewReader(buf.Bytes())

	records, err := readCSV(reader)
	if err != nil {
		return 0, nil, 0, 0, err
	}

	transactions := parseTransactions(records)
	return calculateMetrics(transactions)
}

func readCSV(dataStream *bytes.Reader) ([][]string, error) {
	reader := csv.NewReader(dataStream)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}
	return records, nil
}

func parseTransactions(records [][]string) map[string][]float64 {
	transactions := make(map[string][]float64)
	currentYear := time.Now().Year()

	for i, record := range records {
		if i == 0 {
			continue // Skip headers
		}
		date, amount, err := parseRecord(record)
		if err != nil {
			continue
		}
		month := fmt.Sprintf("%d-%s", currentYear, date)
		transactions[month] = append(transactions[month], amount)
	}
	return transactions
}

func parseRecord(record []string) (string, float64, error) {
	if len(record) < 3 {
		return "", 0, fmt.Errorf("invalid record format")
	}

	date := strings.TrimSpace(record[1])
	amount, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing amount: %w", err)
	}

	return date, amount, nil
}

func calculateMetrics(transactions map[string][]float64) (float64, map[string]Summary, float64, float64, error) {
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

func calculateSummary(transactions map[string][]float64) map[string]Summary {
	summary := make(map[string]Summary)
	for month, amounts := range transactions {
		numTransactions := len(amounts)
		totalCredits, totalDebits, numCredits, numDebits := 0.0, 0.0, 0, 0

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
