package models

type Verification struct {
	Email      string `json:"email" binding:"required"`
	Name       string `json:"name"`
	Password   string `json:"password"`
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
type VerifyOtp struct {
	Email string `json:"email" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	CreatedAT int64  `json:"created_at"`
	UpdatedAT int64  `json:"updated_at"`
}
