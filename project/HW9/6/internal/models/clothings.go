package models

type (
	Clothing struct {
		ID           int    `json:"id" db:"id"`
		Name         string `json:"name" db:"name"`
		Description  string `json:"description" db:"description"`
		Size   string `json:"size" db:"size"`
		Price  float64 `json:"price" db:"price"`
		IsAvailable bool `json:"isAvailable" db:"isAvailable"`
	}
	ClothingsFilter struct {
		Query *string `json:"query"`
	}
)