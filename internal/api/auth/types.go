package auth

type CreateUserDTO struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,password" json:"password"`
	Role     string `validate:"required,oneof=student admin" json:"role"`
}
