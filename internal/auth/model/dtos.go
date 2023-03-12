package auth

type RegistrationDto struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	FullName        string `json:"full_name" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required,e164"`
}

type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
