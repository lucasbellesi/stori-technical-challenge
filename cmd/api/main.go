package main

import (
	"fmt"
	"log"
	"stori-technical-challenge/config"
	"stori-technical-challenge/pkg/db"
	"stori-technical-challenge/pkg/email"
	"stori-technical-challenge/pkg/transactions"
)

const (
	Subject               = "Stori - Transaction Summary"
	FilePath              = "txns.csv"
	FilePathEmailTemplate = "pkg/email/email_template.html"
)

func main() {
	if err := initializeApp(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := processAndSendTransactions(); err != nil {
		log.Fatalf("Failed to process transactions: %v", err)
	}

	log.Println("Application finished successfully")
}

func initializeApp() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := db.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}

func processAndSendTransactions() error {
	reader := transactions.DefaultCSVReader{}
	processor := transactions.NewProcessor(reader)

	// Process transactions and get results
	totalBalance, summary, avgDebit, avgCredit, err := processor.ProcessTransactions(FilePath)
	if err != nil {
		return fmt.Errorf("processing transactions: %w", err)
	}

	// Save processed transactions to database
	if err := db.SaveTransactionsFromCSV(FilePath); err != nil {
		return fmt.Errorf("saving transactions to the database: %w", err)
	}

	// Send the summary email
	if err := sendSummaryEmail(totalBalance, summary, avgDebit, avgCredit); err != nil {
		return fmt.Errorf("sending summary email: %w", err)
	}

	return nil
}

func sendSummaryEmail(totalBalance float64, summary map[string]transactions.Summary, avgDebit, avgCredit float64) error {
	// Prepare email data
	emailData := email.EmailData{
		TotalBalance:    totalBalance,
		NumTransactions: map[string]int{}, // Populate with summary data
		AvgDebitAmount:  avgDebit,
		AvgCreditAmount: avgCredit,
	}

	// Populate NumTransactions from summary
	for date, details := range summary {
		emailData.NumTransactions[date] = details.NumTransactions
	}

	// Render email template
	body, err := email.RenderTemplate(FilePathEmailTemplate, emailData)
	if err != nil {
		return fmt.Errorf("rendering email template: %w", err)
	}

	// Send email using SMTPSender
	sender := email.SMTPSender{}
	if err := sender.SendEmail(Subject, body); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
