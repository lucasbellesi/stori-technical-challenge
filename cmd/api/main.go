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

	totalBalance, summary, avgDebit, avgCredit, err := transactions.ProcessTransactions(FilePath)
	if err != nil {
		log.Fatalf("Error processing transactions: %v", err)
	}

	for month, data := range summary {
		// Guardar sólo la cantidad total de transacciones en la base de datos
		transaction := db.Transaction{
			Date:   month + "-01",                 // Placeholder date
			Amount: float64(data.NumTransactions), // Placeholder amount
		}
		err = db.SaveTransaction(transaction)
		if err != nil {
			log.Fatalf("Error saving transactions: %v", err)
		}
	}

	emailData := email.EmailData{
		TotalBalance:    totalBalance,
		NumTransactions: make(map[string]int),
		AvgDebitAmount:  avgDebit,
		AvgCreditAmount: avgCredit,
	}

	for month, data := range summary {
		emailData.NumTransactions[month] = data.NumTransactions
	}

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
	transactions, err := db.GetAllTransactions()
	if err != nil {
		log.Fatalf("Error retrieving transactions: %v", err)
	}
	for _, transaction := range transactions {
		fmt.Printf("Date: %s, Amount: %.2f\n", transaction.Date, transaction.Amount)
	}
}
