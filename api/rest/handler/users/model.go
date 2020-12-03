package register

type Model struct {
	ID                   *string `json:"id"`
	Username             string  `json:"username"`
	Name                 string  `json:"name"`
	LastName             string  `json:"last_name"`
	EmailNotifications   string  `json:"email_notifications"`
	IdentificationNumber string  `json:"identification_number"`
	IdentificationType   string  `json:"identification_type"`
	Password             string  `json:"password"`
	PasswordConfirm      string  `json:"password_confirm"`
}

type ModelResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
