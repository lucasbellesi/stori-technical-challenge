package email

import (
	"stori-technical-challenge/config"

	gomail "gopkg.in/gomail.v2"
)

func SendEmail(subject, body, toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUser, config.AppConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
