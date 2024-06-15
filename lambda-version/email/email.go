package email

import (
	"bytes"
	"log"
	"os"
	"strconv"

	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmail      string
}

var AppConfig Config

func LoadConfig() error {
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT value: %v", err)
		return err
	}

	AppConfig = Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		ToEmail:      os.Getenv("TO_EMAIL"),
	}

	return nil
}

type Summary struct {
	NumTransactions int
	AvgCredit       float64
	AvgDebit        float64
}

type EmailSender interface {
	SendEmail(subject, body, toEmail string) error
}

type SMTPSender struct{}

type EmailData struct {
	TotalBalance    float64
	NumTransactions map[string]int
	AvgDebitAmount  float64
	AvgCreditAmount float64
}

func (s SMTPSender) SendEmail(subject, body, toEmail string) error {
	err := LoadConfig()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", AppConfig.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(AppConfig.SMTPHost, AppConfig.SMTPPort, AppConfig.SMTPUser, AppConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func LoadTemplate(templatePath string, data EmailData) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return tpl.String(), nil
}

func GenerateEmailData(totalBalance float64, summary map[string]Summary, avgDebit float64, avgCredit float64) EmailData {
	emailData := EmailData{
		TotalBalance:    totalBalance,
		NumTransactions: make(map[string]int),
		AvgDebitAmount:  avgDebit,
		AvgCreditAmount: avgCredit,
	}

	for month, data := range summary {
		emailData.NumTransactions[month] = data.NumTransactions
	}

	return emailData
}
