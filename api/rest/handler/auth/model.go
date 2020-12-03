package auth

type Model struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID int    `json:"client_id"`
	HostName string `json:"host_name"`
	RealIP   string `json:"real_ip"`
}

type ModelResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
