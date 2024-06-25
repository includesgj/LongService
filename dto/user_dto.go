package dto

import "GinProject12/model"

// UserDto 暴露的用户信息
type UserDto struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func ToAdminDto(user model.Admin) *UserDto {
	return &UserDto{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

func ToUserDto(user model.User) *UserDto {
	return &UserDto{
		Username: user.Username,
		Email:    user.Email,
		Role:     "普通用户",
	}
}
