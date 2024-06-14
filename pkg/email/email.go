package email

import (
	"fmt"
	"stori-technical-challenge/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(subject, body, toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.FromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

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
			<img src="Stori_Logo_2023.png" alt="company logo"/>
            <h2>Transaction Summary</h2>
            %s
        </div>
    </body>
    </html>
    `, body)

	m.SetBody("text/html", htmlBody)
	m.Embed("Stori_Logo_2023.png")

	d := gomail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUser, config.AppConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
