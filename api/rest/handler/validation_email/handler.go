package validation_email

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/ciphers"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/env"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/msgs"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/password"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/response"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/sendmail"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/pkg/auth/login"
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
		logger.Error.Printf(h.TxID, "couldn't bind model validate email: %v", err)
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
		logger.Error.Printf(h.TxID, "couldn't create verify code: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	parameters["@access-code"] = emailCode
	parameters["@TEMPLATE-PATH"] = e.Template.EmailCode
	tos := []string{m.Email}

	email := sendmail.Model{
		From:        e.Template.EmailSender,
		To:          tos,
		CC:          nil,
		Subject:     e.Template.EmailCodeSubject,
		Body:        fmt.Sprintf("<h1>%s</h1>", emailCode),
		Attach:      "",
		Attachments: nil,
	}
	tpl, err := email.GenerateTemplateMail(parameters)
	if err != nil {
		logger.Error.Println(h.TxID, "error when parse template: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	email.Body = tpl
	err = email.SendMail()
	if err != nil {
		logger.Error.Println(h.TxID, "error when execute send email: %v", err)
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
		logger.Error.Printf(h.TxID, "couldn't bind model verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	dataVerify, code, err := srvUser.SrvVerificationEmail.GetVerificationEmailByID(m.Id)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify == nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !password.Compare(dataVerify.Email, dataVerify.VerificationCode, m.Code) {
		logger.Error.Printf(h.TxID, "the verification code is not correct: %v", err)
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
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "successful email validation"
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerValidationEmail) GenerateOtp(c *fiber.Ctx) error {
	var parameters = make(map[string]string, 0)
	e := env.NewConfiguration()
	var msg msgs.Model
	res := response.Model{Error: true}
	m := ReqGenerateOtp{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't bind model validate email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)
	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	emailCode := strconv.Itoa(rand.Intn(max-min+1) + min)
	verifiedCode := password.Encrypt(emailCode)

	user, code, err := srvUser.SrvUsers.GetUserByIdentificationNumber(m.IdentificationNumber)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get user by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	codVerify, code, err := srvUser.SrvVerificationEmail.CreateVerificationEmail(m.Email, verifiedCode, user.IdentificationNumber, nil)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't create verify code: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	jwt, code, err := login.GenerateJWTOtp(ciphers.Encrypt(emailCode), codVerify.ID)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't generate otp jwt: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	parameters["@access-code"] = jwt
	parameters["@TEMPLATE-PATH"] = e.Template.EmailCode
	tos := []string{m.Email}

	email := sendmail.Model{
		From:        e.Template.EmailSender,
		To:          tos,
		CC:          nil,
		Subject:     e.Template.EmailCodeSubject,
		Attach:      "",
		Attachments: nil,
	}
	tpl, err := email.GenerateTemplateMail(parameters)
	if err != nil {
		logger.Error.Println(h.TxID, "error when parse template: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	email.Body = tpl
	err = email.SendMail()
	if err != nil {
		logger.Error.Println(h.TxID, "error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = codVerify.ID
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerValidationEmail) verifyOtp(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model

	m := VerificationDataRequest{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't bind model verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	dataVerify, code, err := srvUser.SrvVerificationEmail.GetVerificationEmailByID(m.Id)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify == nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !password.Compare(dataVerify.Email, dataVerify.VerificationCode, ciphers.Decrypt(m.Code)) {
		logger.Error.Printf(h.TxID, "the verification code is not correct: %v", err)
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
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "successful email validation"
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
