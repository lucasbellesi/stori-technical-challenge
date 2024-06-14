package email

import (
	"fmt"
	"stori-technical-challenge/config"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(subject, body, toEmail string) error
}

type SMTPSender struct{}

func (s SMTPSender) SendEmail(subject, body, toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	logoPath := "Stori_Logo_2023-min.png" // Ruta al archivo del logo
	m.Embed(logoPath)

	htmlBody := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" rel="stylesheet">
        <style>
            .email-body {
                padding: 20px;
            }
            .logo {
                width: 100px;
                height: auto;
            }
            .summary-table {
                margin-top: 20px;
            }
        </style>
    </head>
    <body>
        <div class="email-body">
            <img src="cid:logo" class="logo" alt="Company Logo">
            <h2>Transaction Summary</h2>
            %s
        </div>
    </body>
    </html>`, body)

	m.SetBody("text/html", htmlBody)
	m.Embed(logoPath, gomail.SetHeader(map[string][]string{"Content-ID": {"<logo>"}}))

	d := gomail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUser, config.AppConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
