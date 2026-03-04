package validation

type ProductInput struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description" binding:"required,min=10"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name" binding:"omitempty,min=3"`
	Description *string  `json:"description" binding:"omitempty,min=10"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
	Category    *string  `json:"category" binding:"omitempty"`
}
