package responses

import "time"

type CustomerResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	Address  string
	RoleID   int   `json:"role_id"`
	IsActive int16 `json:"is_active"`

	CreatedBy int        `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
