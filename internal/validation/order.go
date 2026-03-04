package validation

type CreateOrderInput struct {
	Items []struct {
		ProductID string `json:"product_id" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required,min=1"`
	} `json:"items" validate:"required,dive"`
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required,oneof=pending shipped delivered cancelled"`
}