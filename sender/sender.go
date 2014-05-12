package sender

import (
	"bytes"
	"encoding/json"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

const (
	endpoint = "email-smtp.us-east-1.amazonaws.com"
)

type configuration struct {
	SmtpUsername string
	SmtpPassword string
	FromAddress  string
	ToAddresses  []string
}

const emailTemplateBody = `From: {{fromEmail}}
To: {{toEmail}}
Subject: {{.Subject}}
Body:

{{.Name}} sent a message:

{{.Body}}

Yours truly,
Mr. Contact Form Robot`

var config *configuration
var emailTemplate *template.Template

func init() {
	config = &configuration{}
	file, err := os.Open("../server/config.json")
	if err != nil {
		file, err = os.Open("./config.json")
	}
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}

	emailTemplate = template.New("emailBody").Funcs(template.FuncMap{
		"fromEmail": func() string { return config.FromAddress },
		"toEmail":   func() string { return strings.Join(config.ToAddresses, ", ") },
	})
	emailTemplate = template.Must(emailTemplate.Parse(emailTemplateBody))
}

// Basic structure of an email message.  Name = User's name (like John) that they enter on a contact form.
type EmailMessage struct {
	Name    string
	Subject string
	Body    string
}

// Sends smtp email
func SendEmail(message EmailMessage, originHost string) error {
	var body bytes.Buffer
	err := emailTemplate.Execute(&body, message)
	if err != nil {
		return err
	}
	return smtpSendEmail(body.Bytes())
}

var smtpSendEmail = func(body []byte) error {
	auth := smtp.PlainAuth(
		"",
		config.SmtpUsername,
		config.SmtpPassword,
		endpoint,
	)

	return smtp.SendMail(
		endpoint+":587",
		auth,
		config.FromAddress,
		config.ToAddresses,
		body,
	)
}
