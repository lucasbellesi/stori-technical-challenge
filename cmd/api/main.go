package main

import (
	"fmt"
	"log"
	"stori-technical-challenge/config"
	"stori-technical-challenge/pkg/db"
	"stori-technical-challenge/pkg/email"
	"stori-technical-challenge/pkg/transactions"
)

const SendEmailTo = "alejobellesi@hotmail.com"

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = db.InitDB()
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
	body := fmt.Sprintf(`
    <div class="container">
        <div class="row">
            <div class="col-12">
                <p><strong>Total balance:</strong> %.2f</p>
            </div>
        </div>
        <div class="row">
            <div class="col-12">
                <table class="table table-striped summary-table">
                    <thead>
                        <tr>
                            <th>Month</th>
                            <th>Number of Transactions</th>
                            <th>Average Credit</th>
                            <th>Average Debit</th>
                        </tr>
                    </thead>
                    <tbody>`, totalBalance)

	for month, data := range summary {
		body += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%.0f</td>
                <td>%.2f</td>
                <td>%.2f</td>
            </tr>`,
			month, data["num_transactions"], data["avg_credit"], data["avg_debit"])
	}

	body += `
                    </tbody>
                </table>
            </div>
        </div>
    </div>`

	err = email.SendEmail(subject, body, SendEmailTo)
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
