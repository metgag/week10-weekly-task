package models

type UserInf struct {
	UID         uint16  `db:"user_id" json:"user_id"`
	FirstName   string  `db:"first_name" json:"first_name"`
	LastName    string  `db:"last_name" json:"last_name"`
	PhoneNumber string  `db:"phone_number" json:"phone_number" binding:"min=10.numeric" example:"08224422765"`
	PointCount  float32 `db:"point_count" json:"point_count" example:"4.2"`
}

type UserinfResponse struct {
	Result  UserInf `json:"result"`
	Success bool    `json:"success"`
	Error   string  `json:"error"`
}

type NewInf struct {
	FirstName   string  `db:"first_name" json:"first_name"`
	LastName    string  `db:"last_name" json:"last_name"`
	PhoneNumber string  `db:"phone_number" json:"phone_number" binding:"min=10,numeric" example:"08667728761"`
	PointCount  float32 `db:"point_count" json:"point_count" example:"4.8"`
}

type UpdateResponse struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type UserOrder struct {
	UID          uint16         `db:"user_id" json:"user_id"`
	OrderHistory []OrderHistory `db:"order_history" json:"order_history"`
}

type HistoryResponse struct {
	Result  UserOrder
	Success bool
	Error   string
}
