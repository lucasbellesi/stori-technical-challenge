package main

import (
	"context"
	"fmt"
	"log"
	"stori-technical-challenge/lambda-version/email"
	"stori-technical-challenge/lambda-version/transactions"
	"stori-technical-challenge/lambda-version/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	ToEmail string               `json:"toEmail"`
	S3Event events.S3EventRecord `json:"s3Event"`
}

func handleRequest(ctx context.Context, req Request) (string, error) {
	// Validar entrada
	if req.ToEmail == "" {
		return "", fmt.Errorf("missing 'ToEmail' field")
	}
	if req.S3Event.S3.Bucket.Name == "" || req.S3Event.S3.Object.Key == "" {
		return "", fmt.Errorf("invalid S3 event: bucket or key is missing")
	}

	bucket := req.S3Event.S3.Bucket.Name
	key := req.S3Event.S3.Object.Key
	log.Printf("Processing file: s3://%s/%s\n", bucket, key)

	// Crear cliente S3
	s3Client := &transactions.DefaultS3Client{}

	// Procesar transacciones desde S3
	totalBalance, summary, avgDebit, avgCredit, err := transactions.ProcessTransactionsFromS3(bucket, key, s3Client)
	if err != nil {
		return "", fmt.Errorf("error processing transactions: %w", err)
	}

	emailSummary := utils.ConvertTransactionSummaryToEmailSummary(summary)

	// Generar contenido del correo
	emailData := email.GenerateEmailData(totalBalance, emailSummary, avgDebit, avgCredit)
	body, err := email.LoadTemplate("email_template.html", emailData)
	if err != nil {
		return "", fmt.Errorf("error loading email template: %w", err)
	}

	// Enviar correo
	emailSender := email.SMTPSender{}
	err = emailSender.SendEmail("Stori - Transaction Summary", body, req.ToEmail)
	if err != nil {
		return "", fmt.Errorf("error sending email: %w", err)
	}

	log.Println("Email sent successfully")
	return "Success", nil
}

func main() {
	lambda.Start(handleRequest)
}
