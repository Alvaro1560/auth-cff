package sendmail

import (
	"bytes"
	"crypto/tls"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"
	"html/template"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"

	"gopkg.in/gomail.v2"
)

func (e *Model) SendMail() error {

	smtpPort, smtpHost, smtpEmail, smtpPassword := e.getParams()
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", e.To...)
	m.SetHeader("Cc", e.CC...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)
	if len(e.Attach) > 0 {
		//m.Attach(e.Attach)
	}

	for _, v := range e.Attachments {
		m.Attach(v)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		logger.Error.Printf("couldn't emil to: %s, subject: %s, %v", e.To, e.Subject, err)
		return err
	}

	return nil
}

func (e *Model) AddAttach(fn string) {
	if len(e.Attachments) == 0 {
		e.Attachments = make([]string, 0)
	}

	e.Attachments = append(e.Attachments, fn)
}

func (e *Model) getParams() (int, string, string, string) {
	c := env.NewConfiguration()
	smtpPort := c.Smtp.Port
	smtpHost := c.Smtp.Host
	smtpEmail := c.Smtp.Email
	smtpPassword := c.Smtp.Password
	return smtpPort, smtpHost, smtpEmail, smtpPassword
}

func (e *Model) GenerateTemplateMail(param map[string]string) (string, error) {
	bf := &bytes.Buffer{}
	var tpl *template.Template
	path := param["@TEMPLATE-PATH"]
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
	err := tpl.ExecuteTemplate(bf, path, &param)
	if err != nil {
		logger.Error.Printf("couldn't generate template body mail: %v", err)
		return "", err
	}
	return bf.String(), nil
}
