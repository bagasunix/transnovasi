package responses

import "time"

type UserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	RoleID   int    `json:"role_id"`
	IsActive int16  `json:"is_active"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
