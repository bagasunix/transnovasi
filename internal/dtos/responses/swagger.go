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
