package register

type UserRequest struct {
	ID                   *string `json:"id"`
	Username             string  `json:"username" valid:"required,stringlength(5|50),matches(^[a-zA-Z0-9_]+$)"`
	Name                 string  `json:"name" valid:"required,stringlength(0|255)"`
	LastName             string  `json:"last_name" valid:"required,stringlength(0|255)"`
	EmailNotifications   string  `json:"email_notifications" valid:"required,email,stringlength(5|255)"`
	IdentificationNumber string  `json:"identification_number" valid:"required,stringlength(0|255)"`
	IdentificationType   string  `json:"identification_type"`
	Password             string  `json:"password,omitempty"`
	PasswordConfirm      string  `json:"password_confirm,omitempty"`
}
