package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Transaction struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./transactions.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        amount REAL NOT NULL
    );`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func SaveTransaction(transaction Transaction) error {
	insertQuery := `INSERT INTO transactions (date, amount) VALUES (?, ?)`
	_, err := DB.Exec(insertQuery, transaction.Date, transaction.Amount)
	if err != nil {
		return fmt.Errorf("error saving transaction: %v", err)
	}
	return nil
}

func GetAllTransactions() ([]Transaction, error) {
	rows, err := DB.Query("SELECT date, amount FROM transactions")
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %v", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.Date, &transaction.Amount); err != nil {
			return nil, fmt.Errorf("error scanning transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func SaveTransactionsFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening transactions file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading transactions file: %v", err)
	}

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		dateStr := record[1]
		amountStr := record[2]

		// Parse the date
		date, err := time.Parse("1/2", dateStr)
		if err != nil {
			return fmt.Errorf("error parsing date: %v", err)
		}

		// Assume current year
		currentYear := time.Now().Year()
		date = date.AddDate(currentYear-date.Year(), 0, 0)
		formattedDate := date.Format("2006-01-02")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return fmt.Errorf("error parsing transaction amount: %v", err)
		}
		transaction := Transaction{
			Date:   formattedDate,
			Amount: amount,
		}
		err = SaveTransaction(transaction)
		if err != nil {
			return fmt.Errorf("error saving transaction: %v", err)
		}
	}
	return nil
}
