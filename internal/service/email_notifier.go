package service

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	templateFile string
	from         string
	password     string
	host         string
	port         int
}

func NewEmailNotifier(templateFile, from, password, host string, port int) (*EmailNotifier, error) {
	return &EmailNotifier{templateFile, from, password, host, port}, nil
}

type EmailTemplateData struct {
	StartedAt    string
	Code         int
	OutputValues []string
	ErrorMessage string
}

func (e *EmailNotifier) Notify(_ context.Context, j *scripts.Job, email scripts.Email) error {
	res, err := j.Result()
	if err != nil {
		return err
	}

	values := res.Output()
	rawValues := make([]string, len(values))
	for i, v := range values {
		rawValues[i] = v.String()
	}

	data := EmailTemplateData{
		StartedAt:    j.CreatedAt().Format("02.01.2006 15:04:05"),
		Code:         int(res.Code()),
		OutputValues: rawValues,
	}

	if msg := res.ErrorMessage(); msg != nil {
		data.ErrorMessage = *msg
	}

	tmplBytes, err := os.ReadFile(e.templateFile)
	if err != nil {
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.from)
	msg.SetHeader("To", string(email))
	msg.SetHeader("Subject", "Результат выполнения")
	msg.SetBody("text/html", body.String())

	dialer := gomail.NewDialer(e.host, e.port, e.from, e.password)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
