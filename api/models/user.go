package models

type Verification struct {
	Email      string `json:"email" binding:"required"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Otp        string `json:"otp"`
	ISverified bool   `json:"isverified"`
	CreatedAT  int64  `json:"created_at"`
	UpdatedAT  int64  `json:"Updated_at"`
}

type UserClient struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
