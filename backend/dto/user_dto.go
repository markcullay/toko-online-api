package dto

type UpdateUserRequest struct {
	Name     *string `json:"name"     binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email"    binding:"omitempty,email"`
	Role     *string `json:"role"     binding:"omitempty,oneof=admin customer"`
	Password *string `json:"password" binding:"omitempty,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
