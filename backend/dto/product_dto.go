package dto

type CreateProductRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,min=2,max=150"`
	Desc       string `json:"desc" binding:"max=1000"`
	Price      int64  `json:"price" binding:"required,gt=0"`
	Stock      int    `json:"stock" binding:"gte=0"`
}

type UpdateProductRequest struct {
	CategoryID *uint   `json:"category_id" binding:"omitempty"`
	Name       *string `json:"name" binding:"omitempty,min=2,max=150"`
	Desc       *string `json:"desc" binding:"omitempty,max=1000"`
	Price      *int64  `json:"price" binding:"omitempty,gt=0"`
	Stock      *int    `json:"stock" binding:"omitempty,gte=0"`
}

// ListProductQuery untuk menangkap parameter URL (Pagination & Search)
type ListProductQuery struct {
	Page     int    `form:"page,default=1" binding:"gte=1"`
	Limit    int    `form:"limit,default=10" binding:"gte=1,lte=100"`
	Q        string `form:"q"`        // Search nama produk
	Category string `form:"category"` // Filter by nama kategori atau ID
	InStock  bool   `form:"in_stock"` // Filter stok > 0
}

// MetaPagination untuk response standar Opsi 1
type MetaPagination struct {
	TotalData   int64 `json:"total_data"`
	TotalPage   int   `json:"total_page"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}
