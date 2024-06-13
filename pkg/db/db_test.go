package db_test

import (
	"stori-technical-challenge/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	err := db.InitDB()
	assert.NoError(t, err, "Error initializing database")

	// Check if the transactions table exists
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='transactions';"
	row := db.DB.QueryRow(query)
	var tableName string
	err = row.Scan(&tableName)
	assert.NoError(t, err, "Error querying database")
	assert.Equal(t, "transactions", tableName, "Table 'transactions' does not exist")
}

func TestSaveTransaction(t *testing.T) {
	err := db.InitDB()
	assert.NoError(t, err, "Error initializing database")

	transaction := db.Transaction{
		Date:   "2024-06-14",
		Amount: 100.50,
	}

	err = db.SaveTransaction(transaction)
	assert.NoError(t, err, "Error saving transaction")

	// Retrieve the transaction from the database
	transactions, err := db.GetAllTransactions()
	assert.NoError(t, err, "Error retrieving transactions")
	assert.NotEmpty(t, transactions, "No transactions found")

	lastTransaction := transactions[len(transactions)-1]
	assert.Equal(t, transaction.Date, lastTransaction.Date, "Transaction date does not match")
	assert.Equal(t, transaction.Amount, lastTransaction.Amount, "Transaction amount does not match")
}
