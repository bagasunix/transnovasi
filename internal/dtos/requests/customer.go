package requests

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Customer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	Vehicle []Vehicle `json:"vehicle,omitempty"`
}

func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Phone, validation.Length(0, 14), validation.Match(regexp.MustCompile(`^\d*$`)).Error("phone must be numeric")),
	)
}
