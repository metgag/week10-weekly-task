package models

type FulfilledResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Result  interface{} `json:"result"`
}

func NewFullfilledResponse(statusCode int, result interface{}) FulfilledResponse {
	return FulfilledResponse{
		Success: true,
		Status:  statusCode,
		Result:  result,
	}
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}
