package models

type User struct {
	ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type Register struct {
	// ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email" example:"foo_bar@mail.com"`
	Password string `db:"password" json:"password" binding:"required,min=8,containsany=!@#$%^&*,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type Login struct {
	Email    string `db:"email" json:"email" binding:"required" example:"foo_bar@mail.com"`
	Password string `db:"password" json:"password" binding:"required"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Result  string `json:"result,omitempty" example:"register succesfully w/ ID: 37"`
	Error   string `json:"message,omitempty" example:"server unable to bind request"`
}

type LoginResponse struct {
	Result  string `json:"message" example:"logged in as UID 123"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
}

type PasswordBody struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `db:"password" json:"new_password" binding:"required,min=8,containsany=!@#$%^&*,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type EditPasswordResponse struct {
	Result  string `json:"result,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type LogoutResponse struct {
	Result  string `json:"result,omitempty" example:"logout succesfully"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" example:"server error while logout"`
}
