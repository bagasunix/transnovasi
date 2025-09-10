package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BaseRequest struct {
	Search string `json:"search,omitempty"`
	Limit  string `json:"limit,omitempty"`
	Offset string `json:"offset,omitempty"`
	Page   string `json:"page,omitempty"`
}

func (c BaseRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Search, is.Alphanumeric),
		validation.Field(&c.Limit, is.UTFNumeric.Error("Limit harus berupa angka")),
		validation.Field(&c.Offset, is.UTFNumeric.Error("Offset harus berupa angka")),
		validation.Field(&c.Page, is.UTFNumeric.Error("Page harus berupa angka")),
	)
}
