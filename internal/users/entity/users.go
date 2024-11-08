package entity

// req
type UserRegisterRequest struct {
	ID           int      `json:"id"`
	PhoneNumber  string   `json:"phone_number" validate:"required,e164"`
	Name         string   `json:"name" validate:"required,min=2,max=50"`
	Password     string   `json:"password" validate:"required,min=8"`
	Role         string   `json:"role" validate:"required,oneof=tukang store"`
	CreatedAt    int64    `json:"created_at"`
	StoreName    *string  `json:"store_name,omitempty" validate:"omitempty,min=2,max=255"`
	Address      *string  `json:"address,omitempty" validate:"omitempty,min=2,max=500"`
	Latitude     *float64 `json:"latitude,omitempty" validate:"omitempty"`
	Longitude    *float64 `json:"longitude,omitempty" validate:"omitempty"`
	WorkingHours *string  `json:"working_hours,omitempty" validate:"omitempty,min=5,max=100"`
	IsVerified   bool     `json:"is_verified"`
}

// -7.968437, 112.596530

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Password    string `json:"password" validate:"required,min=8"`
}

// res
type UserJWT struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

type StoreJWT struct {
	ID          int     `json:"id"`
	PhoneNumber string  `json:"phone_number"`
	Name        string  `json:"name"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	StoreName   string  `json:"store_name"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type UserDetailResponse struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}

type UserListSubResponse struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type UserListResponse struct {
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
	Users []UserListSubResponse `json:"users"`
}
