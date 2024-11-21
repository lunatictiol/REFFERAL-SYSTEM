package types

type User struct {
	ID       string `json:"id"`       // UUID as a string
	Name     string `json:"name"`     // Text field
	Email    string `json:"email"`    // Unique text field
	Password string `json:"password"` // Text field
	Points   int    `json:"points"`   // Integer with a default value
}
