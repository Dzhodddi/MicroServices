package shared

type User struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Email    string `json:"email" validate:"required,min=4,max=32,email"`
	Token    string `json:"token" validate:"required"`
}

type RegisterNewUser struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Email    string `json:"email" validate:"required,min=4,max=32,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
