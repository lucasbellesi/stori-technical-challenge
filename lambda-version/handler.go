package main

import (
	"context"
	"fmt"
	"log"
	"stori-technical-challenge/lambda-version/email"
	"stori-technical-challenge/lambda-version/transactions"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	ToEmail string               `json:"toEmail"`
	S3Event events.S3EventRecord `json:"s3Event"`
}

func handleRequest(ctx context.Context, req Request) {
	bucket := req.S3Event.S3.Bucket.Name
	key := req.S3Event.S3.Object.Key

	fmt.Printf("Processing file: s3://%s/%s\n", bucket, key)

	totalBalance, summary, avgDebit, avgCredit, err := transactions.ProcessTransactionsFromS3(bucket, key)
	if err != nil {
		log.Fatalf("Error processing transactions: %v", err)
		return
	}

	emailData := email.GenerateEmailData(totalBalance, summary, avgDebit, avgCredit)
	body, err := email.LoadTemplate("email_template.html", emailData)
	if err != nil {
		log.Fatalf("Error loading email template: %v", err)
		return
	}

	emailSender := email.SMTPSender{}
	err = emailSender.SendEmail("Transaction Summary", body, req.ToEmail)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
		return
	}

	log.Println("Email sent successfully!")
}

func main() {
	lambda.Start(handleRequest)
}
