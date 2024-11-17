package email

import (
	"bytes"
	"fmt"
	"html/template"
	"stori-technical-challenge/config"
	"stori-technical-challenge/pkg/transactions"

	"gopkg.in/gomail.v2"
)

const LogoPath = "assets/Stori_Logo_2023-min.png"

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

func (s SMTPSender) SendEmail(subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", config.AppConfig.ToEmail)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)
	m.Embed(LogoPath, gomail.SetHeader(map[string][]string{"Content-ID": {"<logo>"}}))

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

func GenerateEmailData(totalBalance float64, summary map[string]transactions.Summary, avgDebit float64, avgCredit float64) EmailData {
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

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
