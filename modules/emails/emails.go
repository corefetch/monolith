package emails

import (
	"bytes"
	"embed"
	"html/template"
	"net/smtp"

	"learnt.io/core"
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

func AfterRegisterStudent(to Recipient, params TemplateParams) (err error) {
	return Send("Registration Complete", to, RegisterStudent, params)
}

func AfterRegisterTutor(to Recipient, params TemplateParams) (err error) {
	return Send("Registration Complete", to, RegisterStudent, params)
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
		"smtp-relay.brevo.com:587",
		smtp.PlainAuth(
			"",
			"cosmin.albulescu@gmail.com",
			"5IHpz6rBwn0ZYqGF",
			"smtp-relay.brevo.com",
		),
		"cosmin.albulescu@gmail.com",
		[]string{to.RecipientEmail()},
		job.Bytes(),
	)
}
