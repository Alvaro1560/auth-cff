package validation_email

type VerificationRequest struct {
	Email string `json:"email"`
}
type VerificationDataRequest struct {
	Id             int64  `json:"id"`
	Identification string `json:"identification"`
	Code           string `json:"code"`
}

type ReqGenerateOtp struct {
	IdentificationNumber string `json:"identification_number"`
	Email                string `json:"email"`
}
