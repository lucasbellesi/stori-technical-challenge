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
		return fmt.Errorf("Error opening database: %v", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        amount REAL NOT NULL
    );`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("Error creating table: %v", err)
	}

	return nil
}

func SaveTransaction(transaction Transaction) error {
	insertQuery := `INSERT INTO transactions (date, amount) VALUES (?, ?)`
	_, err := DB.Exec(insertQuery, transaction.Date, transaction.Amount)
	if err != nil {
		return fmt.Errorf("Error saving transaction: %v", err)
	}
	return nil
}

func GetTransactionSummary() (map[string]map[string]float64, error) {
	summary := make(map[string]map[string]float64)
	query := `SELECT
                SUBSTR(date, 1, 7) AS month,
                COUNT(*) AS num_transactions,
                AVG(CASE WHEN amount > 0 THEN amount ELSE NULL END) AS avg_credit,
                AVG(CASE WHEN amount < 0 THEN amount ELSE NULL END) AS avg_debit
              FROM transactions
              GROUP BY month`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving transaction summary: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var numTransactions float64
		var avgCredit, avgDebit float64

		if err := rows.Scan(&month, &numTransactions, &avgCredit, &avgDebit); err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}

		summary[month] = map[string]float64{
			"num_transactions": numTransactions,
			"avg_credit":       avgCredit,
			"avg_debit":        avgDebit,
		}
	}

	return summary, nil
}
