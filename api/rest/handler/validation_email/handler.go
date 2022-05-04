package validation_email

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/sendmail"

	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type handlerValidationEmail struct {
	DB   *sqlx.DB
	TxID string
}

func (h *handlerValidationEmail) sendCode(c *fiber.Ctx) error {
	var parameters = make(map[string]string, 0)
	e := env.NewConfiguration()
	var msg msgs.Model
	res := response.Model{Error: true}
	m := VerificationRequest{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model validate email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)
	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	emailCode := strconv.Itoa(rand.Intn(max-min+1) + min)
	verifiedCode := password.Encrypt(emailCode)

	codVerify, code, err := srvUser.SrvVerificationEmail.CreateVerificationEmail(m.Email, verifiedCode, "", nil)
	if err != nil {
		logger.Error.Printf("couldn't create verify code: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	parameters["@access-code"] = emailCode
	parameters["@TEMPLATE-PATH"] = e.Template.EmailCode
	tos := []string{m.Email}

	email := sendmail.Model{
		From:        "no-reply@e-capture.co",
		To:          tos,
		CC:          nil,
		Subject:     e.Template.EmailCodeSubject,
		Body:        fmt.Sprintf("<h1>%s</h1>", emailCode),
		Attach:      "",
		Attachments: nil,
	}
	tpl, err := email.GenerateTemplateMail(parameters)
	if err != nil {
		logger.Error.Println("error when parse template: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	email.Body = tpl
	err = email.SendMail()
	if err != nil {
		logger.Error.Println("error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = codVerify.ID
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerValidationEmail) verifyCode(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model

	m := VerificationDataRequest{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	dataVerify, code, err := srvUser.SrvVerificationEmail.GetVerificationEmailByID(m.Id)
	if err != nil {
		logger.Error.Printf("couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify.ID == 0 {
		logger.Error.Printf("couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if !password.Compare(dataVerify.Email, dataVerify.VerificationCode, m.Code) {
		logger.Error.Printf("the verification code is not correct: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify.VerificationDate != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	dateTime := time.Now()
	_, code, err = srvUser.SrvVerificationEmail.UpdateVerificationEmail(dataVerify.ID, dataVerify.Email, "", dataVerify.Identification, &dateTime)
	if err != nil {
		logger.Error.Printf("couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "successful email validation"
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
