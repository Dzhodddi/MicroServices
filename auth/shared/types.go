package shared

type UserWithToken struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Email    string `json:"email" validate:"required,min=4,max=32,email"`
	Token    string `json:"token" validate:"required,min=8,max=32"`
}

type RegisterNewUser struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Email    string `json:"email" validate:"required,min=4,max=32,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,min=4,max=32,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type User struct {
	Username  string `json:"username" validate:"required,min=4,max=32"`
	Email     string `json:"email" validate:"required,min=4,max=32,email"`
	CreatedAt string `json:"createdAt"`
	IsActive  bool   `json:"isActive"`
	RoleID    int64  `json:"role_id"`
}
