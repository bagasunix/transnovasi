package requests

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateCustomer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	Vehicle []CreateVehicle `json:"vehicle,omitempty"`
}

func (c CreateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Phone, validation.Length(0, 14), validation.Match(regexp.MustCompile(`^\d*$`)).Error("phone must be numeric")),
	)
}

type UpdateCustomer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (c UpdateCustomer) Validate() error {
	return validation.ValidateStruct(&c,
		// validation.Field(&c.Email, is.Email),
		validation.Field(&c.ID, validation.Match(regexp.MustCompile(`^\d*$`)).Error("id must be numeric")),
		validation.Field(&c.Phone, validation.Length(0, 14), validation.Match(regexp.MustCompile(`^\d*$`)).Error("phone must be numeric")),
	)
}
