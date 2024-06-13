package db

import (
	"database/sql"
	"fmt"

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
