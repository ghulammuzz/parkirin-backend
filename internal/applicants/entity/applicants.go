package entity

import userEntity "github.com/ghulammuzz/backend-parkerin/internal/users/entity"

type Application struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	StoreID   int    `json:"store_id"`
	Status    string `json:"status"`
	AppliedAt int64  `json:"applied_at"`
}

type ApplicationResponse struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Status   string `json:"status"`
}

type ApplicationUserResponse struct {
	ID             int    `json:"id"`
	StoreID        int    `json:"store_id"`
	StoreName      string `json:"store_name"`
	IsHiring       bool   `json:"is_hiring"`
	IsDirectHiring bool   `json:"is_direct_hiring"`
	Address        string `json:"address"`
	UrlImage       string `json:"url_image"`
	WorkingHours   string `json:"working_hours"`
	Status         string `json:"status"`
}

type ApplicationUserResponseDetail struct {
	ID        int                           `json:"id"`
	StoreName string                        `json:"store_name"`
	Address   string                        `json:"address"`
	Status    string                        `json:"status"`
	User      userEntity.UserDetailResponse `json:"user"`
}
