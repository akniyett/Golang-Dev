package models

type Clothing struct {
	ID           int    `json:"id" db: "id"`
	Name         string `json:"name" db: "namr"`
	Description  string `json:"description" db: "description"`
	Size   string `json:"size" db: "size"`
	Price  float64 `json:"price" db: "price"`
	IsAvailable bool `json:"isAvailable" db: "isAvailable"`
}

type Accessory struct {
	ID           int    `json:"id" db: "id"`
	Name         string `json:"name" db: "name"`
	Description  string `json:"description" db: "description"`
	Size   string `json:"size" db: "size"`
	Price  float64 `json:"price" db: "price"`
	IsAvailable bool `json:"isAvailable" db: "isAvailable"`
	Manufacturer string `json:"manufacturer" db: "manufacturer"`
	Material string `json:"material" db: "material"`
}

type User struct {
	ID           int    `json:"id" db: "id"`
	Nick         string `json:"nick" db: "nick"`
	Password string `json:"password" db: "password"`
	Bio  string `json:"bio" db: "bio"`
	Email string `json:"email" db: "email"`
}
