package entity

import userEntity "github.com/ghulammuzz/backend-parkerin/internal/users/entity"

type ListStoreSubResponse struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	StoreName    string `json:"store_name"`
	Address      string `json:"address"`
	UrlImage     string `json:"url_image"`
	WorkingHours string `json:"working_hours"`
	IsHiring     bool   `json:"is_hiring"`
	IsPaid       bool   `json:"is_paid"`
}

type ListStoreResponse struct {
	Page   int                    `json:"page"`
	Limit  int                    `json:"limit"`
	Stores []ListStoreSubResponse `json:"stores"`
}

type DashboardStoreResponse struct {
	ID           int                           `json:"id"`
	User         userEntity.UserDetailResponse `json:"user"`
	StoreName    string                        `json:"store_name"`
	Address      string                        `json:"address"`
	UrlImage     string                        `json:"url_image"`
	Latitude     float64                       `json:"latitude"`
	Longitude    float64                       `json:"longitude"`
	WorkingHours string                        `json:"working_hours"`
	IsHiring     bool                          `json:"is_hiring"`
	IsPaid       bool                          `json:"is_paid"`
	CreatedAt    int64                         `json:"created_at"`
	IsVerified   bool                          `json:"is_verified"`
}

type DetailStoreResponse struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	StoreName    string  `json:"store_name"`
	Address      string  `json:"address"`
	UrlImage     string  `json:"url_image"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	WorkingHours string  `json:"working_hours"`
	IsPaid       bool    `json:"is_paid"`
	IsVerified   bool    `json:"is_verified"`
	IsHiring     bool    `json:"is_hiring"`
	CreatedAt    int64   `json:"created_at"`
}

type UpdateIsHiringRequest struct {
	IsHiring bool `json:"is_hiring"`
}
