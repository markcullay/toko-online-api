package services

import (
	"errors"
	"fmt"
	"strings"

	"backend/dto"
	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

// Error sentinel — dipakai handler untuk memetakan status code.
var (
	ErrNotFound  = errors.New("tidak ditemukan")
	ErrForbidden = errors.New("akses ditolak")
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// GetByID mengambil user berdasarkan ID. Mengembalikan ErrNotFound jika tidak ada.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := s.DB.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// UpdateUser menerapkan aturan bisnis update user.
// actorID & actorRole adalah pelaku (dari JWT), bukan dari *gin.Context.
func (s *UserService) UpdateUser(actorID uint, actorRole string, targetID uint, req dto.UpdateUserRequest) (*models.User, error) {
	u, err := s.GetByID(targetID)
	if err != nil {
		return nil, err
	}

	isAdmin := actorRole == string(models.RoleAdmin)
	isSelf := actorID == u.ID

	// Hanya admin yang boleh mengubah role orang lain
	if req.Role != nil {
		if !isAdmin {
			return nil, fmt.Errorf("%w: hanya admin yang dapat mengubah peran", ErrForbidden)
		}
		newRole := strings.ToLower(strings.TrimSpace(*req.Role))
		if newRole != string(models.RoleAdmin) && newRole != string(models.RoleCustomer) {
			return nil, fmt.Errorf("peran tidak valid")
		}
		if isSelf && newRole != string(models.RoleAdmin) {
			return nil, fmt.Errorf("anda tidak dapat menurunkan peran anda sendiri")
		}
		u.Role = models.Role(newRole)
	}

	// Password hanya lewat endpoint khusus
	if req.Password != nil {
		return nil, fmt.Errorf("gunakan endpoint PUT /api/users/:id/password untuk mengganti kata sandi")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, fmt.Errorf("nama tidak boleh kosong")
		}
		u.Name = name
	}

	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		if email == "" {
			return nil, fmt.Errorf("email tidak boleh kosong")
		}
		var dup models.User
		err := s.DB.Where("LOWER(email) = ? AND id <> ?", email, u.ID).First(&dup).Error
		if err == nil {
			return nil, fmt.Errorf("email sudah dipakai pengguna lain")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		u.Email = email
	}

	if err := s.DB.Save(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// ChangePassword mengganti password milik sendiri.
func (s *UserService) ChangePassword(actorID, targetID uint, req dto.ChangePasswordRequest) error {
	if actorID != targetID {
		return fmt.Errorf("%w: anda hanya dapat mengganti kata sandi anda sendiri", ErrForbidden)
	}

	u, err := s.GetByID(targetID)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(u.Password, req.OldPassword) {
		return fmt.Errorf("kata sandi lama salah")
	}

	hash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	u.Password = hash

	if err := s.DB.Save(u).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser menghapus user dengan aturan:
// - tidak boleh hapus diri sendiri
// - tidak boleh hapus admin terakhir
func (s *UserService) DeleteUser(actorID, targetID uint) error {
	target, err := s.GetByID(targetID)
	if err != nil {
		return err
	}

	if target.ID == actorID {
		return fmt.Errorf("anda tidak dapat menghapus akun anda sendiri")
	}

	if target.Role == models.RoleAdmin {
		var adminCount int64
		if err := s.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount).Error; err != nil {
			return err
		}
		if adminCount <= 1 {
			return fmt.Errorf("tidak dapat menghapus admin terakhir")
		}
	}

	if err := s.DB.Delete(target).Error; err != nil {
		return err
	}
	return nil
}
