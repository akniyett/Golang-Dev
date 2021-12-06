package models

type (
	User struct {
		ID           int    `json:"id" db: "id"`
		Nick         string `json:"nick" db: "nick"`
		Password string `json:"password" db: "password"`
		Bio  string `json:"bio" db: "bio"`
		Email string `json:"email" db: "email"`
	}
	UsersFilter struct {
		Query *string `json:"query"`
	}
)