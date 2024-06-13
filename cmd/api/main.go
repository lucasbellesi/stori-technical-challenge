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
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	reader := transactions.DefaultCSVReader{}
	processor := transactions.NewProcessor(reader)

	totalBalance, summary, avgDebit, avgCredit, err := processor.ProcessTransactions(FilePath)
	if err != nil {
		log.Fatalf("Error processing transactions: %v", err)
	}

	// Guardar todas las transacciones individuales en la base de datos
	err = db.SaveTransactionsFromCSV(FilePath)
	if err != nil {
		log.Fatalf("Error saving transactions to the database: %v", err)
	}

	emailData := email.GenerateEmailData(totalBalance, summary, avgDebit, avgCredit)

	body, err := email.LoadTemplate(FilePathEmailTemplate, emailData)
	if err != nil {
		log.Fatalf("Error loading email template: %v", err)
	}

	emailSender := email.SMTPSender{}
	err = emailSender.SendEmail(Subject, body)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	fmt.Println("Email sent successfully!")

	// Retrieve and print all transactions
	fmt.Println("The database:")
	transactions, err := db.GetAllTransactions()
	if err != nil {
		log.Fatalf("Error retrieving transactions: %v", err)
	}
	for _, transaction := range transactions {
		fmt.Printf("Date: %s, Amount: %.2f\n", transaction.Date, transaction.Amount)
	}
}
