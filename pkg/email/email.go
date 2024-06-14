package email

import (
	"bytes"
	"fmt"
	"html/template"
	"stori-technical-challenge/config"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(subject, body string) error
}

type SMTPSender struct{}

type EmailData struct {
	TotalBalance float64
	Summary      []MonthSummary
}

type MonthSummary struct {
	Month           string
	NumTransactions int
	AvgCredit       float64
	AvgDebit        float64
}

func (s SMTPSender) SendEmail(subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", config.AppConfig.ToEmail)
	m.SetHeader("Subject", subject)

	logoPath := "Stori_Logo_2023-min.png" // Ruta al archivo del logo
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
