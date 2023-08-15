package register

type UserRequest struct {
	ID              *string `json:"id"`
	Username        string  `json:"username"`
	CodeStudent     string  `json:"code_student"`
	Dni             string  `json:"dni"`
	Names           string  `json:"names"`
	LastnameFather  string  `json:"lastname_father"`
	LastnameMother  string  `json:"lastname_mother"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	PasswordConfirm string  `json:"password_confirm,omitempty"`
}
