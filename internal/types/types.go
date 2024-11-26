package types

type User struct {
	ID       string `json:"id"`       // UUID as a string
	Name     string `json:"name"`     // Text field
	Email    string `json:"email"`    // Unique text field
	Password string `json:"password"` // Text field
	Points   int    `json:"points"`   // Integer with a default value
}

type RegisterUserPayload struct {
	Name        string `json:"name"`     // Text field
	Email       string `json:"email"`    // Unique text field
	Password    string `json:"password"` // Text field
	ReferalCode string `json:"referal_code"`
}

type ReferalData struct {
	Id          string
	ReferalCode string
	ReferedBy   string
	IsUsed      bool
}

type LoginUserPayload struct {
	Email    string `json:"email"`    // Unique text field
	Password string `json:"password"` // Text field
}

type ReferRequest struct {
	ServiceID string `json:"service_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}
