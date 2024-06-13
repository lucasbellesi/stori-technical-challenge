package main

import (
	"fmt"
	"log"
	"stori-technical-challenge/config"
	"stori-technical-challenge/pkg/db"
	"stori-technical-challenge/pkg/email"
	"stori-technical-challenge/pkg/transactions"
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

	filePath := "transactions.csv"
	totalBalance, summary, err := transactions.ProcessTransactions(filePath)
	if err != nil {
		log.Fatalf("Error processing transactions: %v", err)
	}

	for month, data := range summary {
		for _, amount := range data {
			date := month + "-01" // Placeholder date
			transaction := db.Transaction{
				Date:   date,
				Amount: amount,
			}
			err = db.SaveTransaction(transaction)
			if err != nil {
				log.Fatalf("Error saving transaction: %v", err)
			}
		}
	}

	subject := "Transaction Summary"
	body := fmt.Sprintf("Total balance: %.2f\n", totalBalance)
	for month, data := range summary {
		body += fmt.Sprintf("%s: %.0f transactions, avg credit: %.2f, avg debit: %.2f\n",
			month, data["num_transactions"], data["avg_credit"], data["avg_debit"])
	}

	err = email.SendEmail(subject, body, "recipient@example.com")
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
