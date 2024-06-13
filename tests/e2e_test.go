package tests

import (
	"os"
	"stori-technical-challenge/config"
	"stori-technical-challenge/pkg/db"
	"stori-technical-challenge/pkg/email"
	"stori-technical-challenge/pkg/transactions"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestCSV(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("Id,Date,Transaction\n")
	file.WriteString("0,7/15,+60.5\n")
	file.WriteString("1,7/28,-10.3\n")
	file.WriteString("2,8/2,-20.46\n")
	file.WriteString("3,8/13,+10\n")
}

func TestE2E(t *testing.T) {
	// Configurar variables de entorno directamente
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "1025") // Puerto por defecto de MailHog
	os.Setenv("SMTP_USER", "")
	os.Setenv("SMTP_PASSWORD", "")
	os.Setenv("FROM_EMAIL", "test@example.com")

	// Cargar la configuración sin depender del archivo .env
	config.AppConfig = config.Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     1025, // Convertir el puerto de cadena a entero
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
	}

	err := db.InitDB()
	assert.NoError(t, err, "Error initializing database")

	filePath := "test_transactions.csv"
	createTestCSV(filePath)
	defer os.Remove(filePath)

	reader := transactions.DefaultCSVReader{}
	processor := transactions.NewProcessor(reader)

	totalBalance, summary, avgDebit, avgCredit, err := processor.ProcessTransactions(filePath)
	assert.NoError(t, err, "Error processing transactions")

	// Guardar todas las transacciones individuales en la base de datos
	err = db.SaveTransactionsFromCSV(filePath)
	assert.NoError(t, err, "Error saving transactions from CSV")

	emailData := email.EmailData{
		TotalBalance:    totalBalance,
		NumTransactions: make(map[string]int),
		AvgDebitAmount:  avgDebit,
		AvgCreditAmount: avgCredit,
	}

	for month, data := range summary {
		emailData.NumTransactions[month] = data.NumTransactions
	}

	body, err := email.LoadTemplate("email_template_test.html", emailData)
	assert.NoError(t, err, "Error loading email template")

	// Omitir el envío de correo en la prueba
	// emailSender := email.SMTPSender{}
	// err = emailSender.SendEmail(subject, body, "recipient@example.com")
	// assert.NoError(t, err, "Error sending email")

	// Verificar que las transacciones se guardaron en la base de datos
	transactions, err := db.GetAllTransactions()
	assert.NoError(t, err, "Error retrieving transactions")
	assert.NotEmpty(t, transactions, "No transactions found")

	lastTransaction := transactions[len(transactions)-1]
	assert.Equal(t, "2024-08-13", lastTransaction.Date, "Last transaction date does not match")
	assert.Equal(t, 10.0, lastTransaction.Amount, "Last transaction amount does not match")

	// Verificar detalles del email
	assert.Contains(t, body, "Total balance is 39.74", "Email body does not contain correct total balance")
	assert.Contains(t, body, "Number of transactions in 2024-07: 2", "Email body does not contain correct transaction count for July")
	assert.Contains(t, body, "Number of transactions in 2024-08: 2", "Email body does not contain correct transaction count for August")
	assert.Contains(t, body, "Average debit amount: -15.38", "Email body does not contain correct average debit amount")
	assert.Contains(t, body, "Average credit amount: 35.25", "Email body does not contain correct average credit amount")
}
