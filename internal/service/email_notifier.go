package service

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"text/template"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	from       string
	password   string
	host       string
	port       int
	subscriber message.Subscriber
	logger     watermill.LoggerAdapter
}

func NewEmailNotifier(subscriber message.Subscriber, logger watermill.LoggerAdapter) (*EmailNotifier, error) {
	return &EmailNotifier{
		from:       "aaaaaaaaaaaaaaaa@gmail.com", // почта отправителя
		password:   "aaaa aaaa aaaa aaaa",        // App Password
		host:       "smtp.gmail.com",             // SMTP-сервер
		port:       587,                          // SMTP-порт
		subscriber: subscriber,
		logger: logger,
	}, nil
}

const templatePath = "../../resources/template/email_template.html"
const dataFormat = "2006-01-02 15:04:05"

func (e *EmailNotifier) Notify(_ context.Context, r scripts.Result, email scripts.Email) error {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	values := r.Out().Get()

	var errorMsg string
	if r.ErrorMessage() != nil {
		errorMsg = *r.ErrorMessage()
	}

	data := struct {
		Command      string
		StartedAt    string
		Code         int
		OutputValues []string
		ErrorMessage string
	}{
		Command:      r.Job().Command(),
		StartedAt:    r.Job().StartedAt().Format(dataFormat),
		Code:         int(r.Code()),
		OutputValues: values,
		ErrorMessage: errorMsg,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Результат выполнения")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(e.host, e.port, e.from, e.password)
	return d.DialAndSend(m)
}

func (e *EmailNotifier) Listen(ctx context.Context) {
	messages, err := e.subscriber.Subscribe(ctx, "script-finished")
	if err != nil {
		e.logger.Error("Subscribe error", err, nil)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				e.logger.Error("panic recovered in email notifier", nil, watermill.LogFields{"recover": r})
			}
		}()

		for {
			select {
			case <-ctx.Done():
				e.logger.Info("EmailNotifier stopped due to context cancel", nil)
				return
			case msg, ok := <-messages:
				if !ok {
					e.logger.Info("EmailNotifier channel closed", nil)
					return
				}

				var event ScriptFinishedEvent
				if err := json.Unmarshal(msg.Payload, &event); err != nil {
					e.logger.Error("Unmarshal error", err, nil)
					msg.Nack()
					continue
				}

				if err := e.Notify(ctx, event.result, event.userEmail); err != nil {
					e.logger.Error("Notify error", err, nil)
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}
	}()
}
