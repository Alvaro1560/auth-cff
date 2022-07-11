package sendmail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"text/template"
	"time"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"

	"gopkg.in/gomail.v2"
)

func (e *Model) SendMail() error {

	c := env.NewConfiguration()
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

	mp := c.Smtp.Port
	d := gomail.NewDialer(c.Smtp.Host, mp, c.Smtp.Email, c.Smtp.Password)
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

func (e *Model) SendMailNotification(template string, userID string, toMail string, subject string) {
	myMail := Model{}
	param := make(map[string]string)

	param["TEMPLATE-PATH"] = template
	param["user"] = userID
	param["FECHA-EXECUTE"] = time.Now().String()
	param["TO-MAIL"] = toMail
	param["FROM-EMAIL"] = "no-reply@e-capture.co"
	param["SUBJECT-EMAIL"] = subject

	body, err := e.generateTemplateMail(param)
	if err != nil {
		logger.Error.Printf("couldn't generate body in NotificationEmail: %v", err)
		return
	}

	email := param["TO-MAIL"]
	var tos = []string{email}

	myMail.From = param["FROM-EMAIL"]
	myMail.To = tos
	myMail.Subject = fmt.Sprintf(`%s`, param["SUBJECT-EMAIL"])
	myMail.Body = body

	err = myMail.SendMail()
	if err != nil {
		logger.Error.Printf("couldn't sendMail NotificationEmail: %v", err)
		return
	}

	return
}
func (e *Model) generateTemplateMail(param map[string]string) (string, error) {
	bf := &bytes.Buffer{}
	tpl := &template.Template{}

	tpl = template.Must(template.New("").ParseGlob("notifications/*.gohtml"))
	err := tpl.ExecuteTemplate(bf, param["TEMPLATE-PATH"], &param)
	if err != nil {
		logger.Error.Printf("couldn't generate template body mail: %v", err)
		return "", err
	}
	return bf.String(), nil
}
