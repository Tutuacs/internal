package auth

// TODO: Create types and dtos for auth

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Role     int    `json:"role"`
	Password string `json:"-" validate:"required,min=3,max=130"`
}

type RegisterUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}
