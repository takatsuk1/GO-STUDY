package email

import (
	"gin-mall/config"
	"gopkg.in/mail.v2"
)

type EmailSender struct {
	SmtpHost      string `json:"smtp_host"`
	SmtpEmailFrom string `json:"smtp_email_from"`
	SmtpPass      string `json:"smti_pass"`
}

func NewEmailSender() *EmailSender {
	eConfig := config.Config.Email
	return &EmailSender{
		SmtpHost:      eConfig.SmtpHost,
		SmtpEmailFrom: eConfig.SmtpEmail,
		SmtpPass:      eConfig.SmtpPass,
	}
}

func (s *EmailSender) Send(data, emailTo, subject string) error {
	m := mail.NewMessage()
	//发件人
	m.SetHeader("From", s.SmtpEmailFrom)
	//收件人
	m.SetHeader("To", emailTo)
	//主题
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)
	d := mail.NewDialer(s.SmtpHost, 465, s.SmtpEmailFrom, s.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
