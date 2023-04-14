package user

type CreateUserDto struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	FullName        string `json:"full_name" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required,e164"`
}

type GetUsersDto struct {
	Id              string  `json:"id,omitempty"`
	Email           string  `json:"email"`
	Password        string  `json:"-"`
	FullName        string  `json:"full_name"`
	TelephoneNumber string  `json:"telephone_number"`
	IsPremium       bool    `json:"is_premium"`
	City            *string `json:"city"`
	Rating          float32 `json:"rating"`
}

type UpdateUserDto struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	FullName        string `json:"full_name" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required,e164"`
	City            string `json:"city"`
}

type RateDto struct {
	Rating  float32 `json:"rating" binding:"required"`
	Comment string  `json:"comment"`
}

type GetAllRatingsDto struct {
	Rating  float32 `json:"rating"`
	Comment string  `json:"comment"`
	User    string  `json:"user"`
	RatedBy string  `json:"rated_by"`
}
