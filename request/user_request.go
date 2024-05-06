package request

type RegisterRequest struct {
	Name                 string `json:"name" form:"name"`
	Email                string `json:"email" form:"email"`
	Password             string `json:"password" form:"password"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
