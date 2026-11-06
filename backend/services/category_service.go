package services

import (
	"errors"
	"fmt"
	"strings"

	"backend/dto"
	"backend/models"

	"gorm.io/gorm"
)

// Error sentinel khusus Category
var (
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
	ErrCategoryConflict = errors.New("kategori dengan nama ini sudah ada")
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	var c models.Category
	if err := s.DB.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

// List mengambil semua kategori (biasanya tidak perlu pagination karena jumlahnya sedikit)
func (s *CategoryService) List() ([]models.Category, error) {
	var categories []models.Category
	if err := s.DB.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []models.Category{}
	}
	return categories, nil
}

func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*models.Category, error) {
	req.Name = strings.TrimSpace(req.Name)

	// Cek duplikat (case-insensitive)
	var existing models.Category
	err := s.DB.Where("LOWER(name) = ?", strings.ToLower(req.Name)).First(&existing).Error
	if err == nil {
		return nil, ErrCategoryConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := models.Category{Name: req.Name}
	if err := s.DB.Create(&category).Error; err != nil {
		if isDuplicateErr(err) {
			return nil, ErrCategoryConflict
		}
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Update(id uint, req dto.UpdateCategoryRequest) (*models.Category, error) {
	c, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, fmt.Errorf("nama kategori tidak boleh kosong")
		}

		// Cek duplikat (kecuali dirinya sendiri)
		var dup models.Category
		err := s.DB.Where("LOWER(name) = ? AND id <> ?", strings.ToLower(name), c.ID).First(&dup).Error
		if err == nil {
			return nil, ErrCategoryConflict
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		c.Name = name
	}

	if err := s.DB.Save(c).Error; err != nil {
		if isDuplicateErr(err) {
			return nil, ErrCategoryConflict
		}
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Delete(id uint) error {
	c, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// CEK: Apakah ada produk yang menggunakan kategori ini?
	var productCount int64
	if err := s.DB.Model(&models.Product{}).Where("category_id = ?", c.ID).Count(&productCount).Error; err != nil {
		return err
	}
	if productCount > 0 {
		return fmt.Errorf("tidak dapat menghapus kategori: masih ada %d produk yang menggunakan kategori ini", productCount)
	}

	return s.DB.Delete(c).Error
}
