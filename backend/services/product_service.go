package services

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"backend/dto"
	"backend/models"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("produk tidak ditemukan")
	ErrProductConflict = errors.New("produk dengan nama ini sudah ada")
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	var p models.Product
	// Gunakan Preload("Category") agar data kategori ikut terbawa
	if err := s.DB.Preload("Category").First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

// List menangani Pagination, Search, dan Filter Kategori
func (s *ProductService) List(query dto.ListProductQuery) ([]models.Product, dto.MetaPagination, error) {
	var products []models.Product
	var totalData int64

	db := s.DB.Model(&models.Product{}).Preload("Category")

	// 1. Search berdasarkan nama produk
	if q := strings.TrimSpace(query.Q); q != "" {
		db = db.Where("name ILIKE ?", "%"+q+"%")
	}

	// 2. Filter berdasarkan nama Kategori (bisa juga ID jika diubah)
	if cat := strings.TrimSpace(query.Category); cat != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name ILIKE ?", "%"+cat+"%")
	}

	// 3. Filter stok
	if query.InStock {
		db = db.Where("products.stock > 0")
	}

	// 4. Hitung Total Data (sebelum di-limit)
	if err := db.Count(&totalData).Error; err != nil {
		return nil, dto.MetaPagination{}, err
	}

	// 5. Pagination (Limit & Offset)
	offset := (query.Page - 1) * query.Limit
	if err := db.Order("products.id ASC").Limit(query.Limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, dto.MetaPagination{}, err
	}

	if products == nil {
		products = []models.Product{}
	}

	// 6. Susun Metadata
	totalPage := int(math.Ceil(float64(totalData) / float64(query.Limit)))
	meta := dto.MetaPagination{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: query.Page,
		Limit:       query.Limit,
	}

	return products, meta, nil
}

func (s *ProductService) Create(req dto.CreateProductRequest) (*models.Product, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("nama produk wajib diisi")
	}

	// Validasi CategoryID
	if _, err := NewCategoryService(s.DB).GetByID(req.CategoryID); err != nil {
		return nil, fmt.Errorf("kategori tidak valid: %w", err)
	}

	// Cek duplikat nama
	var existing models.Product
	err := s.DB.Where("LOWER(name) = ?", strings.ToLower(req.Name)).First(&existing).Error
	if err == nil {
		return nil, ErrProductConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: strings.TrimSpace(req.Desc),
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.DB.Create(&product).Error; err != nil {
		if isDuplicateErr(err) {
			return nil, ErrProductConflict
		}
		return nil, err
	}

	// Preload kategori untuk response
	s.DB.Preload("Category").First(&product, product.ID)
	return &product, nil
}

func (s *ProductService) Update(id uint, req dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != nil {
		if _, err := NewCategoryService(s.DB).GetByID(*req.CategoryID); err != nil {
			return nil, fmt.Errorf("kategori tidak valid: %w", err)
		}
		product.CategoryID = *req.CategoryID
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, fmt.Errorf("nama produk tidak boleh kosong")
		}
		var dup models.Product
		err := s.DB.Where("LOWER(name) = ? AND id <> ?", strings.ToLower(name), product.ID).First(&dup).Error
		if err == nil {
			return nil, ErrProductConflict
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		product.Name = name
	}
	if req.Desc != nil {
		product.Description = strings.TrimSpace(*req.Desc)
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.DB.Save(product).Error; err != nil {
		if isDuplicateErr(err) {
			return nil, ErrProductConflict
		}
		return nil, err
	}

	// Preload kategori untuk response
	s.DB.Preload("Category").First(product, product.ID)
	return product, nil
}

func isDuplicateErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
