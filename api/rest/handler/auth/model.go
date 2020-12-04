package auth

type LoginRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID int    `json:"client_id"`
	HostName string `json:"host_name"`
	RealIP   string `json:"real_ip"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type PasswordPolicyRequest struct {
	Password string `json:"password"`
}
