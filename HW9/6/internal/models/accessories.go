package models

type (
	Accessory struct {
		ID           int    `json:"id" db: "id"`
		Name         string `json:"name" db: "name"`
		Description  string `json:"description" db: "description"`
		Size   string `json:"size" db: "size"`
		Price  float64 `json:"price" db: "price"`
		IsAvailable bool `json:"isAvailable" db: "isAvailable"`
		Manufacturer string `json:"manufacturer" db: "manufacturer"`
		Material string `json:"material" db: "material"`
	}
	AccessoriesFilter struct {
		Query *string `json:"query"`
	}
)
