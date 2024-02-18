package emails

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"corefetch/core"
)

//go:embed templates/**
var templates embed.FS

type TemplateKey string
type TemplateParams = core.M

const (
	RegisterStudent TemplateKey = "register_student.html"
	RegisterTutor   TemplateKey = "register_tutor.html"
)

type Recipient interface {
	RecipientName() string
	RecipientEmail() string
}

type SendingJob struct {
	bytes.Buffer
	To Recipient
}

func Send(subject string, to Recipient, templateKey TemplateKey, params TemplateParams) (err error) {

	data, err := templates.ReadFile("templates/" + string(templateKey))

	if err != nil {
		return err
	}

	t, err := template.New("root").Parse(string(data))

	if err != nil {
		return err
	}

	job := SendingJob{To: to}
	job.WriteString("To: " + to.RecipientEmail() + "\r\n")
	job.WriteString("Subject: " + subject + "\r\n")
	job.WriteString("\r\n")

	if err := t.Execute(&job, params); err != nil {
		return err
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT")),
		smtp.PlainAuth(
			"",
			os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"),
			os.Getenv("SMTP_HOST"),
		),
		os.Getenv("SMTP_USER"),
		[]string{to.RecipientEmail()},
		job.Bytes(),
	)
}
