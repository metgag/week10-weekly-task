package models

type User struct {
	ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type Register struct {
	ID       uint16 `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Password string `db:"password" json:"password" binding:"min=8,containsany=!@#$%^&*,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type Login struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type AuthResponse struct {
	Success bool
	Result  string
	Error   string
}
