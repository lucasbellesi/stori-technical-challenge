package db

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB

type Transaction struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func InitDB() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return fmt.Errorf("error creating AWS session: %v", err)
	}

	db = dynamodb.New(sess)
	return nil
}

func SaveTransaction(transaction Transaction) error {
	av, err := dynamodbattribute.MarshalMap(transaction)
	if err != nil {
		return fmt.Errorf("error marshalling transaction: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return fmt.Errorf("error putting item into DynamoDB: %v", err)
	}

	return nil
}

func GetTransactionSummary() (map[string]map[string]float64, error) {
	summary := make(map[string]map[string]float64)
	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB table: %v", err)
	}

	transactions := []Transaction{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &transactions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling DynamoDB scan result: %v", err)
	}

	for _, transaction := range transactions {
		month := transaction.Date[:7] // Extract YYYY-MM
		if _, exists := summary[month]; !exists {
			summary[month] = map[string]float64{"num_transactions": 0, "total_credits": 0, "total_debits": 0}
		}
		summary[month]["num_transactions"]++
		if transaction.Amount > 0 {
			summary[month]["total_credits"] += transaction.Amount
		} else {
			summary[month]["total_debits"] += transaction.Amount
		}
	}

	for _, data := range summary {
		data["avg_credit"] = data["total_credits"] / data["num_transactions"]
		data["avg_debit"] = data["total_debits"] / data["num_transactions"]
	}

	return summary, nil
}
