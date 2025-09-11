package responses

// BaseResponseLogin untuk dokumentasi Swagger (menggantikan generic)
type BaseResponseLogin struct {
	Code    int            `json:"code" example:"200"`
	Message string         `json:"message" example:"Pengguna berhasil masuk"`
	Data    *ResponseLogin `json:"data"`
}

// BaseResponseError untuk dokumentasi Swagger (menggantikan BaseResponse[T] kalau error)
type BaseResponseError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"email atau password salah, silahkan coba lagi"`
	Errors  string `json:"errors" example:"record not found"`
}

type BaseResponseUser struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  string        `json:"errors,omitempty"`
	Data    *UserResponse `json:"data,omitempty"`
}

type BaseResponseSwagger struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Pengguna berhasil masuk"`
}

// BaseResponseCustomer digunakan khusus untuk dokumentasi Swagger
type BaseResponseCustomer struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  string            `json:"errors,omitempty"`
	Data    *CustomerResponse `json:"data,omitempty"`
}

// BaseResponseListCustomer khusus untuk dokumentasi Swagger (response list)
type BaseResponseListCustomer struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Errors  string             `json:"errors,omitempty"`
	Data    []CustomerResponse `json:"data,omitempty"`
	Paging  *PageMetadata      `json:"paging,omitempty"`
}

// BaseResponseVehicleList adalah wrapper untuk list vehicle
type BaseResponseVehicleList struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  string            `json:"errors,omitempty"`
	Data    []VehicleResponse `json:"data,omitempty"`
}
