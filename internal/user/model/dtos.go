package user

type CreateUserDto struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FullName        string `json:"full_name"`
	TelephoneNumber string `json:"telephone_number"`
}

type GetUsersDto struct {
	Id              string  `json:"id,omitempty"`
	Email           string  `json:"email"`
	FullName        string  `json:"full_name"`
	TelephoneNumber string  `json:"telephone_number"`
	IsPremium       bool    `json:"is_premium"`
	City            *string `json:"city"`
	Rating          float32 `json:"rating"`
}

type UpdateUserDto struct {
	Id       string  `json:"id,omitempty"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName string  `json:"full_name"`
	City     *string `json:"city"`
}
