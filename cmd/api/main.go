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
	config.LoadConfig()

	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	filePath := "txns.csv"
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
	emailData := email.EmailData{
		TotalBalance: totalBalance,
		Summary:      []email.MonthSummary{},
	}

	for month, data := range summary {
		monthSummary := email.MonthSummary{
			Month:           month,
			NumTransactions: int(data["num_transactions"]),
			AvgCredit:       data["avg_credit"],
			AvgDebit:        data["avg_debit"],
		}
		emailData.Summary = append(emailData.Summary, monthSummary)
	}

	body, err := email.LoadTemplate("pkg/email/email_template.html", emailData)
	if err != nil {
		log.Fatalf("Error loading email template: %v", err)
	}

	emailSender := email.SMTPSender{}
	err = emailSender.SendEmail(subject, body)
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
