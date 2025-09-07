package models

type User struct {
	ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type Register struct {
	ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email" example:"foo_bar@mail.com"`
	Password string `db:"password" json:"password" binding:"required,min=8,containsany=!@#$%^&*,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type Login struct {
	Email    string `db:"email" json:"email" binding:"required" example:"foo_bar@mail.com"`
	Password string `db:"password" json:"password" binding:"required"`
}

type RegisterResponse struct {
	Success bool   `example:"true"`
	Result  string `example:"register succesfully w/ ID: 1"`
	Error   string `example:"server unable to bind request"`
}

type LoginResponse struct {
	Result  string `example:"logged in as UID 1"`
	Success bool   `example:"true"`
	Bearer  string
}
