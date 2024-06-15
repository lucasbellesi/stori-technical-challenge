package email

import (
	"bytes"
	"fmt"
	"html/template"
	"stori-technical-challenge/config"

	"gopkg.in/gomail.v2"
)

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
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	logoPath := "assets/Stori_Logo_2023-min.png" // Ruta al archivo del logo
	m.Embed(logoPath)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUser, config.AppConfig.SMTPPassword)

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
