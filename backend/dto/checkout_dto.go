package dto

type CheckoutItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Qty       int  `json:"qty"        binding:"required,gt=0"`
}

type CheckoutRequest struct {
	Buyer  string         `json:"buyer"  binding:"required,min=2,max=100"`
	Method string         `json:"method" binding:"required"`
	Items  []CheckoutItem `json:"items"  binding:"required,min=1,dive"`
}
